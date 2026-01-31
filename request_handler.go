// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// RequestHandler handles HTTP requests to the Generator Labs API
type RequestHandler struct {
	accountSID string
	authToken  string
	baseURL    string
	client     *retryablehttp.Client
}

// APIError represents an error response from the API
type APIError struct {
	Success bool `json:"success"`
	Error   struct {
		Message string `json:"message"`
	} `json:"error"`
	Message string `json:"message"`
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(accountSID, authToken, baseURL string) *RequestHandler {
	// Create retryable HTTP client
	client := retryablehttp.NewClient()
	client.RetryMax = 3 // Maximum 3 retries
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 4 * time.Second
	client.HTTPClient.Timeout = 30 * time.Second

	// Custom retry policy for 5xx and 429 errors
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on connection errors
		if err != nil {
			return true, err
		}

		// Retry on 5xx server errors and 429 rate limit
		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			return true, nil
		}

		return false, nil
	}

	// Exponential backoff: 1s, 2s, 4s
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if attemptNum == 0 {
			return 0
		}
		backoff := time.Duration(math.Pow(2, float64(attemptNum-1))) * time.Second
		if backoff > max {
			return max
		}
		return backoff
	}

	// Disable default retry logging
	client.Logger = nil

	return &RequestHandler{
		accountSID: accountSID,
		authToken:  authToken,
		baseURL:    baseURL,
		client:     client,
	}
}

// makeRequest performs an HTTP request to the API
func (h *RequestHandler) makeRequest(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("%s%s.json", h.baseURL, path)

	var req *retryablehttp.Request
	var err error

	if method == "GET" && params != nil {
		// Add query parameters
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, fmt.Sprintf("%v", value))
		}
		apiURL = fmt.Sprintf("%s?%s", apiURL, queryParams.Encode())
		req, err = retryablehttp.NewRequest(method, apiURL, nil)
	} else if (method == "POST" || method == "PUT" || method == "DELETE") && params != nil {
		// Add form parameters
		formData := url.Values{}
		for key, value := range params {
			formData.Add(key, fmt.Sprintf("%v", value))
		}
		req, err = retryablehttp.NewRequest(method, apiURL, strings.NewReader(formData.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = retryablehttp.NewRequest(method, apiURL, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", fmt.Sprintf("GeneratorLabs-Go/%s", VERSION))
	req.Header.Set("Accept", "application/json")

	// Set basic auth
	req.SetBasicAuth(h.accountSID, h.authToken)

	// Execute request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for empty response
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check for API error in response
	if success, ok := data["success"].(bool); ok && !success {
		errMsg := "Unknown error"
		if errorData, ok := data["error"].(map[string]interface{}); ok {
			if msg, ok := errorData["message"].(string); ok {
				errMsg = msg
			}
		} else if msg, ok := data["message"].(string); ok {
			errMsg = msg
		}
		return nil, fmt.Errorf("API error: %s", errMsg)
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("HTTP %d error", resp.StatusCode)
		if errorData, ok := data["error"].(map[string]interface{}); ok {
			if msg, ok := errorData["message"].(string); ok {
				errMsg = msg
			}
		} else if msg, ok := data["message"].(string); ok {
			errMsg = msg
		}
		return nil, fmt.Errorf("API error: %s", errMsg)
	}

	return data, nil
}

// Get performs a GET request
func (h *RequestHandler) Get(path string, params map[string]interface{}) (map[string]interface{}, error) {
	return h.makeRequest("GET", path, params)
}

// Post performs a POST request
func (h *RequestHandler) Post(path string, params map[string]interface{}) (map[string]interface{}, error) {
	return h.makeRequest("POST", path, params)
}

// Put performs a PUT request
func (h *RequestHandler) Put(path string, params map[string]interface{}) (map[string]interface{}, error) {
	return h.makeRequest("PUT", path, params)
}

// Delete performs a DELETE request
func (h *RequestHandler) Delete(path string) (map[string]interface{}, error) {
	return h.makeRequest("DELETE", path, nil)
}
