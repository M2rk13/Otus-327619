package api

import (
	"github.com/M2rk13/Otus-327619/internal/model/db"
)

type Auth struct {
	accessKey string
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

func SetAccessKey(account db.Account, auth *Auth) {
	auth.accessKey = account.AccessKey
}

func GetAccessKey(auth *Auth) string {
	return auth.accessKey
}
