package model

import "time"

type ConversionAPIRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	From   string  `json:"from" binding:"required,len=3"`
	To     string  `json:"to" binding:"required,len=3"`
}

type ConversionAPIResponse struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Result float64 `json:"result"`
	Rate   float64 `json:"rate"`
}

type exchangeRateResponse struct {
	Success bool `json:"success"`
	Info    struct {
		Rate float64 `json:"rate"`
	} `json:"info"`
	Result float64 `json:"result"`
	Error  struct {
		Code int    `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

type ConversionHistory struct {
	ID        int64
	From      string
	To        string
	Amount    float64
	Result    float64
	Rate      float64
	CreatedAt time.Time
}
