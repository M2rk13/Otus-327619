package log

import (
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"time"
)

type ConversionLog struct {
	id        string
	timestamp time.Time
	request   api.Request
	response  api.Response
}

func NewConversionLog(id string, req api.Request, resp api.Response) *ConversionLog {
	return &ConversionLog{
		id:        id,
		timestamp: time.Now(),
		request:   req,
		response:  resp,
	}
}

func (cl *ConversionLog) ID() string {
	return cl.id
}

func (cl *ConversionLog) Timestamp() time.Time {
	return cl.timestamp
}

func (cl *ConversionLog) Request() api.Request {
	return cl.request
}

func (cl *ConversionLog) Response() api.Response {
	return cl.response
}

func (cl *ConversionLog) SetResponse(resp api.Response) {
	cl.response = resp
}
