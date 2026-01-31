# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Configuration options for timeout and retry settings
- Pagination support for list endpoints
- Code coverage reporting
- Security policy documentation
- Comprehensive examples

## [2.0.0] - 2026-01-31

### Added
- Initial release with v4.0 API support
- go-retryablehttp with exponential backoff
- Automatic retries on connection errors, 5xx errors, and 429 rate limits
- Exponential backoff retry logic (1s, 2s, 4s delays)
- 30 second request timeout
- RBL monitoring endpoints (hosts, profiles, sources, check, listings)
- Contact management endpoints (contacts, groups)
- Credential validation (account SID and auth token format)
- Comprehensive test suite
- Go 1.21+ support

### Changed
- Switched to plural-only endpoint naming convention
- Used Go conventions: New(), RBL(), Contact()

### Security
- Added User-Agent header for API analytics
- Implemented secure credential validation
