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

	"github.com/M2rk13/Otus-327619/internal/service"
)

func StartWebServer(ctx context.Context, wg *sync.WaitGroup, addr string, storageSvc *service.StorageService) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		// Роуты Swagger возвращены
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.POST("/api/auth/login", loginHandler)

		protected := router.Group("/api")
		protected.Use(authMiddleware())

		// Роуты теперь используют созданный внутри APIHandler
		apiHandler := NewAPIHandler(storageSvc)

		protected.POST("/requests", apiHandler.createRequest)
		protected.PUT("/requests/:id", apiHandler.updateRequest)
		router.GET("/api/requests", apiHandler.getAllRequests)
		router.GET("/api/requests/:id", apiHandler.getRequestByID)
		protected.DELETE("/requests/:id", apiHandler.deleteRequest)

		protected.POST("/responses", apiHandler.createResponse)
		protected.PUT("/responses/:id", apiHandler.updateResponse)
		router.GET("/api/responses", apiHandler.getAllResponses)
		router.GET("/api/responses/:id", apiHandler.getResponseByID)
		protected.DELETE("/responses/:id", apiHandler.deleteResponse)

		protected.POST("/logs", apiHandler.createLog)
		protected.PUT("/logs/:id", apiHandler.updateLog)
		router.GET("/api/logs", apiHandler.getAllLogs)
		router.GET("/api/logs/:id", apiHandler.getLogByID)
		protected.DELETE("/logs/:id", apiHandler.deleteLog)

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
