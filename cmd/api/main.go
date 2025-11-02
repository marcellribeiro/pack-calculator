package main

import (
	"os"

	"github.com/marcellribeiro/awesomeProject/internal/handler"
	"github.com/marcellribeiro/awesomeProject/internal/repository"
	"github.com/marcellribeiro/awesomeProject/internal/router"
	"github.com/marcellribeiro/awesomeProject/internal/service"
	"github.com/marcellribeiro/awesomeProject/pkg/calculator"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Printf("📦 Starting Pack Calculator API...")
	// Repository layer - handles data storage
	packRepo := repository.NewInMemoryPackRepository()

	// Calculator - handles the core algorithm
	packCalc := calculator.NewDynamicPackCalculator()

	// Service layer - handles business logic
	packService := service.NewPackService(packCalc, packRepo)

	// Handler layer - handles HTTP requests
	packHandler := handler.NewPackHandler(packService)

	// Setup Gin router
	ginRouter := router.SetupRouter(packHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Open http://localhost:%s in your browser", port)
	log.Printf("📦 Pack Calculator API is ready!")

	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
