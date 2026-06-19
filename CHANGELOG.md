# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1]

### Changed
- Error detection now reads the API `status_code` and `status_message` fields and treats any
  `status_code` (or HTTP status) of 400 or greater as an error. Errors are now returned as a typed
  `*APIError` carrying the `StatusCode` and `StatusMessage`, instead of a generic formatted error.
- `APIError` was redefined to match the v4.0 API response shape (`status_code`/`status_message`) and
  now implements the `error` interface; the previous `success`/`error.message` fields are removed.

## [2.0.0]

### Added
- Initial release with v4.0 API support
- go-retryablehttp with exponential backoff
- Automatic retries on connection errors, 5xx errors, and 429 rate limits
- `Retry-After` header support on 429 responses — backoff respects the server-specified wait time
- Exponential backoff retry logic (1s, 2s, 4s delays)
- 30 second request timeout
- RBL monitoring endpoints (hosts, profiles, sources, check, listings)
- Contact management endpoints (contacts, groups)
- Certificate monitoring endpoints (monitors, profiles, errors)
- Pagination support for list endpoints
- Configuration options for timeout and retry settings
- `Response` struct wrapping API data (`Data`) with `RateLimit` field for rate limit info
- `RateLimitInfo` struct exposing `Limit`, `Remaining`, and `Reset` from IETF draft rate limit headers
- Webhook signature verification with HMAC-SHA256 and constant-time comparison
- Credential validation (account SID and auth token format)
- Comprehensive test suite
- Comprehensive examples
- Code coverage reporting
- Security policy documentation
- Go 1.21+ support

### Changed
- Switched to plural-only endpoint naming convention
- Used Go conventions: New(), RBL(), Contact()

### Security
- Added User-Agent header for API analytics
- Implemented secure credential validation
