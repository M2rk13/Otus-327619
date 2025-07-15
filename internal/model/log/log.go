package log

import (
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"time"
)

type ConversionLog struct {
	Id        string
	Timestamp time.Time
	Request   api.Request
	Response  api.Response
}

func NewConversionLog(id string, req api.Request, resp api.Response) *ConversionLog {
	return &ConversionLog{
		Id:        id,
		Timestamp: time.Now(),
		Request:   req,
		Response:  resp,
	}
}
