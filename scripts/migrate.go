// File: scripts/migrate.go
// Tạo tại: scripts/migrate.go
// Mục đích: Tool để chạy migration tự động, hỗ trợ up/down migration

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/godiidev/appsynex/config"
	"github.com/godiidev/appsynex/internal/repository/mysql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run scripts/migrate.go [up|down]")
	}

	command := os.Args[1]

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

	switch command {
	case "up":
		runUpMigrations(sqlDB)
	case "down":
		runDownMigrations(sqlDB)
	default:
		log.Fatal("Invalid command. Use 'up' or 'down'")
	}
}

func runUpMigrations(db interface{}) {
	fmt.Println("Running UP migrations...")
	
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migration directory: %v", err)
	}

	var upFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}

	sort.Strings(upFiles)

	for _, file := range upFiles {
		fullPath := filepath.Join(migrationDir, file)
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		fmt.Printf("Running migration: %s\n", file)
		
		// Execute migration
		if sqlDB, ok := db.(interface{ Exec(string) (interface{}, error) }); ok {
			_, err = sqlDB.Exec(string(content))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v", file, err)
			}
		}
		
		fmt.Printf("Successfully ran migration: %s\n", file)
	}

	fmt.Println("All UP migrations completed successfully!")
}

func runDownMigrations(db interface{}) {
	fmt.Println("Running DOWN migrations...")
	
	migrationDir := "migrations"
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migration directory: %v", err)
	}

	var downFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".down.sql") {
			downFiles = append(downFiles, file.Name())
		}
	}

	// Sort in reverse order for down migrations
	sort.Sort(sort.Reverse(sort.StringSlice(downFiles)))

	for _, file := range downFiles {
		fullPath := filepath.Join(migrationDir, file)
		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		fmt.Printf("Running down migration: %s\n", file)
		
		// Execute migration
		if sqlDB, ok := db.(interface{ Exec(string) (interface{}, error) }); ok {
			_, err = sqlDB.Exec(string(content))
			if err != nil {
				log.Fatalf("Failed to execute down migration %s: %v", file, err)
			}
		}
		
		fmt.Printf("Successfully ran down migration: %s\n", file)
	}

	fmt.Println("All DOWN migrations completed successfully!")
}