package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/trading-engine/internal/logx"
	"github.com/yourusername/trading-engine/internal/storage"
)

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the trading database",
		Long: `Initialize the trading database with schema and default settings.

This command creates a new SQLite database file (if it doesn't exist) and
sets up all required tables with default configuration values:
  - Equity: $10,000
  - Risk per trade: 0.75%
  - Portfolio heat cap: 4%
  - Bucket heat cap: 1.5%
  - Stop multiple (K): 2

Running this command multiple times is safe (idempotent).`,
		RunE: runInit,
	}

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	dbPath := cmd.Flag("db").Value.String()
	corrID := cmd.Flag("corr-id").Value.String()
	log := logx.WithCorrelationID(corrID)

	log.WithField("db_path", dbPath).Info("Starting database initialization")

	// Check if database already exists
	existed := false
	if _, err := os.Stat(dbPath); err == nil {
		existed = true
		log.WithField("db_path", dbPath).Info("Database file already exists")
	} else {
		log.WithField("db_path", dbPath).Info("Creating new database file")
	}

	// Open/create database
	db, err := storage.New(dbPath)
	if err != nil {
		log.WithError(err).Error("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	log.Info("Database connection established")

	// Initialize schema and settings
	if err := db.Initialize(); err != nil {
		log.WithError(err).Error("Failed to initialize database schema")
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	log.Info("Database schema initialized")

	if existed {
		fmt.Printf("✓ Database reinitialized: %s\n", dbPath)
		fmt.Println("  Existing data preserved, schema updated if needed")
		log.Info("Database reinitialized (idempotent)")
	} else {
		fmt.Printf("✓ Database created: %s\n", dbPath)
		log.Info("Database created successfully")
	}

	// Show default settings
	fmt.Println("\nDefault settings:")
	settings, err := db.GetAllSettings()
	if err != nil {
		log.WithError(err).Error("Failed to read settings")
		return fmt.Errorf("failed to read settings: %w", err)
	}

	log.WithField("settings_count", len(settings)).Info("Settings retrieved")

	for key, value := range settings {
		fmt.Printf("  %-20s = %s\n", key, value)
	}

	fmt.Println("\n✓ Trading Engine v3 ready")
	fmt.Println("  Run 'tf-engine --help' to see available commands")

	log.Info("Database initialization completed successfully")

	return nil
}
