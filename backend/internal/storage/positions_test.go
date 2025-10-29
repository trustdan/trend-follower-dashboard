package storage

import (
	"os"
	"testing"
	"time"
)

// setupTestDB creates a test database
func setupTestDB(t *testing.T) *DB {
	t.Helper()
	dbPath := "test_positions_" + t.Name() + ".db"
	t.Cleanup(func() { os.Remove(dbPath) })

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.Initialize(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	return db
}

func TestOpenPosition(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a GO decision first
	decision := Decision{
		Date:         time.Now().Format("2006-01-02"),
		Ticker:       "AAPL",
		Action:       "GO",
		Entry:        180.0,
		ATR:          1.5,
		StopDistance: 3.0,
		InitialStop:  177.0,
		Shares:       25,
		RiskDollars:  75.0,
		Banner:       "GREEN",
		Method:       "stock",
	}

	decisionID, err := db.SaveDecision(decision)
	if err != nil {
		t.Fatalf("Failed to save decision: %v", err)
	}

	// Add candidate for bucket
	today := time.Now().Format("2006-01-02")
	err = db.ImportCandidates(today, []string{"AAPL"}, nil, "Technology", "Tech/Comm")
	if err != nil {
		t.Fatalf("Failed to import candidate: %v", err)
	}

	// Open position
	position, err := db.OpenPosition("AAPL")
	if err != nil {
		t.Fatalf("OpenPosition failed: %v", err)
	}

	// Verify position fields
	if position.Ticker != "AAPL" {
		t.Errorf("Expected ticker AAPL, got %s", position.Ticker)
	}
	if position.EntryPrice != 180.0 {
		t.Errorf("Expected entry 180.0, got %.2f", position.EntryPrice)
	}
	if position.CurrentStop != 177.0 {
		t.Errorf("Expected stop 177.0, got %.2f", position.CurrentStop)
	}
	if position.Shares != 25 {
		t.Errorf("Expected shares 25, got %d", position.Shares)
	}
	if position.RiskDollars != 75.0 {
		t.Errorf("Expected risk 75.0, got %.2f", position.RiskDollars)
	}
	if position.Status != "OPEN" {
		t.Errorf("Expected status OPEN, got %s", position.Status)
	}
	if position.Bucket != "Tech/Comm" {
		t.Errorf("Expected bucket Tech/Comm, got %s", position.Bucket)
	}
	if position.DecisionID != decisionID {
		t.Errorf("Expected decision_id %d, got %d", decisionID, position.DecisionID)
	}
}

func TestOpenPosition_NoDecision(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.OpenPosition("NOTHERE")
	if err == nil {
		t.Error("Expected error for ticker without decision, got nil")
	}
}

func TestOpenPosition_NoGoDecision(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a NO-GO decision
	decision := Decision{
		Date:   time.Now().Format("2006-01-02"),
		Ticker: "MSFT",
		Action: "NO-GO",
		Banner: "RED",
		Reason: "Banner is RED",
	}

	_, err := db.SaveDecision(decision)
	if err != nil {
		t.Fatalf("Failed to save decision: %v", err)
	}

	_, err = db.OpenPosition("MSFT")
	if err == nil {
		t.Error("Expected error for NO-GO decision, got nil")
	}
	if err.Error() != "cannot open position for NO-GO decision" {
		t.Errorf("Wrong error message: %v", err)
	}
}

func TestUpdateStop(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create position
	_ = createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)

	// Update stop to 179 (trailing stop up)
	err := db.UpdateStop("AAPL", 179.0)
	if err != nil {
		t.Fatalf("UpdateStop failed: %v", err)
	}

	// Verify updated
	updated, err := db.GetPositionByTicker("AAPL")
	if err != nil {
		t.Fatalf("GetPositionByTicker failed: %v", err)
	}

	if updated.CurrentStop != 179.0 {
		t.Errorf("Expected stop 179.0, got %.2f", updated.CurrentStop)
	}

	// Verify risk was recalculated
	expectedRisk := 25.0 * (180.0 - 179.0) // 25 shares * $1 = $25
	if updated.RiskDollars != expectedRisk {
		t.Errorf("Expected risk %.2f, got %.2f", expectedRisk, updated.RiskDollars)
	}
}

func TestUpdateStop_CannotMoveDown(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create position with stop at 177
	createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)

	// Try to move stop down (not allowed for LONG)
	err := db.UpdateStop("AAPL", 175.0)
	if err == nil {
		t.Error("Expected error for moving stop down, got nil")
	}
}

