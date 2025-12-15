package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"next-ai-gateway/internal/config"
	"next-ai-gateway/pkg/database"
	"strings"

	"gorm.io/gorm"
)

func main() {
	// 1. Load configuration
	// We assume .env is present as per user instruction
	if err := config.LoadConfig("configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect to database
	// We need to allow multi statements for running SQL scripts
	dbConfig := config.GlobalConfig.Database
	// Append multiStatements=true to DSN if not handled by library,
	// but gorm mysql driver usually handles connection, we might need to adjust DSN manually if needed.
	// However, for safety, we will split statements by semicolon and execute one by one.

	if err := database.Init(&dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully.")

	// 3. Execute init.sql (Schema)
	if err := executeSQLFile(database.DB, "docs/init.sql"); err != nil {
		log.Fatalf("Failed to execute init.sql: %v", err)
	}
	log.Println("Schema initialized successfully.")

	// 4. Execute seed.sql (Data)
	if err := executeSQLFile(database.DB, "docs/seed.sql"); err != nil {
		log.Fatalf("Failed to execute seed.sql: %v", err)
	}
	log.Println("Seed data injected successfully.")
}

func executeSQLFile(db *gorm.DB, filepath string) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	sqlScript := string(content)

	// Basic splitting by semicolon.
	// Note: This is a simple parser and might break on semicolons inside strings.
	// For this task's known SQL content, it should be fine.
	// A more robust way is to use the mysql client, but we want a self-contained Go tool.
	statements := strings.Split(sqlScript, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("Error executing statement: %s\nError: %v", stmt, err)
			return err
		}
	}
	return nil
}
