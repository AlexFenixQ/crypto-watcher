package model

type CurrencyRequest struct {
	Coin string `json:"coin" example:"bitcoin"`
}

type PriceRequest struct {
	Coin      string `json:"coin" example:"bitcoin"`
	Timestamp int64  `json:"timestamp" example:"1723041724"`
}

type PriceResponse struct {
	Coin     string  `json:"coin"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
