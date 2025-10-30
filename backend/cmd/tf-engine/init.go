package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/yourusername/trading-engine/internal/storage"
)

// InitCommand initializes the database with schema and default settings
func InitCommand() {
	dbPath := flag.String("db", getDefaultDBPath(), "Path to database file")
	flag.Parse()

	logger := log.New(os.Stdout, "[TF-Engine] ", log.LstdFlags)

	// Resolve to absolute path
	absPath, err := filepath.Abs(*dbPath)
	if err != nil {
		logger.Fatalf("Failed to resolve database path: %v", err)
	}

	logger.Printf("Initializing database: %s", absPath)

	// Create directory if it doesn't exist
	dbDir := filepath.Dir(absPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		logger.Fatalf("Failed to create database directory: %v", err)
	}
	logger.Printf("Database directory ready: %s", dbDir)

	// Open database connection
	db, err := storage.New(absPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize schema and default settings
	logger.Println("Creating database schema...")
	if err := db.Initialize(); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	logger.Println("Database initialized successfully!")
	logger.Println("")
	logger.Println("Next steps:")
	logger.Printf("  1. Start the server: tf-engine server --db \"%s\"", absPath)
	logger.Println("  2. Open your browser to http://localhost:8080")
	logger.Println("  3. Navigate to Settings and configure your account")
}

// getDefaultDBPath returns the default database path based on OS
func getDefaultDBPath() string {
	if runtime.GOOS == "windows" {
		// Use AppData on Windows
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "TF-Engine", "trading.db")
		}
	}

	// Fallback to current directory for non-Windows or if APPDATA not set
	return "trading.db"
}
