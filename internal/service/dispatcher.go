package service

import (
	"fmt"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
)

func DispatchExampleData(
	iteration int,
	requestChan chan<- *api.Request,
	responseChan chan<- *api.Response,
	logChan chan<- *log.ConversionLog,
) {
	amount := float64(105 * iteration)

	req := api.Request{
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
			Success: true,
			Terms:   "https://exchangerate.host/terms",
			Privacy: "https://exchangerate.host/privacy",
			Query:   req,
			Info:    info,
			Result:  info.Quote * amount,
		}
	} else {
		resp = api.Response{
			Success: false,
			Terms:   "https://exchangerate.host/terms",
			Privacy: "https://exchangerate.host/privacy",
			Query:   req,
			Info:    info,
			Result:  0.0,
		}
	}

	convLog := log.NewConversionLog(fmt.Sprintf("log_%03d", iteration+1), req, resp)

	fmt.Println("Dispatching data ...")

	requestChan <- &req
	responseChan <- &resp
	logChan <- convLog
}
