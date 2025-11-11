package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WebSearchTool struct {
	apiKey string
}

type SearchInput struct {
	Query string `json:"query" jsonschema:"description=The search query"`
}

type DuckDuckGoResult struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Link    string `json:"link"`
}

func NewWebSearchTool() *WebSearchTool {
	apiKey := os.Getenv("SEARCH_API_KEY")
	return &WebSearchTool{
		apiKey: apiKey,
	}
}

func (w *WebSearchTool) Name() string {
	return "web_search"
}

func (w *WebSearchTool) Description() string {
	return "Search the web for current information, news, facts, and any information beyond your knowledge cutoff. Use this when users ask about recent events, current data, or anything you're uncertain about."
}

func (w *WebSearchTool) Fn() any {
	return func(input SearchInput) string {
		result, err := w.Execute(input)
		if err != nil {
			return fmt.Sprintf("Search failed: %v", err)
		}
		return result
	}
}

func (w *WebSearchTool) Execute(input interface{}) (string, error) {
	var searchInput SearchInput

	switch v := input.(type) {
	case SearchInput:
		searchInput = v
	case map[string]interface{}:
		if query, ok := v["query"].(string); ok {
			searchInput.Query = query
		} else {
			return "", fmt.Errorf("invalid input format: missing 'query' field")
		}
	case string:
		if err := json.Unmarshal([]byte(v), &searchInput); err != nil {
			searchInput.Query = v
		}
	default:
		jsonBytes, err := json.Marshal(input)
		if err != nil {
			return "", fmt.Errorf("cannot parse input: %v", err)
		}
		if err := json.Unmarshal(jsonBytes, &searchInput); err != nil {
			return "", fmt.Errorf("cannot parse input JSON: %v", err)
		}
	}

	if searchInput.Query == "" {
		return "", fmt.Errorf("search query is empty")
	}

	results, err := w.searchDuckDuckGo(searchInput.Query)
	if err != nil {
		return "", err
	}

	return results, nil
}

func (w *WebSearchTool) searchDuckDuckGo(query string) (string, error) {
	baseURL := "https://api.duckduckgo.com/"
	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")
	params.Add("no_html", "1")
	params.Add("skip_disambig", "1")

	fullURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to search: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Build formatted result
	var resultStr strings.Builder
	resultStr.WriteString(fmt.Sprintf("Search results for '%s':\n\n", query))

	// Abstract (summary)
	if abstract, ok := result["Abstract"].(string); ok && abstract != "" {
		resultStr.WriteString(fmt.Sprintf("Summary: %s\n", abstract))
		if abstractURL, ok := result["AbstractURL"].(string); ok && abstractURL != "" {
			resultStr.WriteString(fmt.Sprintf("Source: %s\n\n", abstractURL))
		}
	}

	// Related topics
	if relatedTopics, ok := result["RelatedTopics"].([]interface{}); ok && len(relatedTopics) > 0 {
		resultStr.WriteString("Related Information:\n")
		count := 0
		for _, topic := range relatedTopics {
			if count >= 5 { // Limit to 5 results
				break
			}
			if topicMap, ok := topic.(map[string]interface{}); ok {
				if text, ok := topicMap["Text"].(string); ok && text != "" {
					resultStr.WriteString(fmt.Sprintf("- %s\n", text))
					if firstURL, ok := topicMap["FirstURL"].(string); ok && firstURL != "" {
						resultStr.WriteString(fmt.Sprintf("  Link: %s\n", firstURL))
					}
					count++
				}
			}
		}
	}

	// If no results found
	if resultStr.Len() == len(fmt.Sprintf("Search results for '%s':\n\n", query)) {
		resultStr.WriteString("No detailed results found. This might be a very recent topic or the query needs to be more specific.\n")
		resultStr.WriteString("Try rephrasing your question or being more specific about what you're looking for.")
	}

	return resultStr.String(), nil
}

// Alternative: Use SerpAPI (requires API key but more reliable)
func (w *WebSearchTool) searchWithSerpAPI(query string) (string, error) {
	if w.apiKey == "" {
		return "", fmt.Errorf("SerpAPI key not configured. Set SEARCH_API_KEY environment variable")
	}

	baseURL := "https://serpapi.com/search"
	params := url.Values{}
	params.Add("q", query)
	params.Add("api_key", w.apiKey)
	params.Add("engine", "google")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Format results
	var resultStr strings.Builder
	resultStr.WriteString(fmt.Sprintf("Search results for '%s':\n\n", query))

	if organic, ok := result["organic_results"].([]interface{}); ok {
		for i, item := range organic {
			if i >= 5 { // Limit to 5 results
				break
			}
			if itemMap, ok := item.(map[string]interface{}); ok {
				if title, ok := itemMap["title"].(string); ok {
					resultStr.WriteString(fmt.Sprintf("%d. %s\n", i+1, title))
				}
				if snippet, ok := itemMap["snippet"].(string); ok {
					resultStr.WriteString(fmt.Sprintf("   %s\n", snippet))
				}
				if link, ok := itemMap["link"].(string); ok {
					resultStr.WriteString(fmt.Sprintf("   %s\n\n", link))
				}
			}
		}
	}

	return resultStr.String(), nil
}
