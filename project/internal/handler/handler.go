package handler

import (
	"context"
	"net/http"

	"currency-converter/internal/model"
	"github.com/gin-gonic/gin"
)

type ConversionServiceProvider interface {
	PerformConversion(ctx context.Context, req *model.ConversionAPIRequest) (*model.ConversionAPIResponse, error)
}

type Handler struct {
	service ConversionServiceProvider
}

func NewHandler(service ConversionServiceProvider) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Convert(c *gin.Context) {
	var req model.ConversionAPIRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	result, err := h.service.PerformConversion(c.Request.Context(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, result)
}
