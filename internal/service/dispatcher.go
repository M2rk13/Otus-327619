package service

import (
	"fmt"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"

	"github.com/google/uuid"
)

func DispatchExampleData(
	iteration int,
	requestChan chan<- *api.Request,
	responseChan chan<- *api.Response,
	logChan chan<- *log.ConversionLog,
) {
	amount := float64(105 * iteration)

	req := api.Request{
		Id:     uuid.New().String(),
		From:   "USD",
		To:     "EUR",
		Amount: amount,
	}

	info := api.Info{
		Timestamp: time.Now().Unix(),
		Quote:     0.9134,
	}

	var resp api.Response

	if iteration%2 == 0 {
		resp = api.Response{
			Id:      uuid.New().String(),
			Success: true,
			Terms:   "https://exchangerate.host/terms",
			Privacy: "https://exchangerate.host/privacy",
			Query:   req,
			Info:    info,
			Result:  info.Quote * amount,
		}
	} else {
		resp = api.Response{
			Id:      uuid.New().String(),
			Success: false,
			Terms:   "https://exchangerate.host/terms",
			Privacy: "https://exchangerate.host/privacy",
			Query:   req,
			Info:    info,
			Result:  0.0,
		}
	}

	convLog := log.NewConversionLog(uuid.New().String(), req, resp)

	fmt.Println("Dispatching data ...")

	requestChan <- &req
	responseChan <- &resp
	logChan <- convLog
}
