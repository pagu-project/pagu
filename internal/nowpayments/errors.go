package nowpayments

import "fmt"

type NowPaymentsError struct {
	Status     string
	StatusCode int
}

func (e NowPaymentsError) Error() string {
	return fmt.Sprintf("error on calling NowPayments API. Status code: %v, status: %v",
		e.StatusCode, e.Status)
}
