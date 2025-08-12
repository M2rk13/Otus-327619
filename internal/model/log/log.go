package log

import (
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
)

type ConversionLog struct {
	Id        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	Request   api.Request  `json:"request"`
	Response  api.Response `json:"response"`
}

func (c *ConversionLog) GetId() string {
	return c.Id
}

func NewConversionLog(id string, req api.Request, resp api.Response) *ConversionLog {
	return &ConversionLog{
		Id:        id,
		Timestamp: time.Now(),
		Request:   req,
		Response:  resp,
	}
}
