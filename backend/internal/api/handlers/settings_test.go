package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/trading-engine/internal/storage"
)

// TestSettingsHandler_GetSettings tests the GET /api/settings endpoint
func TestSettingsHandler_GetSettings(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Initialize with default settings
	if err := db.Initialize(); err != nil {
		t.Fatalf("Failed to initialize settings: %v", err)
	}

	// Create handler
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewSettingsHandler(db, logger)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectSettings bool
	}{
		{
			name:           "GET returns settings",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectSettings: true,
		},
		{
			name:           "POST returns method not allowed",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectSettings: false,
		},
		{
			name:           "PUT returns method not allowed",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
			expectSettings: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(tt.method, "/api/settings", nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.GetSettings(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body if expecting settings
			if tt.expectSettings {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Check response structure (API returns {data: T})
				data, ok := response["data"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected data to be a map")
				}

				// Just verify we got a data object back
				if data == nil {
					t.Errorf("Expected settings data, got nil")
				}
			}
		})
	}
}

// TestSettingsHandler_GetSettings_DatabaseError tests error handling
func TestSettingsHandler_GetSettings_DatabaseError(t *testing.T) {
	// Create handler with closed database to simulate error
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	db.Close() // Close immediately to cause errors

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewSettingsHandler(db, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	w := httptest.NewRecorder()

	handler.GetSettings(w, req)

	// Should return internal server error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestSettingsHandler_UpdateSettings tests the PUT /api/settings endpoint
func TestSettingsHandler_UpdateSettings(t *testing.T) {
	// Create test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewSettingsHandler(db, logger)

	req := httptest.NewRequest(http.MethodPut, "/api/settings", nil)
	w := httptest.NewRecorder()

	handler.UpdateSettings(w, req)

	// Should return not implemented (TODO in code)
	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}
}
