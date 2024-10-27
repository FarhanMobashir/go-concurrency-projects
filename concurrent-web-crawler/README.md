# Go Web Crawler

This repository provides a simple, concurrent web crawler built in Go, designed to fetch web pages from a seed URL while adhering to a specified depth and rate limit. It's intended for developers to learn about key web crawling concepts and techniques in Go.

## Overview

A **web crawler** is a program that systematically browses the web to collect and index information. Typically, web crawlers are the basis of search engines, aggregators, and other tools that require large-scale data from the internet.

This crawler includes the following features:

- Concurrency to crawl multiple URLs simultaneously.
- Depth control to limit the crawler’s scope.
- Rate limiting to avoid overwhelming servers.
- Domain filtering to restrict crawls to the starting domain.
- URL normalization to reduce duplicate processing.

## Concepts Covered

### Concurrency with Goroutines

Go's goroutines allow concurrent execution, letting us fetch multiple pages simultaneously. This significantly speeds up crawling.

### Rate Limiting

Using `time.Ticker`, we control the rate at which requests are sent, preventing the crawler from overloading servers.

### Depth Limiting

Depth limiting stops the crawler from venturing too far into linked pages. For example, a depth of `2` allows crawling the seed URL and its direct links only.

### URL Normalization

URLs often vary in form but point to the same resource (e.g., `/path`, `/path/`, or `http://example.com/path`). Normalizing them prevents redundant crawls.

### Domain Filtering

This ensures the crawler stays within the same domain as the seed URL, avoiding unrelated websites.

## Implementation Overview

### Basic Structure

1. **Seed URL**: Start with a single URL to initiate the crawl.
2. **Concurrency**: Use goroutines to handle URL fetching and parsing.
3. **Rate Limiting**: Control the request rate with `time.Ticker`.
4. **Data Structures**:
   - `visited` map: Keeps track of visited URLs.
   - `urlChan` channel: Manages URLs to be processed.
5. **Depth Limiting**: Define a maximum crawl depth.
6. **Domain Filtering and URL Normalization**: Clean URLs and ensure they're within the base domain.

### Key Functions

#### `fetchURL`

Fetches the HTML content from a URL while respecting the rate limit.

#### `extractLinks`

Parses the HTML content to extract all URLs on the page.

#### `normalizeAndFilterURL`

Normalizes URLs to avoid duplicates and ensures they’re within the same domain.

#### `crawl`

Handles the main crawling logic:

- Checks if the URL has already been visited.
- Fetches and parses the URL's content.
- Normalizes and filters links.
- Spawns new goroutines for each new URL within the allowed depth.

### Example Code

The following code snippet demonstrates the main components. See the `main.go` file in the repo for the full implementation.

```go
package main

import (
    "fmt"
    "net/url"
    "sync"
    "time"
    "golang.org/x/net/html"
)

// Define necessary structs and functions for crawling...
// - fetchURL
// - extractLinks
// - normalizeAndFilterURL
// - crawl

func main() {
    // Set the initial seed URL, rate limit, and max depth
    seedURL := "http://example.com"
    maxDepth := 2
    rateLimit := time.Millisecond * 200

    // Initialize the base domain, rate limiter, and goroutines
    ...
}
```

## Running the Crawler

1. Clone the repository:

   ```sh
   git clone https://github.com/your-username/go-web-crawler.git
   cd go-web-crawler
   ```

2. Run the crawler:
   ```sh
   go run main.go
   ```

## Roadmap for Further Improvements

Here are some ways to expand and improve this basic crawler:

1. **Custom Headers & User-Agent**: Mimic browser requests by adding headers.
2. **Politeness Policy**: Parse `robots.txt` to respect site permissions.
3. **Error Handling & Retries**: Add retries for transient errors like timeouts.
4. **Data Storage**: Store crawled data in a database or file for further processing.
5. **Analytics & Stats**: Track stats like crawl speed, URLs visited, and depth reached.

## Conclusion

This crawler is a powerful learning project for understanding concurrency, rate limiting, and other core concepts in web crawling. By expanding it with more advanced features, you can turn it into a robust tool for real-world data collection tasks.

Feel free to contribute and expand on this crawler!

---
