package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/M2rk13/Otus-327619/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func StartWebServer(ctx context.Context, wg *sync.WaitGroup, addr string) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.POST("/api/auth/login", loginHandler)

		protected := router.Group("/api")
		protected.Use(authMiddleware())

		protected.POST("/requests", createRequest)
		protected.PUT("/requests/:id", updateRequest)
		router.GET("/api/requests", getAllRequests)
		router.GET("/api/requests/:id", getRequestByID)
		protected.DELETE("/requests/:id", deleteRequest)

		protected.POST("/responses", createResponse)
		protected.PUT("/responses/:id", updateResponse)
		router.GET("/api/responses", getAllResponses)
		router.GET("/api/responses/:id", getResponseByID)
		protected.DELETE("/responses/:id", deleteResponse)

		protected.POST("/logs", createLog)
		protected.PUT("/logs/:id", updateLog)
		router.GET("/api/logs", getAllLogs)
		router.GET("/api/logs/:id", getLogByID)
		protected.DELETE("/logs/:id", deleteLog)

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
