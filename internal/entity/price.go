package entity

type Price struct {
	AzbitPacToUSDT     AzbitPriceResponse
	TradeOgrePacToUSDT TradeOgrePriceResponse
}

//nolint:tagliatelle // External dependency
type AzbitPriceResponse struct {
	Timestamp                int     `json:"timestamp"`
	CurrencyPairCode         string  `json:"currencyPairCode"`
	Price                    float64 `json:"price"`
	Price24HAgo              float64 `json:"price24hAgo"`
	PriceChangePercentage24H float64 `json:"priceChangePercentage24h"`
	Volume24H                float64 `json:"volume24h"`
	BidPrice                 float64 `json:"bidPrice"`
	AskPrice                 float64 `json:"askPrice"`
	Low24H                   float64 `json:"low24h"`
	High24H                  float64 `json:"high24h"`
}

type TradeOgrePriceResponse struct {
	InitialPrice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}
