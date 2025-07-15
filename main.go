package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/service"
)

type chanItem[T any] struct {
	ch    chan T
	state int
}

var (
	requestChan  *chanItem[*api.Request]
	responseChan *chanItem[*api.Response]
	logChan      *chanItem[*log.ConversionLog]
)

func init() {
	requestChan = &chanItem[*api.Request]{}
	requestChan.ch = make(chan *api.Request, 10)
	requestChan.state = 1

	responseChan = &chanItem[*api.Response]{}
	responseChan.ch = make(chan *api.Response, 10)
	responseChan.state = 1

	logChan = &chanItem[*log.ConversionLog]{}
	logChan.ch = make(chan *log.ConversionLog, 10)
	logChan.state = 1
}

func main() {
	var wg sync.WaitGroup

	appTimeOut := 10 * time.Second
	ctx, cancelFunc := context.WithTimeout(context.Background(), appTimeOut)
	defer cancelFunc()

	wg.Add(1)

	go func() {
		defer wg.Done()

		osSignalChan := make(chan os.Signal, 1)
		signal.Notify(osSignalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-osSignalChan:
			fmt.Printf("\nOS signal: %v. Shutting down...\n", sig)
			cancelFunc()
		case <-ctx.Done():
			if errors.Is(context.DeadlineExceeded, ctx.Err()) {
				fmt.Println("Stopped by timeout.")
			}
		}
	}()

	service.StartStorageService(&wg, ctx, requestChan.ch, responseChan.ch, logChan.ch)
	service.StartSliceLogger(&wg, ctx, &requestChan.state, &responseChan.state, &logChan.state)

	wg.Add(1)
	doForever(&wg, ctx)

	go func() {
		wg.Wait()
		time.Sleep(700 * time.Millisecond)
		cancelFunc()
	}()

	<-ctx.Done()
	wg.Wait()
	fmt.Print("Finished application.\n")
}

func doForever(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Dispatcher stopped due context cancel.")
			return
		default:
			service.DispatchExampleData(i, requestChan.ch, responseChan.ch, logChan.ch)
			time.Sleep(500 * time.Millisecond)
			fmt.Print("Iteration", i+1, "dispatched.\n")
		}
	}

	close(requestChan.ch)
	requestChan.state = 0

	close(responseChan.ch)
	responseChan.state = 0

	close(logChan.ch)
	logChan.state = 0

	fmt.Println("All data dispatched and channels closed.")
}
