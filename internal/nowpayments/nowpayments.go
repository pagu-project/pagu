package nowpayments

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pactus-project/pactus/util/logger"
)

type NowPayments struct {
	ctx       context.Context
	apiToken  string
	ipnSecret []byte
	webhook   string
	apiURL    string
	username  string
	password  string
}

func NewNowPayments(ctx context.Context, cfg *Config) (*NowPayments, error) {
	ipnSecret, err := base64.StdEncoding.DecodeString(cfg.IPNSecret)
	if err != nil {
		return nil, err
	}
	s := &NowPayments{
		ctx:       ctx,
		apiToken:  cfg.APIToken,
		ipnSecret: ipnSecret,
		apiURL:    cfg.APIUrl,
		webhook:   cfg.Webhook,
		username:  cfg.Username,
		password:  cfg.Password,
	}

	// Web hook has issue
	// http.HandleFunc("/nowpayments", s.WebhookFunc)
	//
	// go func() {
	// 	for {
	// 		logger.Info("starting NowPayments webhook", "port", cfg.ListenPort)
	// 		err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.ListenPort), nil)
	// 		if err != nil {
	// 			logger.Error("unable to start NowPayments webhook", "error", err)
	// 		}
	// 	}
	// }()

	return s, nil
}

func (s *NowPayments) WebhookFunc(w http.ResponseWriter, r *http.Request) {
	logger.Debug("NowPayment webhook called")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("Callback read error", "error", err)

		return
	}

	logger.Debug("Callback result", "data", data)
	msgMACHex := r.Header.Get("x-nowpayments-sig")
	msgMAC, err := hex.DecodeString(msgMACHex)
	if err != nil {
		logger.Error("Invalid sig hex", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	mac := hmac.New(sha512.New, s.ipnSecret)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("json.Unmarshal read error", "error", err)

		return
	}

	_, err = mac.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("mac.Write read error", "error", err)

		return
	}
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(expectedMAC, msgMAC) {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("HMAC is invalid", "expectedMAC", expectedMAC, "msgMAC", msgMAC)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *NowPayments) CreatePayment(priceUSD int, orderID string) (string, error) {
	url := fmt.Sprintf("%v/v1/invoice", s.apiURL)
	jsonStr := fmt.Sprintf(`{"price_amount":%v,"price_currency":"usd","order_id":"%v"}`,
		priceUSD, orderID)

	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	logger.Debug("CreatePayment Response", "res", string(data))
	if res.StatusCode != http.StatusOK {
		return "", NowPaymentsError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}
	}

	var resultJSON map[string]any
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return "", err
	}

	return resultJSON["id"].(string), nil
}

func (s *NowPayments) IsPaid(invoiceID string) (bool, error) {
	token, err := s.getJWTToken()
	if err != nil {
		return false, err
	}
	url := fmt.Sprintf("%v/v1/payment/?invoiceId=%v", s.apiURL, invoiceID)
	req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false, err
	}

	req.Header.Set("x-api-key", s.apiToken)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	logger.Debug("ListOfPayments Response", "res", string(data))
	if res.StatusCode != http.StatusOK {
		return false, NowPaymentsError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}
	}

	var resultJSON map[string]any
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return false, err
	}

	results := resultJSON["data"].([]any)
	for _, payment := range results {
		paymentStatus := payment.(map[string]any)["payment_status"]

		if paymentStatus == "finished" {
			return true, nil
		}
	}

	return false, nil
}

func (s *NowPayments) getJWTToken() (string, error) {
	url := fmt.Sprintf("%v/v1/auth", s.apiURL)
	jsonStr := fmt.Sprintf(`{"email":"%v","password":"%v"}`, s.username, s.password)
	req, err := http.NewRequestWithContext(s.ctx, http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	logger.Info("calling NowPayments:auth")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", NowPaymentsError{
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}
	}

	var resultJSON map[string]any
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return "", err
	}

	return resultJSON["token"].(string), nil
}
