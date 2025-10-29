# Phase 1 - Step 5: Backend API Foundation

**Phase:** 1 - Dashboard & FINVIZ Scanner
**Step:** 5 of 9 (overall), 1 of 5 (Phase 1)
**Duration:** 2 days
**Dependencies:** Phase 0 complete (dev environment, technology decision, build pipeline)

---

## Objectives

Create the HTTP API layer that wraps the existing backend domain logic:

1. Set up HTTP server (Go standard library or Gin/Chi)
2. Implement initial API endpoints for Phase 1 features
3. Add CORS middleware for development
4. Implement comprehensive structured logging with correlation IDs
5. Create consistent JSON response format with error handling
6. Set up performance metrics logging
7. Write API integration tests
8. Document all endpoints with example requests/responses

---

## Prerequisites

- Phase 0 completed (Go 1.24+, backend tests passing)
- Backend domain logic functional (storage, domain packages)
- Understanding of REST API design principles
- Logging directory created (logs/)

---

## Step-by-Step Instructions

### 1. Choose HTTP Framework

**Decision: Use Go standard library `net/http` with `http.ServeMux`**

**Rationale:**
- No external dependencies (aligns with project simplicity)
- Standard library is mature and well-documented
- Sufficient for this project's needs
- Easy to add Chi/Gin later if needed

**Alternative:** If routing becomes complex, consider Chi (lightweight) or Gin (more features)

For now, we'll use standard library.

---

### 2. Create API Package Structure

```bash
cd /home/kali/fresh-start-trading-platform/backend

# Create API handler package
mkdir -p internal/api/handlers
mkdir -p internal/api/middleware
mkdir -p internal/api/responses

# Verify structure
tree internal/api/ -L 2
```

**Expected structure:**
```
internal/api/
├── handlers/      # HTTP handlers for each resource
│   ├── settings.go
│   ├── positions.go
│   ├── candidates.go
│   ├── checklist.go
│   ├── sizing.go
│   ├── heat.go
│   ├── gates.go
│   └── decisions.go
├── middleware/    # HTTP middleware (CORS, logging, etc.)
│   ├── cors.go
│   ├── logging.go
│   └── recovery.go
└── responses/     # Response types and helpers
    └── response.go
```

---

### 3. Create Response Helpers

Create `internal/api/responses/response.go`:

```go
// Package responses provides standard HTTP response helpers
package responses

import (
    "encoding/json"
    "log"
    "net/http"
)

// SuccessResponse wraps successful API responses
type SuccessResponse struct {
    Data interface{} `json:"data"`
}

// ErrorResponse wraps error API responses
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    int    `json:"code"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("Error encoding JSON response: %v", err)
    }
}

// Success writes a successful JSON response (200 OK)
func Success(w http.ResponseWriter, data interface{}) {
    JSON(w, http.StatusOK, SuccessResponse{Data: data})
}

// Created writes a created response (201 Created)
func Created(w http.ResponseWriter, data interface{}) {
    JSON(w, http.StatusCreated, SuccessResponse{Data: data})
}

// NoContent writes a no content response (204 No Content)
func NoContent(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNoContent)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, status int, err error) {
    response := ErrorResponse{
        Error: http.StatusText(status),
        Code:  status,
    }

    if err != nil {
        response.Message = err.Error()
    }

    JSON(w, status, response)
}

// BadRequest writes a 400 Bad Request error
func BadRequest(w http.ResponseWriter, err error) {
    Error(w, http.StatusBadRequest, err)
}

// NotFound writes a 404 Not Found error
func NotFound(w http.ResponseWriter, err error) {
    Error(w, http.StatusNotFound, err)
}

// InternalError writes a 500 Internal Server Error
func InternalError(w http.ResponseWriter, err error) {
    Error(w, http.StatusInternalServerError, err)
}
```

---

### 4. Create CORS Middleware

Create `internal/api/middleware/cors.go`:

```go
// Package middleware provides HTTP middleware
package middleware

import (
    "net/http"
)

// CORS adds CORS headers to allow cross-origin requests
// This is needed during development when frontend runs on different port
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow requests from localhost during development
        origin := r.Header.Get("Origin")
        if origin == "" {
            origin = "*"
        }

        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Correlation-ID")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

