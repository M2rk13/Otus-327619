package main

import (
	"fmt"
	repository "github.com/M2rk13/Otus-327619/internal/repository"
	dispatcher "github.com/M2rk13/Otus-327619/internal/service"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	repository.StartStorageGoRoutines(&wg)
	dispatcher.StartSliceLogger(&wg)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			dispatcher.DispatchExampleData(i)
			time.Sleep(500 * time.Millisecond)
			fmt.Print("Iteration", i+1, "dispatched.\n")
		}

		close(repository.RequestChan)
		repository.RequestChanState = 0

		close(repository.ResponseChan)
		repository.ResponseChanState = 0

		close(repository.LogChan)
		repository.LogChanState = 0

		fmt.Println("All data dispatched and channels closed.")
	}()

	wg.Wait()

	fmt.Print("Finished application. All goroutines completed.\n")
}
