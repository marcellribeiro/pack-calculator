package router

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/marcellribeiro/awesomeProject/internal/handler"
)

func SetupRouter(handler *handler.PackHandler) *gin.Engine {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("web/templates/*")

	// Serve static files
	router.Static("/static", "./web/static")

	// Web UI routes
	router.GET("/", handler.RenderHome)
	router.POST("/calculate", handler.CalculatePacksForm)
	router.POST("/pack-sizes", handler.UpdatePackSizesForm)

	// Documentation routes
	router.GET("/docs", handler.GetDocs)
	router.GET("/docs/json", handler.GetDocsJSON)

	// API routes
	api := router.Group("/api")
	{
		api.POST("/calculate", handler.CalculatePacks)
		api.GET("/pack-sizes", handler.GetPackSizes)
		api.PUT("/pack-sizes", handler.UpdatePackSizes)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "pack-calculator",
		})
	})

	return router
}