### 5. Create Logging Middleware with Correlation IDs

Create `internal/api/middleware/logging.go`:

```go
package middleware

import (
    "log"
    "net/http"
    "time"

    "github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
    // CorrelationIDKey is the context key for correlation IDs
    CorrelationIDKey contextKey = "correlationID"
)

// Logging logs HTTP requests with correlation IDs and performance metrics
func Logging(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Generate or extract correlation ID
            correlationID := r.Header.Get("X-Correlation-ID")
            if correlationID == "" {
                correlationID = uuid.New().String()
            }

            // Add correlation ID to response header
            w.Header().Set("X-Correlation-ID", correlationID)

            // Create response writer wrapper to capture status code
            rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

            // Log request
            logger.Printf("[%s] --> %s %s %s",
                correlationID,
                r.Method,
                r.URL.Path,
                r.RemoteAddr,
            )

            // Process request
            next.ServeHTTP(rw, r)

            // Log response with duration
            duration := time.Since(start)
            logger.Printf("[%s] <-- %s %s %d %s",
                correlationID,
                r.Method,
                r.URL.Path,
                rw.statusCode,
                duration,
            )

            // Log performance warning if slow
            if duration > 500*time.Millisecond {
                logger.Printf("[%s] SLOW REQUEST: %s %s took %s",
                    correlationID,
                    r.Method,
                    r.URL.Path,
                    duration,
                )
            }
        })
    }
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

**Note:** This uses `github.com/google/uuid` for correlation IDs. Add to `go.mod`:

```bash
cd /home/kali/fresh-start-trading-platform/backend
go get github.com/google/uuid
```

---

### 6. Create Recovery Middleware

Create `internal/api/middleware/recovery.go`:

```go
package middleware

import (
    "log"
    "net/http"
    "runtime/debug"
)

// Recovery recovers from panics and returns 500 Internal Server Error
func Recovery(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    // Log panic with stack trace
                    logger.Printf("PANIC: %v\n%s", err, debug.Stack())

                    // Return 500 Internal Server Error
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()

            next.ServeHTTP(w, r)
        })
    }
}
```

---

### 7. Create Settings Handler

Create `internal/api/handlers/settings.go`:

```go
// Package handlers provides HTTP handlers for API endpoints
package handlers

import (
    "log"
    "net/http"

    "github.com/fresh-start-trading-platform/backend/internal/api/responses"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

// SettingsHandler handles settings-related API requests
type SettingsHandler struct {
    db     *storage.DB
    logger *log.Logger
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(db *storage.DB, logger *log.Logger) *SettingsHandler {
    return &SettingsHandler{
        db:     db,
        logger: logger,
    }
}

// GetSettings handles GET /api/settings
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    settings, err := h.db.GetSettings()
    if err != nil {
        h.logger.Printf("Error getting settings: %v", err)
        responses.InternalError(w, err)
        return
    }

    responses.Success(w, settings)
}

// UpdateSettings handles PUT /api/settings
func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    // TODO: Implement settings update (Phase 2+)
    responses.Error(w, http.StatusNotImplemented, nil)
}
```

---

### 8. Create Positions Handler

Create `internal/api/handlers/positions.go`:

```go
package handlers

