package webserver

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/M2rk13/Otus-327619/internal/model/api"
	"github.com/M2rk13/Otus-327619/internal/model/log"
	"github.com/M2rk13/Otus-327619/internal/repository"

	"github.com/gin-gonic/gin"
)

type crudService[T any] interface {
	Create(T)
	GetByID(string) T
	GetAll() []T
	Update(T) bool
	Delete(string) bool
}

func createHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newItem T

		if err := c.ShouldBindJSON(&newItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

			return
		}

		repo.Create(newItem)
		c.JSON(http.StatusCreated, newItem)
	}
}

func getByIDHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		item := repo.GetByID(id)

		if reflect.ValueOf(item).IsNil() {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})

			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func getAllHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		items := repo.GetAll()
		c.JSON(http.StatusOK, items)
	}
}

func updateHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedItem T

		if err := c.ShouldBindJSON(&updatedItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})

			return
		}

		v := reflect.ValueOf(updatedItem).Elem().FieldByName("Id")

		if v.IsValid() && v.String() != id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Id in body must match Id in path"})

			return
		}

		if repo.Update(updatedItem) {
			c.JSON(http.StatusOK, updatedItem)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
	}
}

func deleteHandler[T any](repo crudService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if repo.Delete(id) {
			c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		}
	}
}

type RequestService struct{}

func (s *RequestService) Create(req *api.Request)        { repository.CreateRequest(req) }
func (s *RequestService) GetByID(id string) *api.Request { return repository.GetRequestByID(id) }
func (s *RequestService) GetAll() []*api.Request         { return repository.GetAllRequests() }
func (s *RequestService) Update(req *api.Request) bool   { return repository.UpdateRequest(req) }
func (s *RequestService) Delete(id string) bool          { return repository.DeleteRequest(id) }

type ResponseService struct{}

func (s *ResponseService) Create(resp *api.Response)       { repository.CreateResponse(resp) }
func (s *ResponseService) GetByID(id string) *api.Response { return repository.GetResponseByID(id) }
func (s *ResponseService) GetAll() []*api.Response         { return repository.GetAllResponses() }
func (s *ResponseService) Update(resp *api.Response) bool  { return repository.UpdateResponse(resp) }
func (s *ResponseService) Delete(id string) bool           { return repository.DeleteResponse(id) }

type LogService struct{}

func (s *LogService) Create(log *log.ConversionLog) { repository.CreateConversionLog(log) }
func (s *LogService) GetByID(id string) *log.ConversionLog {
	return repository.GetConversionLogByID(id)
}
func (s *LogService) GetAll() []*log.ConversionLog       { return repository.GetAllConversionLogs() }
func (s *LogService) Update(log *log.ConversionLog) bool { return repository.UpdateConversionLog(log) }
func (s *LogService) Delete(id string) bool              { return repository.DeleteConversionLog(id) }

var (
	requestService  = &RequestService{}
	responseService = &ResponseService{}
	logService      = &LogService{}
)

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
func createRequest(c *gin.Context) {
	createHandler[*api.Request](requestService)(c)
}

// @Summary      Get request by ID
// @Description  Retrieves a request by its ID
// @Tags         requests
// @Produce      json
// @Param        id  path  string  true  "Request ID"
// @Success      200 {object} api.Request
// @Failure      404 {object} object{error=string}
// @Router       /requests/{id} [get]
func getRequestByID(c *gin.Context) {
	getByIDHandler[*api.Request](requestService)(c)
}

// @Summary      Get all requests
// @Description  Retrieves all requests
// @Tags         requests
// @Produce      json
// @Success      200 {array}  api.Request
// @Router       /requests [get]
func getAllRequests(c *gin.Context) {
	getAllHandler[*api.Request](requestService)(c)
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
func updateRequest(c *gin.Context) {
	updateHandler[*api.Request](requestService)(c)
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
func deleteRequest(c *gin.Context) {
	deleteHandler[*api.Request](requestService)(c)
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
func createResponse(c *gin.Context) {
	createHandler[*api.Response](responseService)(c)
}

// @Summary      Get response by ID
// @Description  Retrieves a response by its ID
// @Tags         responses
// @Produce      json
// @Param        id  path  string  true  "Response ID"
// @Success      200 {object} api.Response
// @Failure      404 {object} object{error=string}
// @Router       /responses/{id} [get]
func getResponseByID(c *gin.Context) {
	getByIDHandler[*api.Response](responseService)(c)
}

// @Summary      Get all responses
// @Description  Retrieves all responses
// @Tags         responses
// @Produce      json
// @Success      200 {array}  api.Response
// @Router       /responses [get]
func getAllResponses(c *gin.Context) {
	getAllHandler[*api.Response](responseService)(c)
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
func updateResponse(c *gin.Context) {
	updateHandler[*api.Response](responseService)(c)
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
func deleteResponse(c *gin.Context) {
	deleteHandler[*api.Response](responseService)(c)
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
func createLog(c *gin.Context) {
	createHandler[*log.ConversionLog](logService)(c)
}

// @Summary      Get log by ID
// @Description  Retrieves a conversion log by its ID
// @Tags         logs
// @Produce      json
// @Param        id  path  string  true  "Log ID"
// @Success      200 {object} log.ConversionLog
// @Failure      404 {object} object{error=string}
// @Router       /logs/{id} [get]
func getLogByID(c *gin.Context) {
	getByIDHandler[*log.ConversionLog](logService)(c)
}

// @Summary      Get all logs
// @Description  Retrieves all conversion logs
// @Tags         logs
// @Produce      json
// @Success      200 {array}  log.ConversionLog
// @Router       /logs [get]
func getAllLogs(c *gin.Context) {
	getAllHandler[*log.ConversionLog](logService)(c)
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
func updateLog(c *gin.Context) {
	updateHandler[*log.ConversionLog](logService)(c)
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
func deleteLog(c *gin.Context) {
	deleteHandler[*log.ConversionLog](logService)(c)
}
