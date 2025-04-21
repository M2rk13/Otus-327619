package api

import (
	"github.com/M2rk13/Otus-327619/internal/model/db"
)

type auth struct {
	AccessKey string
}

type Query struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type Info struct {
	Timestamp int64   `json:"timestamp"`
	Quote     float64 `json:"quote"`
}

type ConvertResponse struct {
	Success bool    `json:"success"`
	Terms   string  `json:"terms"`
	Privacy string  `json:"privacy"`
	Query   Query   `json:"query"`
	Info    Info    `json:"info"`
	Result  float64 `json:"result"`
}

func GetResult(response *ConvertResponse) float64 {
	return response.Result
}

func GetRate(response *ConvertResponse) float64 {
	return response.Info.Quote
}

func setFrom(query *Query, currencyFrom string) {
	query.From = currencyFrom
}

func setTo(query *Query, currencyTo string) {
	query.To = currencyTo
}

func setAmount(query *Query, amount float64) {
	query.Amount = amount
}

func SetAccessKey(account db.Account, auth *auth) {
	auth.AccessKey = account.AccessKey
}
