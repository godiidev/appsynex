package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/repository/mysql"
)

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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	// Read migration files
	files := []string{
		"migrations/000001_init_schema.up.sql",
		"migrations/000002_product_tables.up.sql",
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		fmt.Printf("Running migration: %s\n", file)
		_, err = sqlDB.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
		fmt.Printf("Successfully ran migration: %s\n", file)
	}

	fmt.Println("All migrations completed successfully!")
}
