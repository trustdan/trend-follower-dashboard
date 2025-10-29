package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCORS tests the CORS middleware
func TestCORS(t *testing.T) {
	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with CORS middleware
	handler := CORS(testHandler)

	tests := []struct {
		name           string
		method         string
		originHeader   string
		expectedOrigin string
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:           "OPTIONS preflight request",
			method:         http.MethodOptions,
			originHeader:   "http://localhost:5173",
			expectedOrigin: "http://localhost:5173",
			expectedStatus: http.StatusNoContent,
			shouldCallNext: false,
		},
		{
			name:           "GET request with origin",
			method:         http.MethodGet,
			originHeader:   "http://localhost:5173",
			expectedOrigin: "http://localhost:5173",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "POST request with origin",
			method:         http.MethodPost,
			originHeader:   "http://localhost:8080",
			expectedOrigin: "http://localhost:8080",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "Request without origin defaults to *",
			method:         http.MethodGet,
			originHeader:   "",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.originHeader != "" {
				req.Header.Set("Origin", tt.originHeader)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check CORS headers
			if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != tt.expectedOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin '%s', got '%s'", tt.expectedOrigin, origin)
			}

			methods := w.Header().Get("Access-Control-Allow-Methods")
			if methods == "" {
				t.Error("Expected Access-Control-Allow-Methods header to be set")
			}

			headers := w.Header().Get("Access-Control-Allow-Headers")
			if headers == "" {
				t.Error("Expected Access-Control-Allow-Headers header to be set")
			}

			credentials := w.Header().Get("Access-Control-Allow-Credentials")
			if credentials != "true" {
				t.Errorf("Expected Access-Control-Allow-Credentials 'true', got '%s'", credentials)
			}

			// Check if next handler was called (for non-OPTIONS requests)
			if tt.shouldCallNext {
				if w.Body.String() != "OK" {
					t.Errorf("Expected next handler to be called and return 'OK', got '%s'", w.Body.String())
				}
			} else {
				if w.Body.String() == "OK" {
					t.Error("Expected next handler NOT to be called for OPTIONS request")
				}
			}
		})
	}
}

// TestCORS_AllowsExpectedMethods tests that CORS allows the expected HTTP methods
func TestCORS_AllowsExpectedMethods(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(testHandler)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("Origin", "http://localhost:5173")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			// All methods should get CORS headers
			if origin := w.Header().Get("Access-Control-Allow-Origin"); origin == "" {
				t.Errorf("Expected CORS headers for %s method", method)
			}

			// OPTIONS should return 204, others should proceed to handler (200)
			expectedStatus := http.StatusOK
			if method == http.MethodOptions {
				expectedStatus = http.StatusNoContent
			}

			if w.Code != expectedStatus {
				t.Errorf("Expected status %d for %s, got %d", expectedStatus, method, w.Code)
			}
		})
	}
}
