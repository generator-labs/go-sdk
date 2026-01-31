# Generator Labs Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/generator-labs/go-sdk.svg)](https://pkg.go.dev/github.com/generator-labs/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/generator-labs/go-sdk)](https://goreportcard.com/report/github.com/generator-labs/go-sdk)
[![Tests](https://github.com/generator-labs/go-sdk/workflows/Tests/badge.svg)](https://github.com/generator-labs/go-sdk/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for the [Generator Labs API](https://generatorlabs.com). This library provides a simple and intuitive interface for interacting with the Generator Labs v4.0 API, including RBL monitoring, contact management, and more.

## Features

- Full support for Generator Labs API v4.0
- Automatic retry logic with exponential backoff
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

    generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
    // Initialize the client
    client, err := generatorlabs.New(
        "YOUR_ACCOUNT_SID",
        "YOUR_AUTH_TOKEN",
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

    // Get all contacts
    contacts, err := client.Contact().Contacts().Get()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Contacts: %+v\n", contacts)
}
```

## API Reference

### Client Initialization

```go
client, err := generatorlabs.New(accountSID, authToken)
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
- Maximum 3 retry attempts
- Retries on connection errors, 5xx server errors, and 429 rate limits
- Exponential backoff delays: 1s, 2s, 4s
- 30 second request timeout

## Requirements

- Go 1.21 or higher
- Valid Generator Labs API credentials (account SID and auth token)

## Testing

```bash
go test -v
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Documentation: https://docs.generatorlabs.com/api/v4/
- Issues: https://github.com/generator-labs/go-sdk/issues
- Email: support@generatorlabs.com

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
