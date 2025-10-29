package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/trading-engine/internal/storage"
)

// TestHeatHandler_CheckHeat tests the POST /api/heat/check endpoint
func TestHeatHandler_CheckHeat(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Initialize database
	if err := db.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewHeatHandler(db, logger)

	t.Run("Heat check with no existing positions", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"add_risk_dollars": 750.0,
			"add_bucket":       "Tech/Comm",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/heat/check", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CheckHeat(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check response structure (API returns {data: T})
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map")
		}

		// Check expected fields
		expectedFields := []string{"allowed", "current_portfolio_heat", "new_portfolio_heat", "current_bucket_heat", "new_bucket_heat"}
		for _, field := range expectedFields {
			if _, exists := data[field]; !exists {
				t.Errorf("Expected field '%s' in heat result", field)
			}
		}

		// With no positions, allowed should be true
		if allowed, ok := data["allowed"].(bool); !ok || !allowed {
			t.Errorf("Expected allowed=true with no positions, got '%v'", data["allowed"])
		}
	})

	t.Run("Heat check with existing positions", func(t *testing.T) {
		// Add existing position
		_, err := db.OpenPosition("NVDA")
		if err != nil {
			t.Fatalf("Failed to open test position: %v", err)
		}

		// Now try to add another $750 in Tech/Comm
		reqBody := map[string]interface{}{
			"add_risk_dollars": 750.0,
			"add_bucket":       "Tech/Comm",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/heat/check", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CheckHeat(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected data to be a map")
		}

		// Check if bucket cap is exceeded
		bucketCapExceeded, ok := data["bucket_cap_exceeded"].(bool)
		if !ok {
			t.Fatalf("Expected bucket_cap_exceeded to be a bool, got %T", data["bucket_cap_exceeded"])
		}

		// If bucket cap is exceeded, allowed should be false
		if bucketCapExceeded {
			if allowed, ok := data["allowed"].(bool); !ok || allowed {
				t.Errorf("Expected allowed=false when bucket cap exceeded, got '%v'", data["allowed"])
			}
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/heat/check", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CheckHeat(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/heat/check", nil)
		w := httptest.NewRecorder()

		handler.CheckHeat(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}
