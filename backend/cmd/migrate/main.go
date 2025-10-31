package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	// Get database path
	dbPath := "./trading.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	log.Printf("Migrating database: %s", dbPath)

	// Check if database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file not found: %s", dbPath)
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Read migration SQL
	migrationPath := filepath.Join("backend", "migrations", "002_add_options_columns.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	migrationSQL := string(sqlBytes)

	// Split SQL into individual statements
	statements := strings.Split(migrationSQL, ";")
	
	// Execute each statement individually, ignoring "duplicate column" errors
	log.Println("Executing migration...")
	successCount := 0
	skippedCount := 0
	
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		
		_, err := db.Exec(stmt)
		if err != nil {
			// Check if it's a "duplicate column" error - that's OK, column already exists
			errStr := strings.ToLower(err.Error())
			if strings.Contains(errStr, "duplicate column") || 
			   strings.Contains(errStr, "already exists") {
				// Extract column name from statement for logging
				columnName := extractColumnName(stmt)
				log.Printf("  ✓ Column %s already exists, skipping", columnName)
				skippedCount++
				continue
			}
			// Other errors are real problems
			log.Printf("Error executing statement: %s", stmt)
			log.Fatalf("Migration failed: %v", err)
		}
		successCount++
	}

	log.Println("✅ Migration completed successfully!")
	log.Printf("  - Added %d columns", successCount)
	if skippedCount > 0 {
		log.Printf("  - Skipped %d columns (already exist)", skippedCount)
	}
	log.Println("")
	log.Println("Your database now supports:")
	log.Println("  - Options trading metadata (instrument_type, options_strategy)")
	log.Println("  - Expiration tracking (entry_date, dte, roll_threshold_dte)")
	log.Println("  - Multi-leg options (legs_json)")
	log.Println("  - Pyramid add-on tracking (max_units, add_step_n, add_price_1/2/3)")
	log.Println("  - Breakout system selection (entry_lookback, exit_lookback)")
	log.Println("")
	log.Println("You can now create trade sessions with options!")
}

// extractColumnName tries to extract the column name from an ALTER TABLE statement
func extractColumnName(stmt string) string {
	// Look for "ADD COLUMN <name>" pattern
	parts := strings.Fields(stmt)
	for i, part := range parts {
		if strings.ToUpper(part) == "COLUMN" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown"
}
