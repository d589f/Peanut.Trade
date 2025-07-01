package model

type EstimateRequest struct {
	InputAmount    float64 `json:"inputAmount"`
	InputCurrency  string  `json:"inputCurrency"`
	OutputCurrency string  `json:"outputCurrency"`
}

type EstimateResponse struct {
	ExchangeName string  `json:"exchangeName"`
	OutputAmount float64 `json:"outputAmount"`
}

type GetRatesRequest struct {
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
}

type Rate struct {
	ExchangeName string  `json:"exchangeName"`
	Rate         float64 `json:"rate"`
}

type GetRatesResponse struct {
	Rates []Rate `json:"rates"`
}
