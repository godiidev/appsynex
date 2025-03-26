package main

import (
	"fmt"
	"log"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/api/router"
	"github.com/godiidev/appsynex/internal/repository/mysql"
)

// @title           AppSynex API
// @version         1.0
// @description     API Server for AppSynex Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@appsynex.vn

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := mysql.NewDBConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup router
	r := router.SetupRouter(db, cfg)

	// Start server
	port := cfg.Server.Port
	log.Printf("Server started on port %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
