package main

import (
	"fmt"
	dispatcher "github.com/M2rk13/Otus-327619/internal/service"
	"time"
)

func main() {
	for i := 0; i < 3; i++ {
		dispatcher.DispatchExampleData(i)
		time.Sleep(1 * time.Second)
		fmt.Print("Iteration", i+1, "completed.\n")
	}

	fmt.Print("Finished dispatching data.")
}
