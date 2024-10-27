package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	visited    = make(map[string]bool)  // Map to track visited URLs
	mu         sync.Mutex               // Mutex to synchronize access to the visited map
	maxDepth   = 2                      // Set max depth here
	rateLimit  = time.Millisecond * 200 // setting rate limit here (200ms)
	baseDomain *url.URL                 // Base domain to restrict crawling
)

// Fetch the page content
func fetchURL(url string, rateLimiter *time.Ticker) (string, error) {
	<-rateLimiter.C
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// normalise url and filter out external domain
func normaliseAndFilterUrl(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	// Convert relative URLs to absolute URLs
	if !parsedURL.IsAbs() {
		parsedURL = baseDomain.ResolveReference(parsedURL)
	}

	// Skip if URL is outside the base domain
	if parsedURL.Host != baseDomain.Host {
		return "", fmt.Errorf("external domain")
	}

	// Strip URL fragment
	parsedURL.Fragment = ""

	// Return the normalized URL as a string
	return parsedURL.String(), nil
}

// Parse HTML and extract all links
func extractLinks(htmlContent string) ([]string, error) {
	links := []string{}
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Helper function to recursively find <a> tags with href attributes
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}

	findLinks(doc)
	return links, nil
}

// Crawl function to fetch a URL, extract links, and add new URLs to the channel
func crawl(url string, depth int, urlChan chan string, wg *sync.WaitGroup, rateLimiter *time.Ticker) {
	defer wg.Done()

	if depth > maxDepth {
		return // Stop if max depth is reached
	}

	// Skip if the URL is already visited
	mu.Lock()
	if visited[url] {
		mu.Unlock()
		return
	}
	visited[url] = true
	mu.Unlock()

	fmt.Printf("Depth %d: Crawling %s\n", depth, url)

	// Fetch and parse the page
	content, err := fetchURL(url, rateLimiter)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	links, err := extractLinks(content)
	if err != nil {
		fmt.Println("Error extracting links:", err)
		return
	}

	// Send new URLs to the channel with increased depth
	for _, link := range links {
		normalizedLink, err := normaliseAndFilterUrl(link)
		if err != nil {
			continue // Skip links that don't pass the filter
		}

		wg.Add(1)
		go crawl(normalizedLink, depth+1, urlChan, wg, rateLimiter)

	}
}

func main() {
	urlChan := make(chan string) // Channel to hold URLs to be processed
	var wg sync.WaitGroup
	rateLimiter := time.NewTicker(rateLimit)
	defer rateLimiter.Stop()

	// Start with a seed URL
	seedURL := "https://www.google.com/"
	var err error

	// Parse the seed URL and set it as the base domain
	baseDomain, err = url.Parse(seedURL)
	if err != nil {
		fmt.Println("Invalid seed URL:", err)
		return
	}
	wg.Add(1)
	go crawl(seedURL, 0, urlChan, &wg, rateLimiter)

	// Wait for all crawl operations to complete
	wg.Wait()
	close(urlChan) // Close the channel when done
}