import (
    "log"
    "net/http"

    "github.com/fresh-start-trading-platform/backend/internal/api/responses"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

// PositionsHandler handles positions-related API requests
type PositionsHandler struct {
    db     *storage.DB
    logger *log.Logger
}

// NewPositionsHandler creates a new positions handler
func NewPositionsHandler(db *storage.DB, logger *log.Logger) *PositionsHandler {
    return &PositionsHandler{
        db:     db,
        logger: logger,
    }
}

// GetPositions handles GET /api/positions
func (h *PositionsHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    positions, err := h.db.GetPositions()
    if err != nil {
        h.logger.Printf("Error getting positions: %v", err)
        responses.InternalError(w, err)
        return
    }

    // Return empty array instead of null if no positions
    if positions == nil {
        positions = []storage.Position{}
    }

    responses.Success(w, positions)
}
```

---

### 9. Create Candidates Handler

Create `internal/api/handlers/candidates.go`:

```go
package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/fresh-start-trading-platform/backend/internal/api/responses"
    "github.com/fresh-start-trading-platform/backend/internal/scrape"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

// CandidatesHandler handles candidates-related API requests
type CandidatesHandler struct {
    db     *storage.DB
    logger *log.Logger
}

// NewCandidatesHandler creates a new candidates handler
func NewCandidatesHandler(db *storage.DB, logger *log.Logger) *CandidatesHandler {
    return &CandidatesHandler{
        db:     db,
        logger: logger,
    }
}

// GetCandidates handles GET /api/candidates?date=YYYY-MM-DD
func (h *CandidatesHandler) GetCandidates(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    // Get date from query params (default to today)
    dateStr := r.URL.Query().Get("date")
    if dateStr == "" {
        dateStr = time.Now().Format("2006-01-02")
    }

    candidates, err := h.db.GetCandidates(dateStr)
    if err != nil {
        h.logger.Printf("Error getting candidates for %s: %v", dateStr, err)
        responses.InternalError(w, err)
        return
    }

    // Return empty array instead of null
    if candidates == nil {
        candidates = []storage.Candidate{}
    }

    responses.Success(w, candidates)
}

// ScanRequest represents the request body for scanning FINVIZ
type ScanRequest struct {
    Preset string `json:"preset"`
}

// ScanResponse represents the response for a FINVIZ scan
type ScanResponse struct {
    Count   int      `json:"count"`
    Tickers []string `json:"tickers"`
    Date    string   `json:"date"`
}

// ScanCandidates handles POST /api/candidates/scan
func (h *CandidatesHandler) ScanCandidates(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    // Parse request body
    var req ScanRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Printf("Error decoding scan request: %v", err)
        responses.BadRequest(w, err)
        return
    }

    // Validate preset
    if req.Preset == "" {
        req.Preset = "TF_BREAKOUT_LONG" // Default preset
    }

    h.logger.Printf("Starting FINVIZ scan with preset: %s", req.Preset)

    // Get FINVIZ URL for preset (this should be configurable)
    // For now, use a hardcoded map
    presetURLs := map[string]string{
        "TF_BREAKOUT_LONG": "https://finviz.com/screener.ashx?v=111&f=ta_pattern_channelup,ta_perf_1w10o",
        // Add more presets as needed
    }

    url, ok := presetURLs[req.Preset]
    if !ok {
        h.logger.Printf("Unknown preset: %s", req.Preset)
        responses.BadRequest(w, nil)
        return
    }

    // Scrape FINVIZ
    tickers, err := scrape.ScrapeFinviz(url)
    if err != nil {
        h.logger.Printf("Error scraping FINVIZ: %v", err)
        responses.InternalError(w, err)
        return
    }

    h.logger.Printf("FINVIZ scan found %d tickers", len(tickers))

    // Return results (don't save yet, user will review and import)
    result := ScanResponse{
        Count:   len(tickers),
        Tickers: tickers,
        Date:    time.Now().Format("2006-01-02"),
    }

    responses.Success(w, result)
}

// ImportRequest represents the request body for importing candidates
type ImportRequest struct {
    Tickers []string `json:"tickers"`
    Date    string   `json:"date"`
}

// ImportCandidates handles POST /api/candidates/import
func (h *CandidatesHandler) ImportCandidates(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    // Parse request body
    var req ImportRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Printf("Error decoding import request: %v", err)
        responses.BadRequest(w, err)
        return
    }

    // Validate
    if len(req.Tickers) == 0 {
        h.logger.Printf("No tickers provided for import")
        responses.BadRequest(w, nil)
        return
    }

    if req.Date == "" {
        req.Date = time.Now().Format("2006-01-02")
    }

    h.logger.Printf("Importing %d candidates for %s", len(req.Tickers), req.Date)

    // Save to database
    if err := h.db.AddCandidates(req.Tickers, req.Date); err != nil {
        h.logger.Printf("Error importing candidates: %v", err)
        responses.InternalError(w, err)
        return
    }

    responses.Success(w, map[string]interface{}{
        "imported": len(req.Tickers),
        "date":     req.Date,
    })
}

