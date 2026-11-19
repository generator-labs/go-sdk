// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// RequestHandler handles HTTP requests to the Generator Labs API.
//
// This type manages HTTP client configuration, authentication, retry logic,
// and error handling for all API requests. It automatically retries failed
// requests using exponential backoff on connection errors, 5xx server errors,
// and 429 rate limit errors.
type RequestHandler struct {
	accountSID string
	authToken  string
	baseURL    string
	config     *Config
	client     *retryablehttp.Client
}

// MaxRetryAfter is the longest Retry-After (in seconds) that is waited out before a
// retry. A longer wait is not slept: the response is returned as the error instead.
const MaxRetryAfter = 60

// requestOptions holds the options of one request.
type requestOptions struct {
	idempotencyKey string
	requestID      string
	ifNoneMatch    string
	fields         string
	idempotent     bool
}

// RequestOption sets an option of one request. Options are passed as the last
// arguments of a resource method.
//
// Example:
//
//	host, err := client.RBL().Hosts().Get("HTxxxxxxxx", generatorlabs.WithFields("id", "host"))
type RequestOption func(*requestOptions)

// WithIdempotencyKey sets the Idempotency-Key of a request that creates something or
// starts a check. A repeat with the same key returns the first answer. One is generated
// when this option is left out.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *requestOptions) {
		o.idempotencyKey = key
	}
}

// WithRequestID sets the X-Request-Id to send; the API generates one otherwise.
func WithRequestID(id string) RequestOption {
	return func(o *requestOptions) {
		o.requestID = id
	}
}

// WithIfNoneMatch sends the ETag of an earlier response; the answer has NotModified set
// when nothing changed.
func WithIfNoneMatch(etag string) RequestOption {
	return func(o *requestOptions) {
		o.ifNoneMatch = etag
	}
}

// WithFields selects the keys to keep in each object of the response data (GET requests).
func WithFields(fields ...string) RequestOption {
	return func(o *requestOptions) {
		o.fields = strings.Join(fields, ",")
	}
}

// idempotent marks a request to an endpoint that accepts an Idempotency-Key, so one is
// generated when the caller gave none.
func idempotent(o *requestOptions) {
	o.idempotent = true
}

// requestArgs sorts the variadic arguments of a Get call: an id or a map of query
// parameters, followed by any request options.
func requestArgs(args []interface{}) (interface{}, map[string]interface{}, []RequestOption, error) {
	var id interface{}
	var params map[string]interface{}
	var opts []RequestOption

	for _, arg := range args {
		switch v := arg.(type) {
		case RequestOption:
			opts = append(opts, v)
		case map[string]interface{}:
			params = v
		default:
			if id != nil {
				return nil, nil, nil, fmt.Errorf("only one id can be requested at a time")
			}
			id = v
		}
	}

	return id, params, opts, nil
}

// getResource lists a resource, or retrieves one of them by id.
func (h *RequestHandler) getResource(path string, args []interface{}) (*Response, error) {
	id, params, opts, err := requestArgs(args)
	if err != nil {
		return nil, err
	}

	if id != nil {
		return h.Get(fmt.Sprintf("%s/%v", path, id), nil, opts...)
	}
	return h.Get(path, params, opts...)
}

