// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

// Package generatorlabs provides a client library for the Generator Labs API.
//
// The Generator Labs API is a RESTful web service API that lets customers manage
// their RBL monitoring hosts, certificate monitoring, contacts, and retrieve
// listing information.
//
// The SDK provides three main namespaces:
//   - RBL: RBL monitoring operations (hosts, profiles, sources, checks, listings)
//   - Contact: Contact management operations (contacts, contact groups)
//   - Cert: Certificate monitoring operations (monitors, profiles, errors)
//
// Basic usage:
//
//	client, err := generatorlabs.New("your_account_sid", "your_auth_token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// RBL monitoring
//	hosts, err := client.RBL().Hosts().Get()
//
//	// Certificate monitoring
//	monitors, err := client.Cert().Monitors().Get()
//
//	// Contact management
//	contacts, err := client.Contact().Contacts().Get()
//
// The SDK includes automatic retry logic with exponential backoff for
// connection errors, 5xx server errors, and 429 rate limit errors.
package generatorlabs

import (
	"fmt"
	"regexp"
)

// VERSION is the current version of the SDK.
const VERSION = "2.0.0"

// Client represents the Generator Labs API client and provides access to
// RBL monitoring, contact management, and certificate monitoring operations.
//
// The client handles authentication, request retries, and provides namespace
// methods to access different parts of the API.
type Client struct {
	// AccountSID is the account identifier (format: 2 uppercase + 32 hex chars)
	AccountSID string
	// AuthToken is the authentication token (64 hex characters)
	AuthToken string
	// Config contains timeout, retry, and base URL settings
	Config  *Config
	handler *RequestHandler
	rbl     *RBL
	contact *Contact
	cert    *Cert
}

// New creates a new Generator Labs API client with the provided credentials.
//
// The accountSID must be in the format of 2 uppercase letters followed by
// 32 hexadecimal characters (e.g., "AC" + 32 hex chars).
//
// The authToken must be 64 hexadecimal characters.
//
// An optional Config can be provided to customize timeout and retry behavior.
// If no config is provided, DefaultConfig() values are used.
//
// Example:
//
//	client, err := generatorlabs.New("ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "your_auth_token")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Example with custom configuration:
//
//	config := &generatorlabs.Config{
//	    Timeout: 45 * time.Second,
//	    MaxRetries: 5,
//	}
//	client, err := generatorlabs.New("ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "your_auth_token", config)
func New(accountSID, authToken string, config ...*Config) (*Client, error) {
	// Validate account SID format
	sidPattern := regexp.MustCompile(`^[A-Z]{2}[0-9a-fA-F]{32}$`)
	if !sidPattern.MatchString(accountSID) {
		return nil, fmt.Errorf("invalid account SID format: %s", accountSID)
	}

	// Validate auth token format
	tokenPattern := regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
	if !tokenPattern.MatchString(authToken) {
		return nil, fmt.Errorf("invalid auth token format")
	}

	// Use provided config or default
	var cfg *Config
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}

	client := &Client{
		AccountSID: accountSID,
		AuthToken:  authToken,
		Config:     cfg,
	}

	// Initialize request handler with config
	client.handler = NewRequestHandler(accountSID, authToken, cfg)

	return client, nil
}

// RBL returns the RBL monitoring API namespace.
//
// The RBL namespace provides access to:
//   - Hosts(): Manage monitored hosts (IP addresses and domains)
//   - Profiles(): Manage monitoring profiles (which RBLs to check)
//   - Sources(): Manage RBL sources
//   - Check(): Perform ad-hoc RBL checks
//   - Listings(): Retrieve current RBL listings
//
// Example:
//
//	hosts, err := client.RBL().Hosts().Get()
func (c *Client) RBL() *RBL {
	if c.rbl == nil {
		c.rbl = &RBL{handler: c.handler}
	}
	return c.rbl
}

// Contact returns the Contact management API namespace.
//
// The Contact namespace provides access to:
//   - Contacts(): Manage individual contacts
//   - Groups(): Manage contact groups
//
// Example:
//
//	contacts, err := client.Contact().Contacts().Get()
func (c *Client) Contact() *Contact {
	if c.contact == nil {
		c.contact = &Contact{handler: c.handler}
	}
	return c.contact
}

// Cert returns the Certificate monitoring API namespace.
//
// The Cert namespace provides access to:
//   - Monitors(): Manage certificate monitors
//   - Profiles(): Manage certificate monitoring profiles
//   - Errors(): Retrieve current certificate errors
//
// Example:
//
//	monitors, err := client.Cert().Monitors().Get()
func (c *Client) Cert() *Cert {
	if c.cert == nil {
		c.cert = &Cert{handler: c.handler}
	}
	return c.cert
}