// DeleteCandidate handles DELETE /api/candidates/:ticker
func (h *CandidatesHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        responses.Error(w, http.StatusMethodNotAllowed, nil)
        return
    }

    // TODO: Implement candidate deletion (Phase 1 Step 9)
    responses.Error(w, http.StatusNotImplemented, nil)
}
```

---

### 10. Create HTTP Server

Create `backend/cmd/tf-engine/server.go` (or add to existing main.go):

```go
package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/fresh-start-trading-platform/backend/internal/api/handlers"
    "github.com/fresh-start-trading-platform/backend/internal/api/middleware"
    "github.com/fresh-start-trading-platform/backend/internal/logx"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
    "github.com/fresh-start-trading-platform/backend/internal/webui"
)

// ServerCommand runs the HTTP server
func ServerCommand() {
    // Parse flags
    listen := flag.String("listen", "127.0.0.1:8080", "Address to listen on")
    dbPath := flag.String("db", "trading.db", "Path to database file")
    flag.Parse()

    // Initialize logger
    logCfg := logx.DefaultConfig()
    logger, err := logx.NewLogger(logCfg)
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }

    logger.Println("Starting TF-Engine HTTP Server...")

    // Initialize database
    absDBPath, _ := filepath.Abs(*dbPath)
    logger.Printf("Opening database: %s", absDBPath)

    db, err := storage.NewDB(*dbPath)
    if err != nil {
        logger.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Initialize handlers
    settingsHandler := handlers.NewSettingsHandler(db, logger)
    positionsHandler := handlers.NewPositionsHandler(db, logger)
    candidatesHandler := handlers.NewCandidatesHandler(db, logger)

    // Create router
    mux := http.NewServeMux()

    // API routes
    mux.HandleFunc("/api/settings", settingsHandler.GetSettings)
    mux.HandleFunc("/api/positions", positionsHandler.GetPositions)
    mux.HandleFunc("/api/candidates", candidatesHandler.GetCandidates)
    mux.HandleFunc("/api/candidates/scan", candidatesHandler.ScanCandidates)
    mux.HandleFunc("/api/candidates/import", candidatesHandler.ImportCandidates)

    // Serve embedded Svelte UI
    sfs, err := webui.Sub()
    if err != nil {
        logger.Printf("Warning: Could not load embedded UI: %v", err)
        logger.Println("API endpoints will still work")
    } else {
        mux.Handle("/", http.FileServer(http.FS(sfs)))
        logger.Println("Embedded UI loaded successfully")
    }

    // Apply middleware
    handler := middleware.Recovery(logger)(
        middleware.Logging(logger)(
            middleware.CORS(mux),
        ),
    )

    // Create server
    srv := &http.Server{
        Addr:         *listen,
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        logger.Printf("Server listening on http://%s", *listen)
        logger.Println("Press Ctrl+C to stop")

        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Printf("Server forced to shutdown: %v", err)
    }

    logger.Println("Server stopped")
}
```

**Update `cmd/tf-engine/main.go` to add server command:**

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: tf-engine <command> [options]")
        fmt.Println("Commands: init, settings, size, checklist, heat, server")
        os.Exit(1)
    }

    command := os.Args[1]
    os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
    flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

    switch command {
    case "server":
        ServerCommand()
    // ... other commands (init, settings, etc.)
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}
```

---

### 11. Test the API Server

```bash
cd /home/kali/fresh-start-trading-platform/backend

# Ensure database exists
ls -la trading.db
# If not, initialize it:
# go run cmd/tf-engine/main.go init

# Build and run server
go run cmd/tf-engine/main.go server --listen 127.0.0.1:8080

# Expected output:
# Starting TF-Engine HTTP Server...
# Opening database: /home/kali/fresh-start-trading-platform/backend/trading.db
# Embedded UI loaded successfully
# Server listening on http://127.0.0.1:8080
# Press Ctrl+C to stop
```

**Test API endpoints (in another terminal):**

