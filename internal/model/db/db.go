package db

import "time"

type Account struct {
	ID        string
	AccessKey string
}

type RequestCount struct {
	AccountId             string
	RequestCountLastMonth int
}

type RateHistory struct {
	From     string
	To       string
	Rate     float64
	DateTime time.Time
}
