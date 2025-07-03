package dispatcher

import (
	"fmt"
	repository "github.com/M2rk13/Otus-327619/internal/repository"
	"sync"
	"time"
)

func StartSliceLogger(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		fmt.Println("Slice logger started.")

		for range ticker.C {
			newRequests := repository.GetNewConversionRequests()

			if len(newRequests) > 0 {
				fmt.Println("--- New Conversion Requests ---")

				for _, req := range newRequests {
					fmt.Printf("  Request: From=%s, To=%s, Amount=%.2f\n", req.From, req.To, req.Amount)
				}
			}

			newResponses := repository.GetNewConversionResponses()

			if len(newResponses) > 0 {
				fmt.Println("--- New Conversion Responses ---")

				for _, resp := range newResponses {
					fmt.Printf("  Response: Success=%t, Result=%.2f\n", resp.Success, resp.Result)
				}
			}

			newLogs := repository.GetNewConversionLogs()

			if len(newLogs) > 0 {
				fmt.Println("--- New Conversion Logs ---")

				for _, l := range newLogs {
					fmt.Printf(
						"  Log: ID=%s, Timestamp=%s, RequestFrom=%s, ResponseResult=%.2f\n",
						l.ID(),
						l.Timestamp().Format(time.RFC3339),
						l.Request().From,
						l.Response().Result)
				}
			}

			if repository.RequestChanState == 0 && repository.ResponseChanState == 0 && repository.LogChanState == 0 {
				return
			}
		}

		fmt.Println("Slice logger stopped.")
	}()
}
