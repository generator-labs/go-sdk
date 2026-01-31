// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "time"

// Config holds configuration options for the API client
type Config struct {
	// Timeout is the request timeout (default: 30 seconds)
	Timeout time.Duration

	// ConnectTimeout is the connection timeout (default: 5 seconds)
	ConnectTimeout time.Duration

	// MaxRetries is the maximum number of retry attempts (default: 3)
	MaxRetries int

	// RetryBackoff is the backoff multiplier for retries (default: 1.0)
	RetryBackoff float64

	// BaseURL is the custom API base URL
	BaseURL string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:        30 * time.Second,
		ConnectTimeout: 5 * time.Second,
		MaxRetries:     3,
		RetryBackoff:   1.0,
		BaseURL:        "https://api.generatorlabs.com/4.0/",
	}
}
