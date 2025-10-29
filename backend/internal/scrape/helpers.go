package scrape

// ScrapeFinviz is a convenience wrapper around the FinvizScraper
func ScrapeFinviz(url string) ([]string, error) {
	config := DefaultFinvizConfig()
	scraper := NewFinvizScraper(config)

	result, err := scraper.Scrape(url)
	if err != nil {
		return nil, err
	}

	return result.Tickers, nil
}
