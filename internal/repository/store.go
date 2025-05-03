package store

import (
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
)

var (
	conversionRequests  []*api.Request
	conversionResponses []*api.Response
	conversionLogs      []*log.ConversionLog
)

func Store(data interface{}) {
	switch val := data.(type) {
	case *api.Request:
		conversionRequests = append(conversionRequests, val)
		fmt.Println("Added ConversionRequest to repository.")
	case *api.Response:
		conversionResponses = append(conversionResponses, val)
		fmt.Println("Added ConversionResponse to repository.")
	case *log.ConversionLog:
		conversionLogs = append(conversionLogs, val)
		fmt.Println("Added ConversionLog to repository.")
	default:
		fmt.Println("Unknown type passed to Store.")
	}
}
