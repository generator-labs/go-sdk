// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"fmt"
	"regexp"
)

// VERSION is the current version of the SDK
const VERSION = "2.0.0"

// Client represents the Generator Labs API client
type Client struct {
	AccountSID string
	AuthToken  string
	Config     *Config
	handler    *RequestHandler
	rbl        *RBL
	contact    *Contact
	cert       *Cert
}

// New creates a new Generator Labs API client
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

// RBL returns the RBL monitoring API namespace
func (c *Client) RBL() *RBL {
	if c.rbl == nil {
		c.rbl = &RBL{handler: c.handler}
	}
	return c.rbl
}

// Contact returns the Contact management API namespace
func (c *Client) Contact() *Contact {
	if c.contact == nil {
		c.contact = &Contact{handler: c.handler}
	}
	return c.contact
}

// Cert returns the Certificate monitoring API namespace
func (c *Client) Cert() *Cert {
	if c.cert == nil {
		c.cert = &Cert{handler: c.handler}
	}
	return c.cert
}
