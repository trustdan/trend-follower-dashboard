package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeTickers_Valid(t *testing.T) {
	tickers, err := NormalizeTickers("AAPL,MSFT,NVDA")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL", "MSFT", "NVDA"}, tickers)
}

func TestNormalizeTickers_WithWhitespace(t *testing.T) {
	tickers, err := NormalizeTickers("aapl, msft , NVDA")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL", "MSFT", "NVDA"}, tickers)
}

func TestNormalizeTickers_Lowercase(t *testing.T) {
	tickers, err := NormalizeTickers("aapl,msft,googl")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL", "MSFT", "GOOGL"}, tickers)
}

func TestNormalizeTickers_MixedCase(t *testing.T) {
	tickers, err := NormalizeTickers("AaPl,MsFt,NvDa")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL", "MSFT", "NVDA"}, tickers)
}

func TestNormalizeTickers_EmptyString(t *testing.T) {
	tickers, err := NormalizeTickers("")
	assert.Error(t, err)
	assert.Nil(t, tickers)
	assert.Contains(t, err.Error(), "at least one ticker required")
}

func TestNormalizeTickers_WhitespaceOnly(t *testing.T) {
	tickers, err := NormalizeTickers("   ")
	assert.Error(t, err)
	assert.Nil(t, tickers)
	assert.Contains(t, err.Error(), "at least one ticker required")
}

func TestNormalizeTickers_EmptyElements(t *testing.T) {
	tickers, err := NormalizeTickers("AAPL,,MSFT")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL", "MSFT"}, tickers)
}

func TestNormalizeTickers_AllEmptyElements(t *testing.T) {
	tickers, err := NormalizeTickers(",,")
	assert.Error(t, err)
	assert.Nil(t, tickers)
	assert.Contains(t, err.Error(), "at least one ticker required")
}

func TestNormalizeTickers_SingleTicker(t *testing.T) {
	tickers, err := NormalizeTickers("AAPL")
	assert.NoError(t, err)
	assert.Equal(t, []string{"AAPL"}, tickers)
}

func TestValidateImportRequest_Valid(t *testing.T) {
	req := ImportCandidatesRequest{
		Tickers: "AAPL,MSFT",
		Preset:  "TF_BREAKOUT_LONG",
		Sector:  "Technology",
		Bucket:  "Tech/Comm",
	}
	err := ValidateImportRequest(req)
	assert.NoError(t, err)
}

func TestValidateImportRequest_NoTickers(t *testing.T) {
	req := ImportCandidatesRequest{
		Tickers: "",
		Preset:  "TF_BREAKOUT_LONG",
	}
	err := ValidateImportRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one ticker required")
}

func TestValidateImportRequest_ValidDate(t *testing.T) {
	req := ImportCandidatesRequest{
		Tickers: "AAPL",
		Date:    "2025-10-26",
	}
	err := ValidateImportRequest(req)
	assert.NoError(t, err)
}

func TestValidateImportRequest_InvalidDate(t *testing.T) {
	req := ImportCandidatesRequest{
		Tickers: "AAPL",
		Date:    "10/26/2025",
	}
	err := ValidateImportRequest(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

func TestValidateImportRequest_NoPresetIsValid(t *testing.T) {
	req := ImportCandidatesRequest{
		Tickers: "AAPL,MSFT",
		Preset:  "",
	}
	err := ValidateImportRequest(req)
	assert.NoError(t, err)
}

func TestGetImportDate_NoDateProvided(t *testing.T) {
	date := GetImportDate("")
	assert.NotEmpty(t, date)
	// Should be today's date in YYYY-MM-DD format
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, date)
}

func TestGetImportDate_DateProvided(t *testing.T) {
	date := GetImportDate("2025-10-26")
	assert.Equal(t, "2025-10-26", date)
}

func TestImportCandidatesResult_JSONFields(t *testing.T) {
	result := ImportCandidatesResult{
		Count:      3,
		Date:       "2025-10-27",
		Tickers:    []string{"AAPL", "MSFT", "NVDA"},
		Preset:     "TF_BREAKOUT_LONG",
		Sector:     "Technology",
		Bucket:     "Tech/Comm",
		Normalized: true,
	}

	assert.Equal(t, 3, result.Count)
	assert.Equal(t, "2025-10-27", result.Date)
	assert.Equal(t, []string{"AAPL", "MSFT", "NVDA"}, result.Tickers)
	assert.Equal(t, "TF_BREAKOUT_LONG", result.Preset)
	assert.Equal(t, "Technology", result.Sector)
	assert.Equal(t, "Tech/Comm", result.Bucket)
	assert.True(t, result.Normalized)
}
