// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

// RateLimitInfo contains rate limit information from API response headers.
type RateLimitInfo struct {
	// Limit contains the active rate limit policies, e.g. "1000;w=3600, 100;w=1"
	Limit string

	// Remaining is the number of requests remaining in the most restrictive active window
	Remaining int

	// Reset is the number of seconds until the most restrictive window resets
	Reset int
}

// Response wraps an API response with rate limit information and response metadata.
//
// The Data field contains the parsed JSON response body. Rate limit
// information from response headers is available via the RateLimit field.
type Response struct {
	// Data contains the parsed JSON response body
	Data map[string]interface{}

	// RateLimit contains rate limit information from response headers, or nil if not present
	RateLimit *RateLimitInfo

	// RequestID is the X-Request-Id of the response; quote it when contacting support
	RequestID string

	// ETag is the ETag of a GET response, to send back with WithIfNoneMatch
	ETag string

	// Replayed is true when the API answered with the result of an earlier request that
	// carried the same Idempotency-Key
	Replayed bool

	// NotModified is true for a 304 Not Modified answer to WithIfNoneMatch (Data is empty)
	NotModified bool
}
