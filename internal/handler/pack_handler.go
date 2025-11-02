package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marcellribeiro/awesomeProject/internal/model"
	"github.com/marcellribeiro/awesomeProject/internal/service"
	log "github.com/sirupsen/logrus"
)

// PackHandler handles HTTP requests for pack calculations
type PackHandler struct {
	service service.PackService
}

// NewPackHandler creates a new pack handler instance
func NewPackHandler(service service.PackService) *PackHandler {
	return &PackHandler{
		service: service,
	}
}

// CalculatePacks handles POST /api/calculate
func (h *PackHandler) CalculatePacks(c *gin.Context) {
	var request model.PackRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid request", err.Error()))
		return
	}

	response, err := h.service.CalculatePackDistribution(&request)
	if err != nil {
		log.Errorf("Calculation failed: %v", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Calculation failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPackSizes handles GET /api/pack-sizes
func (h *PackHandler) GetPackSizes(c *gin.Context) {
	sizes, err := h.service.GetAvailablePackSizes()
	if err != nil {
		log.Errorf("Failed to get pack sizes: %v", err)
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("Failed to get pack sizes", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pack_sizes": sizes,
	})
}

// UpdatePackSizes handles PUT /api/pack-sizes
func (h *PackHandler) UpdatePackSizes(c *gin.Context) {
	var request struct {
		PackSizes []int `json:"pack_sizes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid request", err.Error()))
		return
	}

	if err := h.service.UpdatePackSizes(request.PackSizes); err != nil {
		log.Errorf("Failed to update pack sizes: %v", err)
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Failed to update pack sizes", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Pack sizes updated successfully",
		"pack_sizes": request.PackSizes,
	})
}

// GetDocs handles GET /docs
// Renders API documentation page
func (h *PackHandler) GetDocs(c *gin.Context) {
	c.HTML(http.StatusOK, "docs.html", gin.H{
		"title": "API Documentation",
	})
}

// GetDocsJSON handles GET /docs/json
// Returns API documentation in JSON format
func (h *PackHandler) GetDocsJSON(c *gin.Context) {
	docs := GetAPIDocumentation()
	c.JSON(http.StatusOK, docs)
}

// RenderHome handles GET /
// Renders the main UI page
func (h *PackHandler) RenderHome(c *gin.Context) {
	sizes, _ := h.service.GetAvailablePackSizes()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Pack Calculator",
		"pack_sizes": sizes,
	})
}

// UpdatePackSizesForm handles POST /pack-sizes (form submission)
func (h *PackHandler) UpdatePackSizesForm(c *gin.Context) {
	// Get all pack sizes from form
	err := c.Request.ParseForm()
	if err != nil {
		sizes, _ := h.service.GetAvailablePackSizes()
		log.Errorf("Failed to parse form: %v", err)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":      "Pack Calculator",
			"pack_sizes": sizes,
			"error":      "Failed to parse pack sizes",
		})
		return
	}

	packSizesStr := c.Request.Form["pack_size"]
	if len(packSizesStr) == 0 {
		sizes, _ := h.service.GetAvailablePackSizes()
		log.Error("No pack sizes provided")
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":      "Pack Calculator",
			"pack_sizes": sizes,
			"error":      "At least one pack size is required",
		})
		return
	}

	var packSizes []int
	for _, sizeStr := range packSizesStr {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size <= 0 {
			continue // Skip invalid sizes
		}
		packSizes = append(packSizes, size)
	}

	if len(packSizes) == 0 {
		sizes, _ := h.service.GetAvailablePackSizes()
		log.Error("No valid pack sizes provided")
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":      "Pack Calculator",
			"pack_sizes": sizes,
			"error":      "At least one valid pack size is required (positive numbers only)",
		})
		return
	}

	err = h.service.UpdatePackSizes(packSizes)
	if err != nil {
		sizes, _ := h.service.GetAvailablePackSizes()
		log.Errorf("Failed to update pack sizes: %v", err)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":      "Pack Calculator",
			"pack_sizes": sizes,
			"error":      err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Pack Calculator",
		"pack_sizes": packSizes,
		"success":    "Pack sizes updated successfully!",
	})
}

// CalculatePacksForm handles POST /calculate (form submission)
func (h *PackHandler) CalculatePacksForm(c *gin.Context) {
	quantityStr := c.PostForm("quantity")
	quantity, err := strconv.Atoi(quantityStr)

	if err != nil || quantity <= 0 {
		log.WithField("quantity", quantityStr).Error("Invalid quantity provided")
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title": "Pack Calculator",
			"error": "Please enter a valid quantity (positive number)",
		})
		return
	}

	request := &model.PackRequest{
		Quantity: quantity,
	}

	response, err := h.service.CalculatePackDistribution(request)
	if err != nil {
		sizes, _ := h.service.GetAvailablePackSizes()
		log.WithField("request", request).Errorf("Calculation failed: %v", err)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"title":      "Pack Calculator",
			"pack_sizes": sizes,
			"error":      err.Error(),
			"quantity":   quantity,
		})
		return
	}

	sizes, _ := h.service.GetAvailablePackSizes()
	log.WithFields(log.Fields{
		"response": response,
		"quantity": quantity,
		"sizes":    sizes,
	}).Info("Pack calculation successful")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Pack Calculator",
		"pack_sizes": sizes,
		"quantity":   quantity,
		"result":     response,
	})
}
