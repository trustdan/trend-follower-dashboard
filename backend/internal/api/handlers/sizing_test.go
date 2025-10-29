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

// TestSizingHandler_CalculateSize tests the POST /api/sizing/calculate endpoint
func TestSizingHandler_CalculateSize(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewSizingHandler(db, logger)

	tests := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		expectedStatus int
		expectData     bool
	}{
		{
			name:   "Valid stock sizing request",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":     "AAPL",
				"entry":      180.0,
				"atr_n":      1.5,
				"method":     "stock",
				"k":          2.0,
				"risk_pct":   0.0075,
				"equity":     100000.0,
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:   "Valid option delta-atr sizing",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":     "AAPL",
				"entry":      180.0,
				"atr_n":      1.5,
				"method":     "opt-delta-atr",
				"k":          2.0,
				"risk_pct":   0.0075,
				"equity":     100000.0,
				"delta":      0.5,
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:   "Valid option max-loss sizing",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker":     "AAPL",
				"entry":      180.0,
				"atr_n":      1.5,
				"method":     "opt-maxloss",
				"k":          2.0,
				"risk_pct":   0.0075,
				"equity":     100000.0,
				"max_loss":   250.0,
			},
			expectedStatus: http.StatusOK,
			expectData:     true,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			requestBody:    nil, // Will send invalid JSON
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
		{
			name:   "Missing required fields",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"ticker": "AAPL",
				// Missing other required fields
			},
			expectedStatus: http.StatusBadRequest,
			expectData:     false,
		},
		{
			name:           "Method not allowed (GET)",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectData:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.requestBody != nil {
				bodyBytes, _ := json.Marshal(tt.requestBody)
				reqBody = bytes.NewBuffer(bodyBytes)
			} else if tt.name == "Invalid JSON" {
				reqBody = bytes.NewBufferString("{invalid json")
			} else {
				reqBody = bytes.NewBuffer([]byte{})
			}

			req := httptest.NewRequest(tt.method, "/api/sizing/calculate", reqBody)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CalculateSize(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectData {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Check response structure (API returns {data: T})
				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected data to be a map")
				}

				// Check for expected sizing result fields
				expectedFields := []string{"risk_dollars", "shares", "actual_risk"}
				for _, field := range expectedFields {
					if _, exists := data[field]; !exists {
						t.Errorf("Expected field '%s' in sizing result", field)
					}
				}

				// Verify calculations for stock method
				if tt.requestBody["method"] == "stock" {
					shares, ok := data["shares"].(float64)
					if !ok || shares != 250 {
						t.Errorf("Expected shares 250 for stock sizing (750 / 3), got %v", data["shares"])
					}
				}
			}
		})
	}
}
