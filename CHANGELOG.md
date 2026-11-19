# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0]

### Changed
- The module path is now `github.com/generator-labs/go-sdk/v2`, as Go requires for a v2 module.
  Install with `go get github.com/generator-labs/go-sdk/v2` and change the import path to match;
  the package name is still `generatorlabs`. The `v2.0.x` tags could not be installed by version.
- A `Retry-After` longer than 60 seconds (`MaxRetryAfter`) is not waited out: the 429 is returned
  as the error.
- `Get` takes a single ID. The form with several IDs is removed and returns an error: it sent an
  `ids` parameter the API does not have.

### Added
- `client.Account()`: `Get`, `Balance`, `Usage`, `PauseChecks` and `ResumeChecks`.
- `client.Cert().Findings()`: `Get`, `GetAll`, `Dismiss` and `Restore`.
- `client.Cert().Monitors().Check(id)` queues a check now, and `client.Cert().Errors().Count()`
  returns the number of current certificate errors.
- Request options on every method: `WithIdempotencyKey`, `WithRequestID`, `WithIfNoneMatch` and
  `WithFields`.
- An `Idempotency-Key` is sent automatically on requests that create something or start a check
  (creates, contact resend and confirm, certificate monitor delete, RBL check start, certificate
  check now). It is generated once per call and reused on every retry, so a retried request is not
  applied twice. A key you pass yourself is used instead.
- `Response` exposes `RequestID`, `ETag`, `Replayed` and `NotModified`, and `APIError` carries the
  `RequestID`. A 304 answer to `WithIfNoneMatch` is a normal response with `NotModified` set.
- `client.RBL().Listings()` accepts query parameters, for paging.

### Fixed
- `GetAll` returned an empty list for every resource: it read the wrong key of the response.
- The first retry now waits for `Retry-After` when the API sends it, and 1 second otherwise. It
  was sent at once.
- When the retries run out, the API's own error is returned as an `*APIError`, in place of
  "giving up after N attempt(s)".

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