```bash
# Test GET /api/settings
curl http://localhost:8080/api/settings

# Expected: {"data":{"equity":100000,"riskPct":0.75,...}}

# Test GET /api/positions
curl http://localhost:8080/api/positions

# Expected: {"data":[]} or list of positions

# Test GET /api/candidates
curl http://localhost:8080/api/candidates

# Expected: {"data":[]} or list of candidates

# Test POST /api/candidates/scan (FINVIZ scan)
curl -X POST http://localhost:8080/api/candidates/scan \
  -H "Content-Type: application/json" \
  -d '{"preset":"TF_BREAKOUT_LONG"}'

# Expected: {"data":{"count":X,"tickers":[...],"date":"2025-10-29"}}

# Test POST /api/candidates/import
curl -X POST http://localhost:8080/api/candidates/import \
  -H "Content-Type: application/json" \
  -d '{"tickers":["AAPL","MSFT"],"date":"2025-10-29"}'

# Expected: {"data":{"imported":2,"date":"2025-10-29"}}

# Verify import worked
curl http://localhost:8080/api/candidates?date=2025-10-29

# Expected: {"data":[{"ticker":"AAPL",...},{"ticker":"MSFT",...}]}
```

---

### 12. Write Integration Tests

Create `backend/internal/api/handlers/handlers_test.go`:

```go
package handlers_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/fresh-start-trading-platform/backend/internal/api/handlers"
    "github.com/fresh-start-trading-platform/backend/internal/storage"
)

func TestSettingsHandler_GetSettings(t *testing.T) {
    // Create in-memory database for testing
    db, err := storage.NewDB(":memory:")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    defer db.Close()

    // Initialize database
    if err := db.InitDB(); err != nil {
        t.Fatalf("Failed to initialize database: %v", err)
    }

    // Create handler
    handler := handlers.NewSettingsHandler(db, nil)

    // Create request
    req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
    w := httptest.NewRecorder()

    // Execute
    handler.GetSettings(w, req)

    // Assert
    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    var response map[string]interface{}
    if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    data, ok := response["data"]
    if !ok {
        t.Fatal("Response missing 'data' field")
    }

    settings, ok := data.(map[string]interface{})
    if !ok {
        t.Fatal("Data is not a settings object")
    }

    if settings["equity"] == nil {
        t.Error("Settings missing 'equity' field")
    }
}

// Add more tests for other handlers...
```

Run tests:

```bash
cd /home/kali/fresh-start-trading-platform/backend
go test ./internal/api/handlers/... -v
```

---

### 13. Document API Endpoints

Create `docs/api-reference.md`:

```markdown
# API Reference

**Base URL:** `http://localhost:8080/api`
**Content-Type:** `application/json`

---

## Settings

### GET /api/settings

Get current account settings.

**Response:**
```json
{
  "data": {
    "equity": 100000.00,
    "riskPct": 0.75,
    "portfolioCap": 4.00,
    "bucketCap": 1.50,
    "maxUnits": 4
  }
}
```

---

## Positions

### GET /api/positions

Get all open positions.

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "ticker": "AAPL",
      "bucket": "Tech/Comm",
      "entry": 180.50,
      "shares": 159,
      "riskPerUnit": 750.00,
      "currentStop": 175.80,
      "entryDate": "2025-10-29"
    }
  ]
}
```

---

## Candidates

### GET /api/candidates?date=YYYY-MM-DD

Get candidates for a specific date.

**Query Parameters:**
- `date` (optional): Date in YYYY-MM-DD format. Defaults to today.

**Response:**
```json
{
  "data": [
    {
      "ticker": "AAPL",
      "date": "2025-10-29",
      "importedAt": "2025-10-29T09:15:00Z"
    }
  ]
}
```

### POST /api/candidates/scan

Scan FINVIZ for candidates using a preset.

**Request Body:**
```json
{
  "preset": "TF_BREAKOUT_LONG"
}
```

**Response:**
```json
{
  "data": {
    "count": 23,
    "tickers": ["AAPL", "MSFT", "NVDA", ...],
    "date": "2025-10-29"
  }
}
```

### POST /api/candidates/import

Import selected candidates to database.

**Request Body:**
```json
{
  "tickers": ["AAPL", "MSFT"],
  "date": "2025-10-29"
}
```

**Response:**
```json
{
  "data": {
    "imported": 2,
    "date": "2025-10-29"
  }
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Bad Request",
  "message": "Detailed error message",
  "code": 400
}
```

**Status Codes:**
- `200` OK - Success
- `201` Created - Resource created
- `204` No Content - Success with no response body
- `400` Bad Request - Invalid request
- `404` Not Found - Resource not found
- `405` Method Not Allowed - HTTP method not supported
- `500` Internal Server Error - Server error
- `501` Not Implemented - Feature not implemented yet
```