func TestClosePosition_Win(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create position
	createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)

	// Close with WIN
	err := db.ClosePosition("AAPL", 185.0, "WIN")
	if err != nil {
		t.Fatalf("ClosePosition failed: %v", err)
	}

	// Verify closed
	position, err := db.GetPosition(1)
	if err != nil {
		t.Fatalf("GetPosition failed: %v", err)
	}

	if position.Status != "CLOSED" {
		t.Errorf("Expected status CLOSED, got %s", position.Status)
	}
	if position.ExitPrice != 185.0 {
		t.Errorf("Expected exit 185.0, got %.2f", position.ExitPrice)
	}
	if position.Outcome != "WIN" {
		t.Errorf("Expected outcome WIN, got %s", position.Outcome)
	}

	expectedPnL := 25.0 * (185.0 - 180.0) // 25 * $5 = $125
	if position.PnL != expectedPnL {
		t.Errorf("Expected PnL %.2f, got %.2f", expectedPnL, position.PnL)
	}

	// Verify NO cooldown triggered
	cooldown, err := db.GetBucketCooldown("Tech/Comm")
	if err != nil {
		t.Fatalf("GetBucketCooldown failed: %v", err)
	}
	if cooldown != nil {
		t.Error("Expected no cooldown for WIN, but cooldown exists")
	}
}

func TestClosePosition_Loss(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create position with bucket
	position := createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)

	// Set bucket
	_, err := db.conn.Exec("UPDATE positions SET bucket = ? WHERE id = ?", "Tech/Comm", position.ID)
	if err != nil {
		t.Fatalf("Failed to set bucket: %v", err)
	}

	// Close with LOSS
	err = db.ClosePosition("AAPL", 176.0, "LOSS")
	if err != nil {
		t.Fatalf("ClosePosition failed: %v", err)
	}

	// Verify closed
	closedPos, err := db.GetPosition(position.ID)
	if err != nil {
		t.Fatalf("GetPosition failed: %v", err)
	}

	if closedPos.Status != "CLOSED" {
		t.Errorf("Expected status CLOSED, got %s", closedPos.Status)
	}
	if closedPos.Outcome != "LOSS" {
		t.Errorf("Expected outcome LOSS, got %s", closedPos.Outcome)
	}

	expectedPnL := 25.0 * (176.0 - 180.0) // 25 * -$4 = -$100
	if closedPos.PnL != expectedPnL {
		t.Errorf("Expected PnL %.2f, got %.2f", expectedPnL, closedPos.PnL)
	}

	// Verify cooldown WAS triggered
	cooldown, err := db.GetBucketCooldown("Tech/Comm")
	if err != nil {
		t.Fatalf("GetBucketCooldown failed: %v", err)
	}
	if cooldown == nil {
		t.Error("Expected cooldown for LOSS, but none found")
	}
}

func TestClosePosition_Scratch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create position
	createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)

	// Close with SCRATCH
	err := db.ClosePosition("AAPL", 180.0, "SCRATCH")
	if err != nil {
		t.Fatalf("ClosePosition failed: %v", err)
	}

	// Verify closed
	position, err := db.GetPosition(1)
	if err != nil {
		t.Fatalf("GetPosition failed: %v", err)
	}

	if position.Outcome != "SCRATCH" {
		t.Errorf("Expected outcome SCRATCH, got %s", position.Outcome)
	}
	if position.PnL != 0.0 {
		t.Errorf("Expected PnL 0.0, got %.2f", position.PnL)
	}

	// Verify NO cooldown triggered
	cooldown, err := db.GetBucketCooldown("Tech/Comm")
	if err != nil {
		t.Fatalf("GetBucketCooldown failed: %v", err)
	}
	if cooldown != nil {
		t.Error("Expected no cooldown for SCRATCH, but cooldown exists")
	}
}

func TestGetAllPositions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create multiple positions
	createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)
	createTestPosition(t, db, "MSFT", 350.0, 345.0, 5)
	createTestPosition(t, db, "GOOGL", 140.0, 137.0, 15)

	// Close one
	db.ClosePosition("MSFT", 355.0, "WIN")

	// Get all
	all, err := db.GetAllPositions("")
	if err != nil {
		t.Fatalf("GetAllPositions failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Expected 3 positions, got %d", len(all))
	}

	// Get only open
	open, err := db.GetAllPositions("OPEN")
	if err != nil {
		t.Fatalf("GetAllPositions(OPEN) failed: %v", err)
	}
	if len(open) != 2 {
		t.Errorf("Expected 2 open positions, got %d", len(open))
	}

	// Get only closed
	closed, err := db.GetAllPositions("CLOSED")
	if err != nil {
		t.Fatalf("GetAllPositions(CLOSED) failed: %v", err)
	}
	if len(closed) != 1 {
		t.Errorf("Expected 1 closed position, got %d", len(closed))
	}
}

