package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godiidev/appsynex/internal/domain/services"
	"github.com/godiidev/appsynex/internal/dto/request"
)

type SampleHandler struct {
	sampleService services.SampleService
}

func NewSampleHandler(sampleService services.SampleService) *SampleHandler {
	return &SampleHandler{
		sampleService: sampleService,
	}
}

// GetAll godoc
// @Summary     Get all sample products
// @Description Get a list of sample products with pagination and filtering
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       page query int false "Page number"
// @Param       limit query int false "Items per page"
// @Param       search query string false "Search term for SKU or product name"
// @Param       category query string false "Filter by category ID"
// @Param       sample_type query string false "Filter by sample type"
// @Param       weight_min query number false "Minimum weight"
// @Param       weight_max query number false "Maximum weight"
// @Param       width_min query number false "Minimum width"
// @Param       width_max query number false "Maximum width"
// @Param       color query string false "Filter by color"
// @Security    BearerAuth
// @Success     200 {object} response.PaginatedResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /samples [get]
func (h *SampleHandler) GetAll(c *gin.Context) {
	var req request.SampleFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.sampleService.GetSamples(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary     Get sample by ID
// @Description Get a sample product by ID
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       id path int true "Sample ID"
// @Security    BearerAuth
// @Success     200 {object} response.SampleResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /samples/{id} [get]
func (h *SampleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	sample, err := h.sampleService.GetSampleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sample)
}

// Create godoc
// @Summary     Create a new sample
// @Description Create a new sample product
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       sample body request.CreateSampleRequest true "Sample to create"
// @Security    BearerAuth
// @Success     201 {object} response.SampleResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /samples [post]
func (h *SampleHandler) Create(c *gin.Context) {
	var req request.CreateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample, err := h.sampleService.CreateSample(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sample)
}

// Update godoc
// @Summary     Update a sample
// @Description Update a sample product by ID
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       id path int true "Sample ID"
// @Param       sample body request.UpdateSampleRequest true "Sample data to update"
// @Security    BearerAuth
// @Success     200 {object} response.SampleResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /samples/{id} [put]
func (h *SampleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req request.UpdateSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample, err := h.sampleService.UpdateSample(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sample)
}

// Delete godoc
// @Summary     Delete a sample
// @Description Delete a sample product by ID
// @Tags        samples
// @Accept      json
// @Produce     json
// @Param       id path int true "Sample ID"
// @Security    BearerAuth
// @Success     204 {object} nil
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Failure     403 {object} response.ErrorResponse
// @Failure     404 {object} response.ErrorResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /samples/{id} [delete]
func (h *SampleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.sampleService.DeleteSample(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
