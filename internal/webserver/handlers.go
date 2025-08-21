package webserver

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/service"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	storageSvc *service.StorageService
}

func NewAPIHandler(storageSvc *service.StorageService) *APIHandler {
	return &APIHandler{storageSvc: storageSvc}
}

// @Summary      Create request
// @Description  Creates a new currency conversion request
// @Tags         requests
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      api.Request  true  "Request to create"
// @Success      201      {object}  api.Request
// @Failure      400      {object}  object{error=string}
// @Router       /requests [post]
func (h *APIHandler) createRequest(c *gin.Context) {
	var req api.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	h.storageSvc.CreateRequest(&req)
	c.JSON(http.StatusCreated, req)
}

// @Summary      Get request by ID
// @Description  Retrieves a request by its ID
// @Tags         requests
// @Produce      json
// @Param        id  path  string  true  "Request ID"
// @Success      200 {object} api.Request
// @Failure      404 {object} object{error=string}
// @Router       /requests/{id} [get]
func (h *APIHandler) getRequestByID(c *gin.Context) {
	id := c.Param("id")
	item := h.storageSvc.GetRequestByID(id)

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})

		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary      Get all requests
// @Description  Retrieves all requests
// @Tags         requests
// @Produce      json
// @Success      200 {array}  api.Request
// @Router       /requests [get]
func (h *APIHandler) getAllRequests(c *gin.Context) {
	items := h.storageSvc.GetAllRequests()
	c.JSON(http.StatusOK, items)
}

// @Summary      Update request
// @Description  Updates a request by ID
// @Tags         requests
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      string       true  "Request ID"
// @Param        request  body      api.Request  true  "Updated request"
// @Success      200      {object}  api.Request
// @Failure      400      {object}  object{error=string}
// @Failure      404      {object}  object{error=string}
// @Router       /requests/{id} [put]
func (h *APIHandler) updateRequest(c *gin.Context) {
	id := c.Param("id")
	var updatedItem api.Request

	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	v := reflect.ValueOf(updatedItem).Elem().FieldByName("Id")

	if v.IsValid() && v.String() != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match Id in path"})

		return
	}

	if h.storageSvc.UpdateRequest(&updatedItem) {
		c.JSON(http.StatusOK, updatedItem)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

// @Summary      Delete request
// @Description  Deletes a request by ID
// @Tags         requests
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  string  true  "Request ID"
// @Success      200 {object} object{message=string}
// @Failure      404 {object} object{error=string}
// @Router       /requests/{id} [delete]
func (h *APIHandler) deleteRequest(c *gin.Context) {
	id := c.Param("id")

	if h.storageSvc.DeleteRequest(id) {
		c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

// @Summary      Create response
// @Description  Creates a new conversion response
// @Tags         responses
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        response  body      api.Response  true  "Response to create"
// @Success      201       {object}  api.Response
// @Failure      400       {object}  object{error=string}
// @Router       /responses [post]
func (h *APIHandler) createResponse(c *gin.Context) {
	var resp api.Response

	if err := c.ShouldBindJSON(&resp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	h.storageSvc.CreateResponse(&resp)
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Get response by ID
// @Description  Retrieves a response by its ID
// @Tags         responses
// @Produce      json
// @Param        id  path  string  true  "Response ID"
// @Success      200 {object} api.Response
// @Failure      404 {object} object{error=string}
// @Router       /responses/{id} [get]
func (h *APIHandler) getResponseByID(c *gin.Context) {
	id := c.Param("id")
	item := h.storageSvc.GetResponseByID(id)

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})

		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary      Get all responses
// @Description  Retrieves all responses
// @Tags         responses
// @Produce      json
// @Success      200 {array}  api.Response
// @Router       /responses [get]
func (h *APIHandler) getAllResponses(c *gin.Context) {
	items := h.storageSvc.GetAllResponses()
	c.JSON(http.StatusOK, items)
}

// @Summary      Update response
// @Description  Updates a response by ID
// @Tags         responses
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id        path      string       true  "Response ID"
// @Param        response  body      api.Response true  "Updated response"
// @Success      200       {object}  api.Response
// @Failure      400       {object}  object{error=string}
// @Failure      404       {object}  object{error=string}
// @Router       /responses/{id} [put]
func (h *APIHandler) updateResponse(c *gin.Context) {
	id := c.Param("id")
	var updatedItem api.Response

	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	if updatedItem.Id != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match Id in path"})

		return
	}

	if h.storageSvc.UpdateResponse(&updatedItem) {
		c.JSON(http.StatusOK, updatedItem)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

// @Summary      Delete response
// @Description  Deletes a response by ID
// @Tags         responses
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  string  true  "Response ID"
// @Success      200 {object} object{message=string}
// @Failure      404 {object} object{error=string}
// @Router       /responses/{id} [delete]
func (h *APIHandler) deleteResponse(c *gin.Context) {
	id := c.Param("id")

	if h.storageSvc.DeleteResponse(id) {
		c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

// @Summary      Create log record
// @Description  Creates a new conversion log record
// @Tags         logs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        log  body      log.ConversionLog  true  "Log to create"
// @Success      201  {object}  log.ConversionLog
// @Failure      400  {object}  object{error=string}
// @Router       /logs [post]
func (h *APIHandler) createLog(c *gin.Context) {
	var logItem log.ConversionLog

	if err := c.ShouldBindJSON(&logItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	h.storageSvc.CreateConversionLog(&logItem)
	c.JSON(http.StatusCreated, logItem)
}

// @Summary      Get log by ID
// @Description  Retrieves a conversion log by its ID
// @Tags         logs
// @Produce      json
// @Param        id  path  string  true  "Log ID"
// @Success      200 {object} log.ConversionLog
// @Failure      404 {object} object{error=string}
// @Router       /logs/{id} [get]
func (h *APIHandler) getLogByID(c *gin.Context) {
	id := c.Param("id")
	item := h.storageSvc.GetConversionLogByID(id)

	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})

		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary      Get all logs
// @Description  Retrieves all conversion logs
// @Tags         logs
// @Produce      json
// @Success      200 {array}  log.ConversionLog
// @Router       /logs [get]
func (h *APIHandler) getAllLogs(c *gin.Context) {
	items := h.storageSvc.GetAllConversionLogs()
	c.JSON(http.StatusOK, items)
}

// @Summary      Update log
// @Description  Updates a conversion log by ID
// @Tags         logs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string            true  "Log ID"
// @Param        log  body      log.ConversionLog true  "Updated log"
// @Success      200  {object}  log.ConversionLog
// @Failure      400  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Router       /logs/{id} [put]
func (h *APIHandler) updateLog(c *gin.Context) {
	id := c.Param("id")
	var updatedItem log.ConversionLog

	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

		return
	}

	if updatedItem.Id != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match Id in path"})

		return
	}

	if h.storageSvc.UpdateConversionLog(&updatedItem) {
		c.JSON(http.StatusOK, updatedItem)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}

// @Summary      Delete log
// @Description  Deletes a conversion log by ID
// @Tags         logs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id  path  string  true  "Log ID"
// @Success      200 {object} object{message=string}
// @Failure      404 {object} object{error=string}
// @Router       /logs/{id} [delete]
func (h *APIHandler) deleteLog(c *gin.Context) {
	id := c.Param("id")

	if h.storageSvc.DeleteConversionLog(id) {
		c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
	}
}