---

## Verification Checklist

Before proceeding to Step 6, verify:

- [ ] API package structure created (`internal/api/handlers`, `middleware`, `responses`)
- [ ] Response helpers created (`responses/response.go`)
- [ ] CORS middleware implemented
- [ ] Logging middleware with correlation IDs implemented
- [ ] Recovery middleware implemented
- [ ] Settings handler created and tested
- [ ] Positions handler created and tested
- [ ] Candidates handler created with scan and import
- [ ] HTTP server created with graceful shutdown
- [ ] Server runs and listens on port 8080
- [ ] All API endpoints respond correctly (curl tests pass)
- [ ] Correlation IDs appear in logs and response headers
- [ ] Performance metrics logged (request duration)
- [ ] Slow requests logged (>500ms)
- [ ] Integration tests written and passing
- [ ] API documentation created (`docs/api-reference.md`)
- [ ] `go.mod` updated with new dependencies (uuid)

---

## Expected Outputs

After completing this step, you should have:

1. **API Package:**
   - `internal/api/handlers/` - HTTP handlers
   - `internal/api/middleware/` - Middleware functions
   - `internal/api/responses/` - Response helpers

2. **Working HTTP Server:**
   - `cmd/tf-engine/server.go` - Server command
   - Listens on `127.0.0.1:8080`
   - Serves API endpoints and embedded UI

3. **API Endpoints (Phase 1):**
   - `GET /api/settings` - Get settings
   - `GET /api/positions` - Get positions
   - `GET /api/candidates?date=...` - Get candidates
   - `POST /api/candidates/scan` - Scan FINVIZ
   - `POST /api/candidates/import` - Import candidates

4. **Logging:**
   - Correlation IDs in all requests/responses
   - Performance metrics (request duration)
   - Slow request warnings (>500ms)
   - All logs in `logs/tf-engine.log`

5. **Tests:**
   - `internal/api/handlers/handlers_test.go` - Integration tests

6. **Documentation:**
   - `docs/api-reference.md` - API documentation

---

## Troubleshooting

### Server won't start

**Problem:** Port 8080 already in use
**Solution:**
```bash
# Check what's using port 8080
lsof -i :8080

# Use different port
go run cmd/tf-engine/main.go server --listen 127.0.0.1:8090
```

### Database errors

**Problem:** Cannot open database
**Solution:**
```bash
# Ensure database exists
ls -la trading.db

# Initialize if missing
go run cmd/tf-engine/main.go init
```

### FINVIZ scan fails

**Problem:** Network errors or empty results
**Solution:**
- Check internet connection
- Verify FINVIZ URL is correct
- Check if FINVIZ changed their page structure
- Test scraper directly: `go run cmd/tf-engine/main.go import-candidates --preset TF_BREAKOUT_LONG`

### CORS errors in browser

**Problem:** Browser blocks requests
**Solution:**
- Ensure CORS middleware is applied
- Check `Access-Control-Allow-Origin` header in response
- For development, middleware allows all origins

---

## Time Estimate

- **Package Structure:** 30 minutes
- **Response Helpers:** 30 minutes
- **Middleware (CORS, Logging, Recovery):** 1-2 hours
- **Handlers (Settings, Positions, Candidates):** 2-3 hours
- **HTTP Server:** 1-2 hours
- **Testing:** 1-2 hours
- **Documentation:** 1 hour

**Total:** ~7-10 hours (1-2 days with testing and troubleshooting)

---

## References

- [Go net/http Package](https://pkg.go.dev/net/http)
- [Go embed Package](https://pkg.go.dev/embed)
- [UUID Package](https://pkg.go.dev/github.com/google/uuid)
- [roadmap.md - Step 5 Details](../plans/roadmap.md#step-5-backend-api-foundation)
- [overview-plan.md - Backend API Endpoints](../plans/overview-plan.md#backend-api-endpoints-to-implement)

---

## Next Step

Proceed to: **[Phase 1 - Step 6: Application Layout & Navigation](phase1-step6-layout-navigation.md)**

---

**Status:** 📋 Ready for Execution
**Created:** 2025-10-29
**Last Updated:** 2025-10-29
