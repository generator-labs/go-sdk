# Generator Labs Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/generator-labs/go-sdk.svg)](https://pkg.go.dev/github.com/generator-labs/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/generator-labs/go-sdk)](https://goreportcard.com/report/github.com/generator-labs/go-sdk)
[![Tests](https://github.com/generator-labs/go-sdk/workflows/Tests/badge.svg)](https://github.com/generator-labs/go-sdk/actions)
[![CodeQL](https://github.com/generator-labs/go-sdk/workflows/CodeQL/badge.svg)](https://github.com/generator-labs/go-sdk/actions?query=workflow%3ACodeQL)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [Generator Labs API](https://generatorlabs.com). This library provides a simple and intuitive interface for interacting with the Generator Labs v4.0 API, including RBL monitoring, contact management, and more.

## Features

- Full support for Generator Labs API v4.0
- Configurable timeouts, retries, and backoff strategies
- Automatic retry logic with exponential backoff
- Automatic pagination for large result sets
- Connection pooling and timeout management
- Type-safe API methods
- Comprehensive error handling
- Security scanning with CodeQL and Dependabot
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

    // Check an IP address
    result, err := client.RBL().Check("8.8.8.8")
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
host, err := client.RBL().Hosts().Get(123)

// Get multiple hosts
hosts, err := client.RBL().Hosts().Get(123, 456, 789)

// Create a host
params := map[string]interface{}{
    "ip": "8.8.8.8",
    "description": "Google DNS",
}
host, err := client.RBL().Hosts().Create(params)

// Update a host
params := map[string]interface{}{"description": "Updated description"}
host, err := client.RBL().Hosts().Update(123, params)

// Delete a host
result, err := client.RBL().Hosts().Delete(123)
```

#### Profiles

```go
// Get all profiles
profiles, err := client.RBL().Profiles().Get()

// Get a specific profile
profile, err := client.RBL().Profiles().Get(1)

// Create/Update/Delete - similar to Hosts
```

#### Sources

```go
// Get all RBL sources
sources, err := client.RBL().Sources().Get()

// Get a specific source
source, err := client.RBL().Sources().Get(10)

// Create/Update/Delete - similar to Hosts
```

#### Check & Listings

```go
// Check an IP address
result, err := client.RBL().Check("8.8.8.8")

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
    "port": 443,
    "protocol": "https",
    "cert_profile": "CP79b597e61a984a35b5eb7dcdbc3de53c",
    "contact_group": "CG4f3e2d1c0b9a8776655443322110fed",
}
monitor, err := client.Cert().Monitors().Create(params)
if err != nil {
    log.Fatal(err)
}

// Update a monitor
params := map[string]interface{}{"name": "Updated Server Name"}
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
    "expiration_warning_days": 30,
    "expiration_critical_days": 7,
}
profile, err := client.Cert().Profiles().Create(params)
if err != nil {
    log.Fatal(err)
}

// Update a profile
params := map[string]interface{}{"expiration_warning_days": 45}
profile, err := client.Cert().Profiles().Update("CP79b597e61a984a35b5eb7dcdbc3de53c", params)
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
contact, err := client.Contact().Contacts().Get(456)

// Get multiple contacts
contacts, err := client.Contact().Contacts().Get(456, 789)

// Create a contact
params := map[string]interface{}{
    "email": "user@example.com",
    "name": "John Doe",
}
contact, err := client.Contact().Contacts().Create(params)

// Update a contact
params := map[string]interface{}{"name": "Jane Doe"}
contact, err := client.Contact().Contacts().Update(456, params)

// Delete a contact
result, err := client.Contact().Contacts().Delete(456)
```

#### Groups

```go
// Get all groups
groups, err := client.Contact().Groups().Get()

// Get a specific group
group, err := client.Contact().Groups().Get(10)

// Create/Update/Delete - similar to Contacts
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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Documentation: https://docs.generatorlabs.com/api/v4/
- Issues: https://github.com/generator-labs/go-sdk/issues
- Email: support@generatorlabs.com

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
