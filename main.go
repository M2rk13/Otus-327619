package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/service"
)

func main() {
	var wg sync.WaitGroup

	var (
		requestChan  chan *api.Request
		responseChan chan *api.Response
		logChan      chan *log.ConversionLog
	)

	var (
		requestChanState  int
		responseChanState int
		logChanState      int
	)

	requestChan = make(chan *api.Request, 10)
	requestChanState = 1

	responseChan = make(chan *api.Response, 10)
	responseChanState = 1

	logChan = make(chan *log.ConversionLog, 10)
	logChanState = 1

	service.StartStorageService(&wg, requestChan, responseChan, logChan)
	service.StartSliceLogger(&wg, &requestChanState, &responseChanState, &logChanState)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			service.DispatchExampleData(i, requestChan, responseChan, logChan)
			time.Sleep(500 * time.Millisecond)
			fmt.Print("Iteration", i+1, "dispatched.\n")
		}

		close(requestChan)
		requestChanState = 0

		close(responseChan)
		responseChanState = 0

		close(logChan)
		logChanState = 0

		fmt.Println("All data dispatched and channels closed.")
	}()

	wg.Wait()

	fmt.Print("Finished application. All goroutines completed.\n")
}
