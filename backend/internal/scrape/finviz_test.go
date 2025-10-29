package scrape

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFinvizURL_Valid(t *testing.T) {
	err := ValidateFinvizURL("https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa")
	assert.NoError(t, err)
}

func TestValidateFinvizURL_ValidWithHTTP(t *testing.T) {
	err := ValidateFinvizURL("http://finviz.com/screener.ashx?v=111")
	assert.NoError(t, err)
}

func TestValidateFinvizURL_EmptyString(t *testing.T) {
	err := ValidateFinvizURL("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "URL cannot be empty")
}

func TestValidateFinvizURL_NoScheme(t *testing.T) {
	err := ValidateFinvizURL("finviz.com/screener.ashx")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing scheme or host")
}

func TestValidateFinvizURL_WrongDomain(t *testing.T) {
	err := ValidateFinvizURL("https://google.com/search")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be from finviz.com domain")
}

func TestValidateFinvizURL_InvalidURL(t *testing.T) {
	err := ValidateFinvizURL("not-a-url")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid FINVIZ URL")
}

func TestNormalizeTickerSymbol_Uppercase(t *testing.T) {
	ticker := NormalizeTickerSymbol("aapl")
	assert.Equal(t, "AAPL", ticker)
}

func TestNormalizeTickerSymbol_Whitespace(t *testing.T) {
	ticker := NormalizeTickerSymbol("  MSFT  ")
	assert.Equal(t, "MSFT", ticker)
}

func TestNormalizeTickerSymbol_DotToDash(t *testing.T) {
	ticker := NormalizeTickerSymbol("BRK.B")
	assert.Equal(t, "BRK-B", ticker)
}

func TestNormalizeTickerSymbol_MultipleDots(t *testing.T) {
	ticker := NormalizeTickerSymbol("test.a.b")
	assert.Equal(t, "TEST-A-B", ticker)
}

func TestNormalizeTickerSymbol_AlreadyNormalized(t *testing.T) {
	ticker := NormalizeTickerSymbol("AAPL")
	assert.Equal(t, "AAPL", ticker)
}

func TestNormalizeTickerSymbol_MixedCase(t *testing.T) {
	ticker := NormalizeTickerSymbol("AaPl")
	assert.Equal(t, "AAPL", ticker)
}

func TestExtractTickerFromURL_Standard(t *testing.T) {
	ticker := extractTickerFromURL("/quote.ashx?t=AAPL")
	assert.Equal(t, "AAPL", ticker)
}

func TestExtractTickerFromURL_WithQueryParams(t *testing.T) {
	ticker := extractTickerFromURL("/quote.ashx?t=MSFT&ty=c&ta=1")
	assert.Equal(t, "MSFT", ticker)
}

func TestExtractTickerFromURL_WithDot(t *testing.T) {
	ticker := extractTickerFromURL("/quote.ashx?t=BRK.B")
	assert.Equal(t, "BRK.B", ticker)
}

func TestExtractTickerFromURL_NoTicker(t *testing.T) {
	ticker := extractTickerFromURL("/screener.ashx?v=111")
	assert.Equal(t, "", ticker)
}

func TestExtractTickerFromURL_Empty(t *testing.T) {
	ticker := extractTickerFromURL("")
	assert.Equal(t, "", ticker)
}

func TestDefaultFinvizConfig(t *testing.T) {
	config := DefaultFinvizConfig()
	assert.Equal(t, 10, config.MaxPages)
	assert.Greater(t, config.RateLimit.Seconds(), 0.0)
	assert.Greater(t, config.RequestTimeout.Seconds(), 0.0)
	assert.Equal(t, 3, config.MaxRetries)
	assert.NotEmpty(t, config.UserAgent)
	assert.Contains(t, config.UserAgent, "Mozilla")
}

func TestNewFinvizScraper(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)
	assert.NotNil(t, scraper)
	assert.NotNil(t, scraper.client)
	assert.Equal(t, config.MaxPages, scraper.config.MaxPages)
}

func TestNormalizeAndDedupe_RemovesDuplicates(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	tickers := []string{"AAPL", "MSFT", "AAPL", "NVDA", "msft"}
	result := scraper.normalizeAndDedupe(tickers)

	assert.Equal(t, 3, len(result))
	assert.Contains(t, result, "AAPL")
	assert.Contains(t, result, "MSFT")
	assert.Contains(t, result, "NVDA")
}

func TestNormalizeAndDedupe_NormalizesCase(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	tickers := []string{"aapl", "AAPL", "AaPl"}
	result := scraper.normalizeAndDedupe(tickers)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "AAPL", result[0])
}

func TestNormalizeAndDedupe_RemovesEmpty(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	tickers := []string{"AAPL", "", "MSFT", "  ", "NVDA"}
	result := scraper.normalizeAndDedupe(tickers)

	assert.Equal(t, 3, len(result))
	assert.NotContains(t, result, "")
}

func TestNormalizeAndDedupe_HandlesDots(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	tickers := []string{"BRK.B", "BRK-B", "brk.b"}
	result := scraper.normalizeAndDedupe(tickers)

	// All should normalize to BRK-B and deduplicate
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "BRK-B", result[0])
}

func TestBuildPageURL_FirstPage(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	baseURL := parseURL("https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa")
	pageURL := scraper.buildPageURL(baseURL, 1)

	assert.Equal(t, "https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa", pageURL)
}

func TestBuildPageURL_SecondPage(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	baseURL := parseURL("https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa")
	pageURL := scraper.buildPageURL(baseURL, 2)

	assert.Contains(t, pageURL, "r=21")
}

func TestBuildPageURL_ThirdPage(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	baseURL := parseURL("https://finviz.com/screener.ashx?v=111&f=ta_sma20_pa")
	pageURL := scraper.buildPageURL(baseURL, 3)

	assert.Contains(t, pageURL, "r=41")
}

func TestParseHTML_EmptyResults(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	html := `<html><body><div>No results found</div></body></html>`
	tickers, hasNext, err := scraper.parseHTML(html)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(tickers))
	assert.False(t, hasNext)
}

func TestParseHTML_WithTickers(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	html := `<html><body>
		<a href="/quote.ashx?t=AAPL">AAPL</a>
		<a href="/quote.ashx?t=MSFT">MSFT</a>
		<a href="/quote.ashx?t=NVDA">NVDA</a>
	</body></html>`

	tickers, hasNext, err := scraper.parseHTML(html)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(tickers))
	assert.Contains(t, tickers, "AAPL")
	assert.Contains(t, tickers, "MSFT")
	assert.Contains(t, tickers, "NVDA")
	assert.False(t, hasNext)
}

func TestParseHTML_WithPagination(t *testing.T) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	html := `<html><body>
		<a href="/quote.ashx?t=AAPL">AAPL</a>
		<a href="/screener.ashx?v=111&r=21">Next</a>
	</body></html>`

	tickers, hasNext, err := scraper.parseHTML(html)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(tickers))
	assert.True(t, hasNext)
}

// Helper function to parse URL for tests
func parseURL(urlStr string) *url.URL {
	u, _ := url.Parse(urlStr)
	return u
}
