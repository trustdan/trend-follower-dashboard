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
    "portfolioCap": 0.04,
    "bucketCap": 0.015,
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
      "entry_price": 180.50,
      "current_stop": 175.80,
      "initial_stop": 177.00,
      "shares": 159,
      "risk_dollars": 750.00,
      "bucket": "Tech/Comm",
      "status": "OPEN",
      "decision_id": 1,
      "opened_at": "2025-10-29T09:15:00Z"
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
      "id": 1,
      "ticker": "AAPL",
      "date": "2025-10-29",
      "sector": "",
      "bucket": ""
    }
  ]
}
```

### POST /api/candidates/scan

Scan FINVIZ for candidates using a preset. **(Not yet tested - requires internet)**

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
    "tickers": ["AAPL", "MSFT", "NVDA"],
    "date": "2025-10-29"
  }
}
```

### POST /api/candidates/import

Import selected candidates to database.

**Request Body:**
```json
{
  "tickers": ["AAPL", "MSFT", "NVDA"],
  "date": "2025-10-29"
}
```

**Response:**
```json
{
  "data": {
    "imported": 3,
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

---

## Correlation IDs

All responses include an `X-Correlation-ID` header for request tracking.

You can also provide your own correlation ID in the request:
```bash
curl -H "X-Correlation-ID: my-custom-id" http://localhost:8080/api/settings
```

---

## Testing

**Start server:**
```bash
./tf-engine server --listen 127.0.0.1:8080 --db trading.db
```

**Test endpoints:**
```bash
# Get settings
curl http://localhost:8080/api/settings

# Get positions
curl http://localhost:8080/api/positions

# Get candidates
curl http://localhost:8080/api/candidates

# Import candidates
curl -X POST http://localhost:8080/api/candidates/import \
  -H 'Content-Type: application/json' \
  -d '{"tickers":["AAPL","MSFT"],"date":"2025-10-29"}'
```

---

**Status:** ✅ Phase 1 Step 5 Complete
**Last Updated:** 2025-10-29
