package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Client wraps HTTP client with additional functionality
type Client struct {
	httpClient *http.Client
}

// Response holds the HTTP response data
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	URL        string
}

// NewClient creates a new HTTP client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Do performs an HTTP request with environment variable substitution
func (c *Client) Do(method, url, body string, headers map[string]string) (*Response, error) {
	// Substitute environment variables in URL
	url = substituteEnvVars(url)

	// Substitute environment variables in body
	body = substituteEnvVars(body)

	// Create request
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers with environment variable substitution
	for key, value := range headers {
		req.Header.Set(key, substituteEnvVars(value))
	}

	// Set default Content-Type if not provided
	if req.Header.Get("Content-Type") == "" && body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
		URL:        url,
	}, nil
}

// substituteEnvVars replaces ${VAR_NAME} or $VAR_NAME with environment variable values
func substituteEnvVars(text string) string {
	// Replace ${VAR_NAME} format
	re1 := regexp.MustCompile(`\$\{([^}]+)\}`)
	text = re1.ReplaceAllStringFunc(text, func(match string) string {
		varName := match[2 : len(match)-1] // Remove ${ and }
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	})

	// Replace $VAR_NAME format (word characters only)
	re2 := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	text = re2.ReplaceAllStringFunc(text, func(match string) string {
		varName := match[1:] // Remove $
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	})

	return text
}
