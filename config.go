// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "time"

// Config holds configuration options for the API client.
//
// All fields are optional. If not provided, default values are used.
// Create a custom config to override timeouts, retry behavior, or API endpoint.
//
// Example:
//
//	config := &generatorlabs.Config{
//	    Timeout:        45 * time.Second,
//	    MaxRetries:     5,
//	    RetryBackoff:   2.0,
//	}
//	client, _ := generatorlabs.New(accountSID, authToken, config)
type Config struct {
	// Timeout is the total request timeout including retries.
	// Default: 30 seconds
	Timeout time.Duration

	// ConnectTimeout is the initial connection timeout.
	// Default: 5 seconds
	ConnectTimeout time.Duration

	// MaxRetries is the maximum number of retry attempts on failure.
	// The SDK automatically retries on connection errors, 5xx server errors,
	// and 429 rate limit errors.
	// Default: 3
	MaxRetries int

	// RetryBackoff is the exponential backoff multiplier for retries.
	// Each retry waits RetryBackoff^attemptNum seconds.
	// Default: 1.0 (results in 1s, 1s, 1s delays)
	// Set to 2.0 for exponential backoff (1s, 2s, 4s, 8s, etc.)
	RetryBackoff float64

	// BaseURL is the API base URL.
	// Default: "https://api.generatorlabs.com/4.0/"
	// Override this for testing or if using a proxy.
	BaseURL string
}

// DefaultConfig returns the default configuration values.
//
// Returns a Config with:
//   - Timeout: 30 seconds
//   - ConnectTimeout: 5 seconds
//   - MaxRetries: 3
//   - RetryBackoff: 1.0
//   - BaseURL: "https://api.generatorlabs.com/4.0/"
func DefaultConfig() *Config {
	return &Config{
		Timeout:        30 * time.Second,
		ConnectTimeout: 5 * time.Second,
		MaxRetries:     3,
		RetryBackoff:   1.0,
		BaseURL:        "https://api.generatorlabs.com/4.0/",
	}
}