// APIError represents a structured error response from the API.
//
// The API returns errors as a JSON body with a status_code (which mirrors the
// HTTP status code) and a status_message describing the problem.
type APIError struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
	// RequestID is the X-Request-Id of the failed request; quote it when contacting support
	RequestID string `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.StatusMessage)
}

// NewRequestHandler creates a new request handler with retry logic.
//
// This initializes an HTTP client with:
//   - Automatic retries on connection errors, 5xx errors, and 429 rate limits
//   - Exponential backoff based on config.RetryBackoff
//   - Configurable timeouts from config.Timeout and config.ConnectTimeout
//   - HTTP Basic Authentication using accountSID and authToken
//
// The handler is used internally by all API resource types.
func NewRequestHandler(accountSID, authToken string, config *Config) *RequestHandler {
	// Create retryable HTTP client
	client := retryablehttp.NewClient()
	client.RetryMax = config.MaxRetries
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = time.Duration(math.Pow(2, float64(config.MaxRetries-1))) * time.Second
	client.HTTPClient.Timeout = config.Timeout

	// Custom retry policy for 5xx and 429 errors
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on connection errors
		if err != nil {
			return true, err
		}

		// Retry on 5xx server errors and 429 rate limit, unless the Retry-After header asks
		// for a longer wait than MaxRetryAfter: that response is returned as the error
		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			seconds, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			return seconds <= MaxRetryAfter, nil
		}

		return false, nil
	}

	// Respect Retry-After header; fall back to exponential backoff. attemptNum counts the
	// requests already made less one, so it is 0 for the wait before the first retry.
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// Use Retry-After header when present (rate limit responses)
		if resp != nil {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 {
					return time.Duration(seconds) * time.Second
				}
			}
		}

		// Exponential backoff with configurable multiplier
		backoff := time.Duration(math.Pow(config.RetryBackoff, float64(attemptNum))) * time.Second
		if backoff > max {
			return max
		}
		return backoff
	}

	// When the retries run out, return the last response (or the connection error) as it
	// is, so the caller gets the API's own status code and message
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	// Disable default retry logging
	client.Logger = nil

	return &RequestHandler{
		accountSID: accountSID,
		authToken:  authToken,
		baseURL:    config.BaseURL,
		config:     config,
		client:     client,
	}
}

// send performs an HTTP request to the API with automatic retries, and returns the
// response with its body read.
//
// This method handles:
//   - Building the request URL with .json extension
//   - Adding query parameters for GET requests
//   - Adding form data for POST/PUT/DELETE requests
//   - Setting authentication and request option headers
//   - Executing the request with retry logic
func (h *RequestHandler) send(method, path string, params map[string]interface{}, opts []RequestOption) (*http.Response, []byte, error) {
	apiURL := fmt.Sprintf("%s%s.json", h.baseURL, path)

	options := &requestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Field selection applies to GET requests
	if method == "GET" && options.fields != "" {
		query := map[string]interface{}{"fields": options.fields}
		for key, value := range params {
			query[key] = value
		}
		params = query
	}

	var req *retryablehttp.Request
	var err error

	// Build request based on method
	if method == "GET" && params != nil {
		// Add query parameters
		queryParams := url.Values{}
		for key, value := range params {
			queryParams.Add(key, fmt.Sprintf("%v", value))
		}
		apiURL = fmt.Sprintf("%s?%s", apiURL, queryParams.Encode())
		req, err = retryablehttp.NewRequest(method, apiURL, nil)
	} else if (method == "POST" || method == "PUT" || method == "DELETE") && params != nil {
		// Add form data
		formData := url.Values{}
		for key, value := range params {
			switch v := value.(type) {
			case []string:
				formData.Add(key, strings.Join(v, ","))
			case []interface{}:
				parts := make([]string, len(v))
				for i, item := range v {
					parts[i] = fmt.Sprintf("%v", item)
				}
				formData.Add(key, strings.Join(parts, ","))
			default:
				formData.Add(key, fmt.Sprintf("%v", value))
			}
		}
		req, err = retryablehttp.NewRequest(method, apiURL, strings.NewReader(formData.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	} else {
		req, err = retryablehttp.NewRequest(method, apiURL, nil)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", fmt.Sprintf("GeneratorLabs-Go/%s", VERSION))
	req.Header.Set("Accept", "application/json")

	// The headers are set once per call, so every retry of the call sends the same
	// Idempotency-Key and the API answers a repeat with the first result
	if options.idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", options.idempotencyKey)
	} else if options.idempotent {
		key := make([]byte, 16)
		if _, err := rand.Read(key); err != nil {
			return nil, nil, fmt.Errorf("failed to generate an idempotency key: %w", err)
		}
		req.Header.Set("Idempotency-Key", hex.EncodeToString(key))
	}
	if options.requestID != "" {
		req.Header.Set("X-Request-Id", options.requestID)
	}
	if options.ifNoneMatch != "" {
		req.Header.Set("If-None-Match", options.ifNoneMatch)
	}

	// Set authentication
	req.SetBasicAuth(h.accountSID, h.authToken)

	// Execute request with retry logic
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, body, nil
}

// makeRequest performs an HTTP request to the API and parses the JSON response.
//
// Errors are returned if:
//   - The request cannot be created
//   - All retry attempts fail
//   - The response body cannot be read or parsed
//   - The HTTP status code or the API status_code is >= 400
func (h *RequestHandler) makeRequest(method, path string, params map[string]interface{}, opts []RequestOption) (*Response, error) {
	resp, body, err := h.send(method, path, params, opts)
	if err != nil {
		return nil, err
	}

	// Parse rate limit headers
	var rateLimit *RateLimitInfo
	if limitHeader := resp.Header.Get("RateLimit-Limit"); limitHeader != "" {
		remaining, _ := strconv.Atoi(resp.Header.Get("RateLimit-Remaining"))
		reset, _ := strconv.Atoi(resp.Header.Get("RateLimit-Reset"))
		rateLimit = &RateLimitInfo{
			Limit:     limitHeader,
			Remaining: remaining,
			Reset:     reset,
		}
	}

	response := &Response{
		RateLimit: rateLimit,
		RequestID: resp.Header.Get("X-Request-Id"),
		ETag:      resp.Header.Get("ETag"),
		Replayed:  resp.Header.Get("Idempotent-Replayed") == "true",
	}

	// 304 Not Modified has no body: the copy the caller's ETag came from is still current
	if resp.StatusCode == http.StatusNotModified {
		response.Data = map[string]interface{}{}
		response.NotModified = true
		return response, nil
	}

	// Check for empty response
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	// Parse JSON response
	if err := json.Unmarshal(body, &response.Data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Determine success vs failure from the API status_code (which mirrors the HTTP status),
	// and surface the API status_message.
	apiCode := 0
	if v, ok := response.Data["status_code"].(float64); ok {
		apiCode = int(v)
	}
	if resp.StatusCode >= 400 || apiCode >= 400 {
		code := apiCode
		if code == 0 {
			code = resp.StatusCode
		}
		msg, _ := response.Data["status_message"].(string)
		if msg == "" {
			msg = fmt.Sprintf("HTTP %d error", resp.StatusCode)
		}
		return nil, &APIError{StatusCode: code, StatusMessage: msg, RequestID: response.RequestID}
	}

	return response, nil
}

// Get performs a GET request to the API.
//
// Parameters are sent as query string parameters. The request includes
// automatic retry logic for failures.
//
// Example:
//
//	response, err := handler.Get("rbl/hosts", map[string]interface{}{"status": "active"})
func (h *RequestHandler) Get(path string, params map[string]interface{}, opts ...RequestOption) (*Response, error) {
	return h.makeRequest("GET", path, params, opts)
}

// GetPlain performs a GET request to an endpoint that answers in plain text, and
// returns the body.
//
// An error is returned when the HTTP status code is >= 400.
func (h *RequestHandler) GetPlain(path string, opts ...RequestOption) (string, error) {
	resp, body, err := h.send("GET", path, nil, opts)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", &APIError{
			StatusCode:    resp.StatusCode,
			StatusMessage: fmt.Sprintf("HTTP %d error", resp.StatusCode),
			RequestID:     resp.Header.Get("X-Request-Id"),
		}
	}

	return strings.TrimSpace(string(body)), nil
}

// Post performs a POST request to the API.
//
// Parameters are sent as application/x-www-form-urlencoded data.
// The request includes automatic retry logic for failures.
//
// Example:
//
//	response, err := handler.Post("rbl/hosts", map[string]interface{}{"name": "My Host", "host": "1.2.3.4"})
func (h *RequestHandler) Post(path string, params map[string]interface{}, opts ...RequestOption) (*Response, error) {
	return h.makeRequest("POST", path, params, opts)
}

// Put performs a PUT request to the API.
//
// Parameters are sent as application/x-www-form-urlencoded data.
// The request includes automatic retry logic for failures.
//
// Example:
//
//	response, err := handler.Put("rbl/hosts/HTxxxxxxxx", map[string]interface{}{"name": "Updated Name"})
func (h *RequestHandler) Put(path string, params map[string]interface{}, opts ...RequestOption) (*Response, error) {
	return h.makeRequest("PUT", path, params, opts)
}

// Delete performs a DELETE request to the API.
//
// No parameters are sent with DELETE requests. The request includes
// automatic retry logic for failures.
//
// Example:
//
//	response, err := handler.Delete("rbl/hosts/HTxxxxxxxx")
func (h *RequestHandler) Delete(path string, opts ...RequestOption) (*Response, error) {
	return h.makeRequest("DELETE", path, nil, opts)
}
