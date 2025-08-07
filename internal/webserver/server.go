package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func StartWebServer(ctx context.Context, wg *sync.WaitGroup, addr string) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		router.POST("/api/requests", createHandler(requestService))
		router.PUT("/api/requests/:id", updateHandler(requestService))
		router.GET("/api/requests", getAllHandler(requestService))
		router.GET("/api/requests/:id", getByIDHandler(requestService))
		router.DELETE("/api/requests/:id", deleteHandler(requestService))

		router.POST("/api/responses", createHandler(responseService))
		router.PUT("/api/responses/:id", updateHandler(responseService))
		router.GET("/api/responses", getAllHandler(responseService))
		router.GET("/api/responses/:id", getByIDHandler(responseService))
		router.DELETE("/api/responses/:id", deleteHandler(responseService))

		router.POST("/api/logs", createHandler(logService))
		router.PUT("/api/logs/:id", updateHandler(logService))
		router.GET("/api/logs", getAllHandler(logService))
		router.GET("/api/logs/:id", getByIDHandler(logService))
		router.DELETE("/api/logs/:id", deleteHandler(logService))

		server := &http.Server{
			Addr:    addr,
			Handler: router,
		}

		go func() {
			fmt.Printf("Web server started on %s\n", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Web server listen error: %v", err)
			}
		}()

		<-ctx.Done()
		fmt.Println("Web server is shutting down...")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Web server graceful shutdown failed: %v\n", err)
		} else {
			fmt.Println("Web server gracefully stopped.")
		}
	}()
}
