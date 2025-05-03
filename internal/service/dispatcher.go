package dispatcher

import (
	"fmt"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	store "github.com/M2rk13/Otus-327619/internal/repository"
	"time"
)

func DispatchExampleData(iteration int) {
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

	convLog := log.NewConversionLog("log_001", req, resp)

	fmt.Println("Dispatching data ...")

	store.Store(req)
	store.Store(resp)
	store.Store(*convLog)
}
