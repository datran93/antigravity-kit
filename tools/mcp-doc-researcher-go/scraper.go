package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	Title string
	Href  string
	Body  string
}

// SearchDDGLite searches DuckDuckGo Lite version and returns snippets.
func SearchDDGLite(query string, maxResults int) ([]SearchResult, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	data := url.Values{}
	data.Set("q", query)

	req, err := http.NewRequest("POST", "https://lite.duckduckgo.com/lite/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	// Parsing DuckDuckGo Lite HTML structure
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if len(results) >= maxResults {
			return
		}

		titleEl := s.Find(".result-snippet")
		if titleEl.Length() > 0 {
			// snippet is usually in the summary row
			parent := s.Prev()
			aEl := parent.Find("a.result-url")
			title := aEl.Text()
			href, _ := aEl.Attr("href")
			body := s.Find(".result-snippet").Text()

			if title != "" && href != "" {
				results = append(results, SearchResult{
					Title: strings.TrimSpace(title),
					Href:  strings.TrimSpace(href),
					Body:  strings.TrimSpace(body),
				})
			}
		}
	})

	// Sometimes the structure varies slightly, if the above doesn't yield results:
	if len(results) == 0 {
		doc.Find("td.result-snippet").Each(func(i int, s *goquery.Selection) {
			if len(results) >= maxResults {
				return
			}
			parent := s.Parent().Prev()
			aEl := parent.Find("a.result-url")
			if aEl.Length() == 0 {
				aEl = parent.Find("a")
			}
			title := aEl.Text()
			href, _ := aEl.Attr("href")
			body := s.Text()

			if title != "" && href != "" {
				results = append(results, SearchResult{
					Title: strings.TrimSpace(title),
					Href:  strings.TrimSpace(href),
					Body:  strings.TrimSpace(body),
				})
			}
		})
	}

	return results, nil
}

// FetchJinaMarkdown fetches markdown content via r.jina.ai for a given URL.
func FetchJinaMarkdown(targetUrl string, cookie string) string {
	jinaUrl := fmt.Sprintf("https://r.jina.ai/%s", targetUrl)

	client := &http.Client{
		Timeout: 25 * time.Second,
	}

	req, err := http.NewRequest("GET", jinaUrl, nil)
	if err != nil {
		return fmt.Sprintf("Error fetching markdown from %s: %v", targetUrl, err)
	}

	if cookie != "" {
		req.Header.Set("x-set-cookie", cookie)
	}

	apiKey := os.Getenv("JINA_API_KEY")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error fetching markdown from %s: %v", targetUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Error fetching markdown from %s: status %d", targetUrl, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response body from %s: %v", targetUrl, err)
	}

	return string(body)
}
