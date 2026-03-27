# Generator Labs Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/generator-labs/go-sdk.svg)](https://pkg.go.dev/github.com/generator-labs/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/generator-labs/go-sdk)](https://goreportcard.com/report/github.com/generator-labs/go-sdk)
[![Tests](https://github.com/generator-labs/go-sdk/workflows/Tests/badge.svg)](https://github.com/generator-labs/go-sdk/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [Generator Labs API](https://generatorlabs.com). This library provides a simple and intuitive interface for interacting with the Generator Labs v4.0 API, including RBL monitoring, contact management, and more.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Webhook Verification](#webhook-verification)
- [API Reference](#api-reference)
  - [Client Initialization](#client-initialization)
  - [Configuration Options](#configuration-options)
  - [Pagination](#pagination)
  - [RBL Monitoring](#rbl-monitoring) - [Hosts](#hosts) | [Profiles](#profiles) | [Sources](#sources) | [Check & Listings](#check--listings)
  - [Certificate Monitoring](#certificate-monitoring) - [Errors](#errors) | [Monitors](#monitors) | [Profiles](#profiles-1)
  - [Contact Management](#contact-management) - [Contacts](#contacts) | [Groups](#groups)
- [Error Handling](#error-handling)
- [Retry Logic](#retry-logic)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)
- [Requirements](#requirements)
- [Testing](#testing)
- [Security](#security)
- [Release History](#release-history)
- [License](#license)
- [Support](#support)
- [Contributing](#contributing)

## Features

- Full support for Generator Labs API v4.0
- Configurable timeouts, retries, and backoff strategies
- Automatic retry logic with exponential backoff
- Automatic pagination for large result sets
- Connection pooling and timeout management
- Type-safe API methods
- Comprehensive error handling
- Go 1.21+ support

## Installation

```bash
go get github.com/generator-labs/go-sdk
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"

    generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
    // Initialize the client with default config
    client, err := generatorlabs.New(
        "YOUR_ACCOUNT_SID",
        "YOUR_AUTH_TOKEN",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Or with custom configuration
    config := &generatorlabs.Config{
        Timeout:        45 * time.Second,
        ConnectTimeout: 10 * time.Second,
        MaxRetries:     5,
        RetryBackoff:   2.0,
        BaseURL:        "https://api.generatorlabs.com/4.0/",
    }
    client, err = generatorlabs.New(
        "YOUR_ACCOUNT_SID",
        "YOUR_AUTH_TOKEN",
        config,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Get all monitored hosts
    hosts, err := client.RBL().Hosts().Get()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hosts: %+v\n", hosts)

    // Start a manual RBL check
    result, err := client.RBL().Check().Start(map[string]interface{}{
        "host": "8.8.8.8",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Check result: %+v\n", result)

    // Get all contacts with automatic pagination
    allContacts, err := client.Contact().Contacts().GetAll(nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total contacts: %d\n", len(allContacts))
}
```

## Webhook Verification

The SDK includes a helper for verifying incoming webhook signatures. Each webhook is assigned a signing secret (available in the Portal), which is used to compute an HMAC-SHA256 signature sent with every request in the `X-Webhook-Signature` header.

```go
import generatorlabs "github.com/generator-labs/go-sdk"

header := r.Header.Get("X-Webhook-Signature")
body, _ := io.ReadAll(r.Body)
secret := os.Getenv("GENERATOR_LABS_WEBHOOK_SECRET")

payload, err := generatorlabs.VerifyWebhook(string(body), header, secret, generatorlabs.DefaultWebhookTolerance)
if err != nil {
    // Signature verification failed
    http.Error(w, `{"error":"Invalid signature"}`, http.StatusForbidden)
    return
}

// payload is the decoded event data
fmt.Println(payload["event"])
```

The default timestamp tolerance is 5 minutes. You can customize it (in seconds), or pass `0` to disable:

```go
payload, err := generatorlabs.VerifyWebhook(string(body), header, secret, 600)  // 10-minute tolerance
payload, err := generatorlabs.VerifyWebhook(string(body), header, secret, 0)    // disable timestamp check
```

See `examples/webhook_verification/main.go` for a complete example.

## API Reference

### Client Initialization

```go
// With default configuration
client, err := generatorlabs.New(accountSID, authToken)

// With custom configuration
config := &generatorlabs.Config{
    Timeout:        45 * time.Second,  // Request timeout (default: 30s)
    ConnectTimeout: 10 * time.Second,  // Connection timeout (default: 5s)
    MaxRetries:     5,                 // Max retry attempts (default: 3)
    RetryBackoff:   2.0,              // Backoff multiplier (default: 1.0)
    BaseURL:        "https://api.generatorlabs.com/4.0/",
}
client, err := generatorlabs.New(accountSID, authToken, config)
```

### Configuration Options

- **Timeout**: Maximum duration for the entire request (default: 30 seconds)
- **ConnectTimeout**: Maximum duration for connection establishment (default: 5 seconds)
- **MaxRetries**: Maximum number of retry attempts for failed requests (default: 3)
- **RetryBackoff**: Multiplier for exponential backoff between retries (default: 1.0)
- **BaseURL**: Custom API base URL (default: https://api.generatorlabs.com/4.0/)

### Pagination

All list operations support automatic pagination using the `GetAll()` method:

```go
// Get all hosts across multiple pages
allHosts, err := client.RBL().Hosts().GetAll(map[string]interface{}{
    "page_size": 50,  // Items per page
})

// Get all profiles with automatic pagination
allProfiles, err := client.RBL().Profiles().GetAll(nil)

// Get all contacts with automatic pagination
allContacts, err := client.Contact().Contacts().GetAll(nil)

// Get all groups with automatic pagination
allGroups, err := client.Contact().Groups().GetAll(nil)
```

### RBL Monitoring

#### Hosts

```go
// Get all hosts
hosts, err := client.RBL().Hosts().Get()

// Get a specific host
host, err := client.RBL().Hosts().Get("HT1a2b3c4d5e6f7890abcdef1234567890")

// Get multiple hosts
hosts, err := client.RBL().Hosts().Get("HT1a2b3c4d5e6f7890abcdef1234567890", "HT2b3c4d5e6f7890abcdef12345678901a")

// Create a host
params := map[string]interface{}{
    "host": "8.8.8.8",
    "name": "Google DNS",
    "profile": "RP9f8e7d6c5b4a3210fedcba0987654321",
    "contact_group": []string{
        "CG4f3e2d1c0b9a8776655443322110fedc",
        "CG5a6b7c8d9e0f1234567890abcdef1234",
    },
    "tags": []string{"production", "web"},
}
host, err := client.RBL().Hosts().Create(params)

// Update a host
params := map[string]interface{}{
    "name": "Updated description",
    "tags": []string{"production", "web"},
}
host, err := client.RBL().Hosts().Update("HT1a2b3c4d5e6f7890abcdef1234567890", params)

// Delete a host
result, err := client.RBL().Hosts().Delete("HT1a2b3c4d5e6f7890abcdef1234567890")
```

#### Profiles

```go
// Get all profiles
profiles, err := client.RBL().Profiles().Get()

// Get a specific profile
profile, err := client.RBL().Profiles().Get("RP9f8e7d6c5b4a3210fedcba0987654321")

// Create a profile
params := map[string]interface{}{
    "name": "My Custom Profile",
    "entries": []string{
        "RB1234567890abcdef1234567890abcdef",
        "RB0987654321fedcba0987654321fedcba",
    },
}
profile, err := client.RBL().Profiles().Create(params)

// Update a profile
params := map[string]interface{}{
    "name": "Updated Profile Name",
    "entries": []string{
        "RB1234567890abcdef1234567890abcdef",
        "RB0987654321fedcba0987654321fedcba",
    },
}
profile, err := client.RBL().Profiles().Update("RP9f8e7d6c5b4a3210fedcba0987654321", params)

// Delete a profile
result, err := client.RBL().Profiles().Delete("RP9f8e7d6c5b4a3210fedcba0987654321")
```

#### Sources

```go
// Get all RBL sources
sources, err := client.RBL().Sources().Get()

// Get a specific source
source, err := client.RBL().Sources().Get("RB18c470cc518a09678bb280960dbdd524")

// Create a custom source
params := map[string]interface{}{
    "host": "custom.rbl.example.com",
    "type": "rbl",
    "custom_codes": []string{"127.0.0.2", "127.0.0.3"},
}
source, err := client.RBL().Sources().Create(params)

// Update a source
params := map[string]interface{}{
    "host": "updated.rbl.example.com",
    "custom_codes": []string{"127.0.0.2", "127.0.0.3"},
}
source, err := client.RBL().Sources().Update("RB18c470cc518a09678bb280960dbdd524", params)

// Delete a source
result, err := client.RBL().Sources().Delete("RB18c470cc518a09678bb280960dbdd524")
```

#### Check & Listings

```go
// Start a manual RBL check
result, err := client.RBL().Check().Start(map[string]interface{}{
    "host": "8.8.8.8",
})

// Get check status
status, err := client.RBL().Check().Status("check_id")

// Get current listings
listings, err := client.RBL().Listings()
```

### Certificate Monitoring

Certificate monitoring allows you to monitor SSL/TLS certificates for expiration, validity, and configuration issues across HTTPS, SMTPS, IMAPS, and other TLS-enabled services.

#### Errors

```go
// Get all certificate errors
errors, err := client.Cert().Errors().Get()
if err != nil {
    log.Fatal(err)
}

// Get a specific error
error, err := client.Cert().Errors().Get("CE5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a")
if err != nil {
    log.Fatal(err)
}
```

#### Monitors

```go
// Get all certificate monitors
monitors, err := client.Cert().Monitors().Get()
if err != nil {
    log.Fatal(err)
}

// Get a specific monitor
monitor, err := client.Cert().Monitors().Get("CM62944aeeee2b46d7a28221164f38976a")
if err != nil {
    log.Fatal(err)
}

// Create a certificate monitor
params := map[string]interface{}{
    "name": "Production Web Server",
    "hostname": "example.com",
    "protocol": "https",
    "profile": "CP79b597e61a984a35b5eb7dcdbc3de53c",
    "contact_group": []string{
        "CG4f3e2d1c0b9a8776655443322110fedc",
        "CG5a6b7c8d9e0f1234567890abcdef1234",
    },
    "tags": []string{"production", "web", "ssl"},
}
monitor, err := client.Cert().Monitors().Create(params)
if err != nil {
    log.Fatal(err)
}

// Update a monitor
params := map[string]interface{}{
    "name": "Updated Server Name",
    "tags": []string{"production", "web", "ssl"},
}
monitor, err := client.Cert().Monitors().Update("CM62944aeeee2b46d7a28221164f38976a", params)
if err != nil {
    log.Fatal(err)
}

// Delete a monitor
result, err := client.Cert().Monitors().Delete("CM62944aeeee2b46d7a28221164f38976a")
if err != nil {
    log.Fatal(err)
}

// Pause monitoring
result, err := client.Cert().Monitors().Pause("CM62944aeeee2b46d7a28221164f38976a")
if err != nil {
    log.Fatal(err)
}

// Resume monitoring
result, err := client.Cert().Monitors().Resume("CM62944aeeee2b46d7a28221164f38976a")
if err != nil {
    log.Fatal(err)
}
```

#### Profiles

```go
// Get all certificate profiles
profiles, err := client.Cert().Profiles().Get()
if err != nil {
    log.Fatal(err)
}

// Get a specific profile
profile, err := client.Cert().Profiles().Get("CP79b597e61a984a35b5eb7dcdbc3de53c")
if err != nil {
    log.Fatal(err)
}

// Create a profile
params := map[string]interface{}{
    "name": "Standard Certificate Profile",
    "expiration_thresholds": []int{30, 14, 7},
    "alert_on_expiration":   true,
    "alert_on_name_mismatch":       true,
    "alert_on_misconfigurations":   true,
    "alert_on_changes":             true,
}
profile, err := client.Cert().Profiles().Create(params)
if err != nil {
    log.Fatal(err)
}

// Update a profile
params = map[string]interface{}{
    "expiration_thresholds":      []int{45, 14, 7},
    "alert_on_misconfigurations": true,
    "alert_on_changes":           true,
}
profile, err = client.Cert().Profiles().Update("CP79b597e61a984a35b5eb7dcdbc3de53c", params)
if err != nil {
    log.Fatal(err)
}

// Delete a profile
result, err := client.Cert().Profiles().Delete("CP79b597e61a984a35b5eb7dcdbc3de53c")
if err != nil {
    log.Fatal(err)
}
```

### Contact Management

#### Contacts

```go
// Get all contacts
contacts, err := client.Contact().Contacts().Get()

// Get a specific contact
contact, err := client.Contact().Contacts().Get("COabcdef1234567890abcdef1234567890")

// Get multiple contacts
contacts, err := client.Contact().Contacts().Get("COabcdef1234567890abcdef1234567890", "CO1234567890abcdef1234567890abcdef")

// Create a contact
params := map[string]interface{}{
    "contact": "user@example.com",
    "type":    "email",
    "schedule": "every_check",
    "contact_group": []string{
        "CG4f3e2d1c0b9a8776655443322110fedc",
        "CG5a6b7c8d9e0f1234567890abcdef1234",
    },
}
contact, err := client.Contact().Contacts().Create(params)

// Update a contact
params := map[string]interface{}{
    "contact": "updated@example.com",
    "contact_group": []string{
        "CG4f3e2d1c0b9a8776655443322110fedc",
        "CG5a6b7c8d9e0f1234567890abcdef1234",
    },
}
contact, err := client.Contact().Contacts().Update("COabcdef1234567890abcdef1234567890", params)

// Delete a contact
result, err := client.Contact().Contacts().Delete("COabcdef1234567890abcdef1234567890")
```

#### Groups

```go
// Get all groups
groups, err := client.Contact().Groups().Get()

// Get a specific group
group, err := client.Contact().Groups().Get("CG4f3e2d1c0b9a8776655443322110fedc")

// Create a group
params := map[string]interface{}{
    "name": "Primary Contacts",
}
group, err := client.Contact().Groups().Create(params)

// Update a group
params := map[string]interface{}{
    "name": "Updated Group Name",
}
group, err := client.Contact().Groups().Update("CG4f3e2d1c0b9a8776655443322110fedc", params)

// Delete a group
result, err := client.Contact().Groups().Delete("CG4f3e2d1c0b9a8776655443322110fedc")
```

## Error Handling

All methods return an error as the second return value. Always check for errors:

```go
result, err := client.RBL().Check("8.8.8.8")
if err != nil {
    log.Printf("API error: %v", err)
    return
}
```

## Retry Logic

The SDK automatically retries failed requests with exponential backoff:
- Configurable maximum retry attempts (default: 3)
- Retries on connection errors, 5xx server errors, and 429 rate limits
- Respects `Retry-After` header on 429 responses before falling back to exponential backoff
- Configurable exponential backoff multiplier (default: 1.0 for 1s, 2s, 4s delays)
- Configurable request timeout (default: 30 seconds)

Customize retry behavior via the `Config` struct:

```go
config := &generatorlabs.Config{
    MaxRetries:   5,    // More retry attempts
    RetryBackoff: 2.0,  // Faster exponential growth
    Timeout:      60 * time.Second,  // Longer timeout
}
client, err := generatorlabs.New(accountSID, authToken, config)
```

## Rate Limiting

The API enforces two layers of rate limiting:

- **Hourly limit**: 1,000 requests per hour per application
- **Per-second limit**: varies by endpoint — 100 RPS for read operations, 50 RPS for write operations, and 20 RPS for manual check start

When a rate limit is exceeded, the API returns HTTP 429 with a `Retry-After` header indicating how many seconds to wait. The SDK automatically respects this header during retries.

All API responses include IETF draft rate limit headers, accessible via the `RateLimit` field on every response:

| Header | Description | Example |
|--------|-------------|---------|
| `RateLimit-Limit` | Active rate limit policies | `1000;w=3600, 100;w=1` |
| `RateLimit-Remaining` | Requests remaining in the most restrictive window | `95` |
| `RateLimit-Reset` | Seconds until the most restrictive window resets | `1` |

```go
response, err := client.RBL().Hosts().Get()
if err != nil {
    log.Fatal(err)
}

// Access response data
hosts := response.Data["data"]

// Access rate limit info
if response.RateLimit != nil {
    fmt.Printf("Remaining: %d, Reset: %ds\n",
        response.RateLimit.Remaining,
        response.RateLimit.Reset)
}
```

## Examples

The `examples/` directory contains complete, runnable examples demonstrating:

- **check_ip.go**: Check if an IP is listed on any RBLs
- **manage_hosts.go**: Create, list, update, and delete monitored hosts
- **pagination.go**: Handle large result sets with automatic pagination
- **error_handling.go**: Proper error handling and custom configuration

Run examples:

```bash
export GENERATOR_LABS_ACCOUNT_SID="your_account_sid"
export GENERATOR_LABS_AUTH_TOKEN="your_auth_token"
go run examples/check_ip.go
```

## Requirements

- Go 1.21 or higher
- Valid Generator Labs API credentials (account SID and auth token)

## Testing

```bash
go test -v
```

## Security

For security best practices and vulnerability reporting, see [SECURITY.md](SECURITY.md).

## Release History

### v2.0.0 (2026-01-31)
* Complete rewrite for Generator Labs API v4.0
* RESTful endpoint design with proper HTTP verbs
* Updated to use Generator Labs branding (formerly RBLTracker)
* Automatic `Retry-After` header support on 429 rate limit responses
* `Response` wrapper exposes per-request rate limit info (`RateLimit`)
* Added `RateLimitInfo` struct with `Limit`, `Remaining`, and `Reset` fields
* Automatic retry with exponential backoff on 429 and 5xx errors
* Webhook signature verification with HMAC-SHA256 and constant-time comparison
* Automatic pagination via `GetAll()` for large result sets
* Standard library only — no external dependencies

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Documentation: https://docs.generatorlabs.com/api/v4/
- Issues: https://github.com/generator-labs/go-sdk/issues
- Email: support@generatorlabs.com

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
