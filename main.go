package main

import (
	_ "github.com/M2rk13/Otus-327619/internal/bootstrap"

	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/M2rk13/Otus-327619/internal/config"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"
	"github.com/M2rk13/Otus-327619/internal/service"
	"github.com/M2rk13/Otus-327619/internal/webserver"
)

type chanItem[T any] struct {
	ch    chan T
	state int
}

var (
	requestChan  *chanItem[*api.Request]
	responseChan *chanItem[*api.Response]
	logChan      *chanItem[*logmodel.ConversionLog]
)

func init() {
	requestChan = &chanItem[*api.Request]{}
	requestChan.ch = make(chan *api.Request, 10)
	requestChan.state = 1

	responseChan = &chanItem[*api.Response]{}
	responseChan.ch = make(chan *api.Response, 10)
	responseChan.state = 1

	logChan = &chanItem[*logmodel.ConversionLog]{}
	logChan.ch = make(chan *logmodel.ConversionLog, 10)
	logChan.state = 1
}

func main() {
	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var store repository.Repository
	var err error

	log.Printf("Using storage type: %s\n", config.AppCfg.StorageType)

	switch config.AppCfg.StorageType {
	case "mongo":
		mongoStore, mongoErr := repository.NewMongoStore(ctx, config.MongoCfg, config.RedisCfg)
		if mongoErr != nil {
			log.Fatalf("Failed to setup mongo persistence: %v", mongoErr)
		}
		store = mongoStore
		defer mongoStore.Close(ctx)
	case "file":
		fileStore := repository.NewFileStore()
		if err = fileStore.SetupPersistence(); err != nil {
			log.Fatalf("Failed to setup file persistence: %v", err)
		}
		store = fileStore
		defer fileStore.ClosePersistence()
	default:
		log.Fatalf("Unknown storage type: %s", config.AppCfg.StorageType)
	}

	dispatcherService := service.NewDispatcherService()
	storageService := service.NewStorageService(store)
	loggerService := service.NewLoggerService(store)

	wg.Add(1)

	go func() {
		defer wg.Done()

		osSignalChan := make(chan os.Signal, 1)
		signal.Notify(osSignalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-osSignalChan:
			fmt.Printf("\nOS signal received: %v. Shutting down...\n", sig)
			cancelFunc()
		case <-ctx.Done():
		}
	}()

	storageService.StartStorageService(&wg, ctx, requestChan.ch, responseChan.ch, logChan.ch)
	loggerService.StartSliceLogger(&wg, ctx, &requestChan.state, &responseChan.state, &logChan.state)
	webserver.StartWebServer(ctx, &wg, ":8080", storageService)

	wg.Add(1)
	go doForever(&wg, ctx, dispatcherService)

	wg.Wait()

	fmt.Print("Finished application. All goroutines completed.\n")
}

func doForever(wg *sync.WaitGroup, ctx context.Context, dispatcher *service.DispatcherService) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Dispatcher stopped by context.")

			return
		default:
			dispatcher.DispatchExampleData(i, requestChan.ch, responseChan.ch, logChan.ch)
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Iteration %d finished.\n", i+1)
		}
	}

	close(requestChan.ch)
	requestChan.state = 0

	close(responseChan.ch)
	responseChan.state = 0

	close(logChan.ch)
	logChan.state = 0

	fmt.Println("All data sent, channels closed.")
}
