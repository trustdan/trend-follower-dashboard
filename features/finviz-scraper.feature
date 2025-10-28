Feature: FINVIZ scraper
  As a trader
  I want to scrape candidate tickers from FINVIZ screener
  So that I can automatically populate my daily candidates list

  Background:
    Given I have initialized the database with default settings

  Scenario: Scrape FINVIZ with valid query URL
    When I scrape FINVIZ with a valid screener URL
    Then I should receive a list of tickers
    And the result should include a count
    And the result should include the scrape date
    And all tickers should be normalized (uppercase, no whitespace)

  Scenario: Scrape with pagination (multiple pages)
    Given FINVIZ returns results across multiple pages
    When I scrape with pagination enabled
    Then I should receive tickers from all pages
    And pagination should be detected automatically
    And the total count should match all pages combined

  Scenario: Rate limiting between pages
    Given FINVIZ returns results across 3 pages
    When I scrape with rate limiting enabled
    Then there should be at least 1 second delay between page requests
    And all pages should be fetched successfully

  Scenario: Normalize ticker symbols
    Given FINVIZ returns tickers with various formats
    When I scrape and normalize tickers
    Then "BRK.B" should become "BRK-B"
    And all tickers should be uppercase
    And duplicate tickers should be removed

  Scenario: Handle empty results
    Given FINVIZ returns no matching tickers
    When I scrape the screener
    Then I should receive an empty ticker list
    And the count should be 0
    And no error should occur

  Scenario: Handle network errors
    Given FINVIZ is unreachable
    When I attempt to scrape
    Then I should receive a clear error message
    And the error should indicate network failure

  Scenario: Handle invalid HTML response
    Given FINVIZ returns malformed HTML
    When I attempt to scrape
    Then I should receive a parsing error
    And the error should indicate HTML parsing failure

  Scenario: Handle rate limit (HTTP 429)
    Given FINVIZ returns HTTP 429 Too Many Requests
    When I attempt to scrape
    Then I should receive a rate limit error
    And the error should suggest waiting before retry

  Scenario: Scrape and auto-import to candidates
    Given I have a valid FINVIZ screener URL
    When I run "scrape-finviz --query <url> --preset TF_BREAKOUT_LONG --import"
    Then tickers should be scraped from FINVIZ
    And tickers should be automatically imported as today's candidates
    And the preset should be "TF_BREAKOUT_LONG"

  Scenario: Scrape without auto-import
    Given I have a valid FINVIZ screener URL
    When I run "scrape-finviz --query <url>"
    Then tickers should be scraped from FINVIZ
    And tickers should NOT be imported to database
    And I should receive JSON output with ticker list

  Scenario: Maximum page limit
    Given FINVIZ returns results across 20 pages
    When I scrape with a max page limit of 10
    Then only the first 10 pages should be fetched
    And a warning should indicate more pages available

  Scenario: User-Agent header
    When I scrape FINVIZ
    Then the HTTP request should include a proper User-Agent header
    And the User-Agent should identify as a legitimate browser

  Scenario: Timeout handling
    Given FINVIZ takes longer than 30 seconds to respond
    When I attempt to scrape
    Then the request should timeout
    And I should receive a timeout error

  Scenario: Retry on transient failures
    Given FINVIZ returns HTTP 503 Service Unavailable
    When I attempt to scrape with retry enabled
    Then the scraper should retry up to 3 times
    And there should be exponential backoff between retries
    And if all retries fail, I should receive a clear error

  Scenario: Validate query URL format
    When I run "scrape-finviz --query 'not-a-url'"
    Then I should receive an error "invalid FINVIZ URL"

  Scenario: Extract tickers from table rows
    Given FINVIZ returns a standard screener results page
    When I parse the HTML
    Then tickers should be extracted from the ticker column
    And company names should be ignored
    And only ticker symbols should be included in results