func TestCalculatePortfolioHeat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// No positions
	heat, err := db.CalculatePortfolioHeat()
	if err != nil {
		t.Fatalf("CalculatePortfolioHeat failed: %v", err)
	}
	if heat != 0.0 {
		t.Errorf("Expected heat 0.0, got %.2f", heat)
	}

	// Create positions with known risk
	createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)  // Risk = 25 * 3 = 75
	createTestPosition(t, db, "MSFT", 350.0, 345.0, 16)  // Risk = 16 * 5 = 80
	createTestPosition(t, db, "GOOGL", 140.0, 137.0, 23) // Risk = 23 * 3 = 69

	// Set exact risk values
	db.conn.Exec("UPDATE positions SET risk_dollars = ? WHERE ticker = ?", 75.0, "AAPL")
	db.conn.Exec("UPDATE positions SET risk_dollars = ? WHERE ticker = ?", 80.0, "MSFT")
	db.conn.Exec("UPDATE positions SET risk_dollars = ? WHERE ticker = ?", 69.0, "GOOGL")

	heat, err = db.CalculatePortfolioHeat()
	if err != nil {
		t.Fatalf("CalculatePortfolioHeat failed: %v", err)
	}

	expectedHeat := 75.0 + 80.0 + 69.0 // 224
	if heat != expectedHeat {
		t.Errorf("Expected heat %.2f, got %.2f", expectedHeat, heat)
	}

	// Close one position
	db.ClosePosition("MSFT", 355.0, "WIN")

	heat, err = db.CalculatePortfolioHeat()
	if err != nil {
		t.Fatalf("CalculatePortfolioHeat after close failed: %v", err)
	}

	expectedHeat = 75.0 + 69.0 // 144 (MSFT closed)
	if heat != expectedHeat {
		t.Errorf("Expected heat %.2f after close, got %.2f", expectedHeat, heat)
	}
}

func TestCalculateBucketHeat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create positions in different buckets
	p1 := createTestPosition(t, db, "AAPL", 180.0, 177.0, 25)
	p2 := createTestPosition(t, db, "MSFT", 350.0, 345.0, 16)
	p3 := createTestPosition(t, db, "XLE", 90.0, 88.0, 35)

	// Set buckets and risk
	db.conn.Exec("UPDATE positions SET bucket = ?, risk_dollars = ? WHERE id = ?", "Tech/Comm", 75.0, p1.ID)
	db.conn.Exec("UPDATE positions SET bucket = ?, risk_dollars = ? WHERE id = ?", "Tech/Comm", 80.0, p2.ID)
	db.conn.Exec("UPDATE positions SET bucket = ?, risk_dollars = ? WHERE id = ?", "Energy", 70.0, p3.ID)

	// Calculate Tech/Comm bucket heat
	heat, err := db.CalculateBucketHeat("Tech/Comm")
	if err != nil {
		t.Fatalf("CalculateBucketHeat failed: %v", err)
	}

	expectedHeat := 75.0 + 80.0 // 155
	if heat != expectedHeat {
		t.Errorf("Expected Tech/Comm heat %.2f, got %.2f", expectedHeat, heat)
	}

	// Calculate Energy bucket heat
	heat, err = db.CalculateBucketHeat("Energy")
	if err != nil {
		t.Fatalf("CalculateBucketHeat failed: %v", err)
	}

	if heat != 70.0 {
		t.Errorf("Expected Energy heat 70.0, got %.2f", heat)
	}

	// Calculate non-existent bucket
	heat, err = db.CalculateBucketHeat("NonExistent")
	if err != nil {
		t.Fatalf("CalculateBucketHeat for empty bucket failed: %v", err)
	}
	if heat != 0.0 {
		t.Errorf("Expected heat 0.0 for empty bucket, got %.2f", heat)
	}
}

// Helper function to create a test position
func createTestPosition(t *testing.T, db *DB, ticker string, entry, stop float64, shares int) *Position {
	t.Helper()

	today := time.Now().Format("2006-01-02")

	// Create decision
	decision := Decision{
		Date:         today,
		Ticker:       ticker,
		Action:       "GO",
		Entry:        entry,
		InitialStop:  stop,
		Shares:       shares,
		RiskDollars:  float64(shares) * (entry - stop),
		Banner:       "GREEN",
		Method:       "stock",
	}

	_, err := db.SaveDecision(decision)
	if err != nil {
		t.Fatalf("Failed to save decision for %s: %v", ticker, err)
	}

	// Add candidate
	err = db.ImportCandidates(today, []string{ticker}, nil, "Technology", "Tech/Comm")
	if err != nil {
		t.Fatalf("Failed to import candidate for %s: %v", ticker, err)
	}

	// Open position
	position, err := db.OpenPosition(ticker)
	if err != nil {
		t.Fatalf("Failed to open position for %s: %v", ticker, err)
	}

	return position
}
