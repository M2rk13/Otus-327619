package repository

import (
	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"
)

type Repository interface {
	CreateRequest(req *api.Request)
	GetRequestByID(id string) *api.Request
	GetAllRequests() []*api.Request
	UpdateRequest(req *api.Request) bool
	DeleteRequest(id string) bool

	CreateResponse(resp *api.Response)
	GetResponseByID(id string) *api.Response
	GetAllResponses() []*api.Response
	UpdateResponse(resp *api.Response) bool
	DeleteResponse(id string) bool

	CreateConversionLog(log *logmodel.ConversionLog)
	GetConversionLogByID(id string) *logmodel.ConversionLog
	GetAllConversionLogs() []*logmodel.ConversionLog
	UpdateConversionLog(log *logmodel.ConversionLog) bool
	DeleteConversionLog(id string) bool

	GetNewConversionRequests() []*api.Request
	GetNewConversionResponses() []*api.Response
	GetNewConversionLogs() []*logmodel.ConversionLog
}
