package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateChecklist_AllItemsChecked(t *testing.T) {
	// Perfect setup - all 6 items checked
	req := ChecklistRequest{
		Ticker:        "AAPL",
		FromPreset:    true,
		TrendPass:     true,
		LiquidityPass: true,
		TVConfirm:     true,
		EarningsOK:    true,
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, BannerGreen, result.Banner)
	assert.Equal(t, 0, result.MissingCount)
	assert.Empty(t, result.MissingItems)
	assert.True(t, result.AllowSave)
	assert.False(t, result.EvaluationTimestamp.IsZero(), "Timestamp should be recorded for GREEN")
	assert.WithinDuration(t, time.Now(), result.EvaluationTimestamp, 1*time.Second)
}

func TestEvaluateChecklist_OneItemMissing(t *testing.T) {
	// One item missing - YELLOW caution
	req := ChecklistRequest{
		Ticker:        "TSLA",
		FromPreset:    true,
		TrendPass:     true,
		LiquidityPass: true,
		TVConfirm:     true,
		EarningsOK:    false, // Missing this one
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, BannerYellow, result.Banner)
	assert.Equal(t, 1, result.MissingCount)
	assert.Len(t, result.MissingItems, 1)
	assert.Contains(t, result.MissingItems, "EarningsOK")
	assert.False(t, result.AllowSave)
	assert.True(t, result.EvaluationTimestamp.IsZero(), "Timestamp should NOT be recorded for YELLOW")
}

func TestEvaluateChecklist_TwoItemsMissing(t *testing.T) {
	// Two items missing - RED no-go
	req := ChecklistRequest{
		Ticker:        "NVDA",
		FromPreset:    true,
		TrendPass:     false, // Missing
		LiquidityPass: true,
		TVConfirm:     false, // Missing
		EarningsOK:    true,
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, BannerRed, result.Banner)
	assert.Equal(t, 2, result.MissingCount)
	assert.Len(t, result.MissingItems, 2)
	assert.Contains(t, result.MissingItems, "TrendPass")
	assert.Contains(t, result.MissingItems, "TVConfirm")
	assert.False(t, result.AllowSave)
	assert.True(t, result.EvaluationTimestamp.IsZero(), "Timestamp should NOT be recorded for RED")
}

func TestEvaluateChecklist_ThreeItemsMissing(t *testing.T) {
	// Three items missing - RED no-go
	req := ChecklistRequest{
		Ticker:        "AMD",
		FromPreset:    false, // Missing
		TrendPass:     false, // Missing
		LiquidityPass: true,
		TVConfirm:     true,
		EarningsOK:    false, // Missing
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, BannerRed, result.Banner)
	assert.Equal(t, 3, result.MissingCount)
	assert.Len(t, result.MissingItems, 3)
	assert.Contains(t, result.MissingItems, "FromPreset")
	assert.Contains(t, result.MissingItems, "TrendPass")
	assert.Contains(t, result.MissingItems, "EarningsOK")
	assert.False(t, result.AllowSave)
}

func TestEvaluateChecklist_AllItemsMissing(t *testing.T) {
	// All items missing - RED no-go
	req := ChecklistRequest{
		Ticker:        "MSFT",
		FromPreset:    false,
		TrendPass:     false,
		LiquidityPass: false,
		TVConfirm:     false,
		EarningsOK:    false,
		JournalOK:     false,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, BannerRed, result.Banner)
	assert.Equal(t, 6, result.MissingCount)
	assert.Len(t, result.MissingItems, 6)

	// Verify all items are in the missing list
	for _, itemName := range ChecklistItemNames {
		assert.Contains(t, result.MissingItems, itemName)
	}

	assert.False(t, result.AllowSave)
}

func TestEvaluateChecklist_BannerDeterminationByMissingCount(t *testing.T) {
	tests := []struct {
		name         string
		missingCount int
		expectedBanner string
		allowSave    bool
	}{
		{"0 missing", 0, BannerGreen, true},
		{"1 missing", 1, BannerYellow, false},
		{"2 missing", 2, BannerRed, false},
		{"3 missing", 3, BannerRed, false},
		{"4 missing", 4, BannerRed, false},
		{"5 missing", 5, BannerRed, false},
		{"6 missing", 6, BannerRed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with specific number of items checked
			checkedCount := 6 - tt.missingCount
			req := ChecklistRequest{
				Ticker:        "TEST",
				FromPreset:    checkedCount > 0,
				TrendPass:     checkedCount > 1,
				LiquidityPass: checkedCount > 2,
				TVConfirm:     checkedCount > 3,
				EarningsOK:    checkedCount > 4,
				JournalOK:     checkedCount > 5,
			}

			result, err := EvaluateChecklist(req)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBanner, result.Banner)
			assert.Equal(t, tt.missingCount, result.MissingCount)
			assert.Equal(t, tt.allowSave, result.AllowSave)

			if tt.expectedBanner == BannerGreen {
				assert.False(t, result.EvaluationTimestamp.IsZero())
			} else {
				assert.True(t, result.EvaluationTimestamp.IsZero())
			}
		})
	}
}

func TestEvaluateChecklist_EmptyTicker(t *testing.T) {
	req := ChecklistRequest{
		Ticker:        "",
		FromPreset:    true,
		TrendPass:     true,
		LiquidityPass: true,
		TVConfirm:     true,
		EarningsOK:    true,
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ticker is required")
}

func TestEvaluateChecklist_WhitespaceTicker(t *testing.T) {
	req := ChecklistRequest{
		Ticker:        "   ",
		FromPreset:    true,
		TrendPass:     true,
		LiquidityPass: true,
		TVConfirm:     true,
		EarningsOK:    true,
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ticker is required")
}

func TestEvaluateChecklist_MissingItemsOrder(t *testing.T) {
	// Verify that missing items are returned in consistent order
	req := ChecklistRequest{
		Ticker:        "GOOG",
		FromPreset:    false, // Missing - should be first
		TrendPass:     true,
		LiquidityPass: false, // Missing - should be third
		TVConfirm:     true,
		EarningsOK:    false, // Missing - should be fifth
		JournalOK:     true,
	}

	result, err := EvaluateChecklist(req)
	require.NoError(t, err)

	assert.Equal(t, 3, result.MissingCount)
	assert.Len(t, result.MissingItems, 3)

	// Items should appear in the order defined in ChecklistItemNames
	assert.Equal(t, "FromPreset", result.MissingItems[0])
	assert.Equal(t, "LiquidityPass", result.MissingItems[1])
	assert.Equal(t, "EarningsOK", result.MissingItems[2])
}

func TestEvaluateChecklist_DifferentTickers(t *testing.T) {
	// Verify the same checklist works for different tickers
	tickers := []string{"AAPL", "TSLA", "NVDA", "AMD", "MSFT", "GOOG"}

	for _, ticker := range tickers {
		req := ChecklistRequest{
			Ticker:        ticker,
			FromPreset:    true,
			TrendPass:     true,
			LiquidityPass: true,
			TVConfirm:     true,
			EarningsOK:    true,
			JournalOK:     true,
		}

		result, err := EvaluateChecklist(req)
		require.NoError(t, err, "ticker=%s", ticker)
		assert.Equal(t, BannerGreen, result.Banner, "ticker=%s", ticker)
	}
}
