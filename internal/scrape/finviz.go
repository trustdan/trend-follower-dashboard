package scrape

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// FinvizConfig holds configuration for FINVIZ scraping
type FinvizConfig struct {
	MaxPages       int           // Maximum pages to scrape (0 = unlimited)
	RateLimit      time.Duration // Delay between page requests
	RequestTimeout time.Duration // HTTP request timeout
	MaxRetries     int           // Maximum retry attempts
	UserAgent      string        // HTTP User-Agent header
}

// DefaultFinvizConfig returns default scraper configuration
func DefaultFinvizConfig() FinvizConfig {
	return FinvizConfig{
		MaxPages:       10,
		RateLimit:      1 * time.Second,
		RequestTimeout: 30 * time.Second,
		MaxRetries:     3,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// ScrapeResult represents the result of a FINVIZ scrape
type ScrapeResult struct {
	Tickers       []string `json:"tickers"`
	Count         int      `json:"count"`
	Date          string   `json:"date"`
	PagesScraped  int      `json:"pages_scraped"`
	MoreAvailable bool     `json:"more_available,omitempty"`
	Normalized    bool     `json:"normalized"`
}

// FinvizScraper scrapes ticker symbols from FINVIZ screener
type FinvizScraper struct {
	config FinvizConfig
	client *http.Client
}

// NewFinvizScraper creates a new FINVIZ scraper with the given config
func NewFinvizScraper(config FinvizConfig) *FinvizScraper {
	return &FinvizScraper{
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}
}

// ValidateFinvizURL validates that a URL is a valid FINVIZ screener URL
func ValidateFinvizURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid FINVIZ URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid FINVIZ URL: missing scheme or host")
	}

	if !strings.Contains(parsedURL.Host, "finviz.com") {
		return fmt.Errorf("invalid FINVIZ URL: must be from finviz.com domain")
	}

	return nil
}

// NormalizeTickerSymbol normalizes a ticker symbol according to standard conventions
func NormalizeTickerSymbol(ticker string) string {
	ticker = strings.TrimSpace(ticker)
	ticker = strings.ToUpper(ticker)

	// Replace dots with dashes (e.g., BRK.B -> BRK-B)
	ticker = strings.ReplaceAll(ticker, ".", "-")

	return ticker
}

// Scrape scrapes tickers from a FINVIZ screener URL
func (s *FinvizScraper) Scrape(queryURL string) (*ScrapeResult, error) {
	if err := ValidateFinvizURL(queryURL); err != nil {
		return nil, err
	}

	var allTickers []string
	pagesScraped := 0
	moreAvailable := false

	// Parse the base URL
	baseURL, err := url.Parse(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Scrape pages with pagination
	for page := 1; ; page++ {
		// Check max pages limit
		if s.config.MaxPages > 0 && pagesScraped >= s.config.MaxPages {
			moreAvailable = true
			break
		}

		// Build paginated URL
		pageURL := s.buildPageURL(baseURL, page)

		// Fetch page with retries
		tickers, hasNext, err := s.scrapePage(pageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scrape page %d: %w", page, err)
		}

		allTickers = append(allTickers, tickers...)
		pagesScraped++

		// Check if there are more pages
		if !hasNext {
			break
		}

		// Rate limiting between pages
		if hasNext && s.config.RateLimit > 0 {
			time.Sleep(s.config.RateLimit)
		}
	}

	// Normalize and deduplicate tickers
	normalizedTickers := s.normalizeAndDedupe(allTickers)

	result := &ScrapeResult{
		Tickers:       normalizedTickers,
		Count:         len(normalizedTickers),
		Date:          time.Now().Format("2006-01-02"),
		PagesScraped:  pagesScraped,
		MoreAvailable: moreAvailable,
		Normalized:    true,
	}

	return result, nil
}

// buildPageURL builds a paginated URL for FINVIZ
func (s *FinvizScraper) buildPageURL(baseURL *url.URL, page int) string {
	if page == 1 {
		return baseURL.String()
	}

	// FINVIZ uses r parameter for pagination (r=1, r=21, r=41, etc.)
	// Each page shows 20 results
	offset := (page - 1) * 20

	query := baseURL.Query()
	query.Set("r", fmt.Sprintf("%d", offset+1))

	newURL := *baseURL
	newURL.RawQuery = query.Encode()

	return newURL.String()
}

// scrapePage scrapes a single page and returns tickers and whether there's a next page
func (s *FinvizScraper) scrapePage(pageURL string) ([]string, bool, error) {
	var lastErr error

	// Retry logic with exponential backoff
	for attempt := 0; attempt < s.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}

		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		req.Header.Set("User-Agent", s.config.UserAgent)

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("network error: %w", err)
			continue
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = fmt.Errorf("rate limited by FINVIZ (HTTP 429), please wait before retrying")
			continue
		}

		if resp.StatusCode == http.StatusServiceUnavailable {
			lastErr = fmt.Errorf("FINVIZ service unavailable (HTTP 503), retrying...")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			continue
		}

		// Parse HTML
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		tickers, hasNext, err := s.parseHTML(string(body))
		if err != nil {
			lastErr = fmt.Errorf("HTML parsing failure: %w", err)
			continue
		}

		return tickers, hasNext, nil
	}

	return nil, false, fmt.Errorf("failed after %d retries: %w", s.config.MaxRetries, lastErr)
}

// parseHTML parses FINVIZ HTML and extracts ticker symbols
func (s *FinvizScraper) parseHTML(htmlContent string) ([]string, bool, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var tickers []string
	hasNext := false

	// Find the screener table and extract tickers
	var findTickers func(*html.Node)
	findTickers = func(n *html.Node) {
		// Look for ticker links in the screener table
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				// FINVIZ ticker links have href="/quote.ashx?t=TICKER"
				if attr.Key == "href" && strings.Contains(attr.Val, "/quote.ashx?t=") {
					ticker := extractTickerFromURL(attr.Val)
					if ticker != "" && n.FirstChild != nil {
						tickers = append(tickers, ticker)
					}
				}
			}
		}

		// Check for "next" pagination link
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && strings.Contains(attr.Val, "&r=") {
					// Check if this is a "next" link (contains higher r value)
					hasNext = true
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTickers(c)
		}
	}

	findTickers(doc)

	if len(tickers) == 0 {
		// This might not be an error - could be empty results
		return []string{}, hasNext, nil
	}

	return tickers, hasNext, nil
}

// extractTickerFromURL extracts ticker symbol from a FINVIZ URL
func extractTickerFromURL(urlStr string) string {
	re := regexp.MustCompile(`[?&]t=([A-Za-z0-9.-]+)`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// normalizeAndDedupe normalizes tickers and removes duplicates
func (s *FinvizScraper) normalizeAndDedupe(tickers []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(tickers))

	for _, ticker := range tickers {
		normalized := NormalizeTickerSymbol(ticker)
		if normalized != "" && !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}

	return result
}
