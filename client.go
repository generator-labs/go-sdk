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
	handler    *RequestHandler
	rbl        *RBL
	contact    *Contact
}

// NewClient creates a new Generator Labs API client
func NewClient(accountSID, authToken string) (*Client, error) {
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

	client := &Client{
		AccountSID: accountSID,
		AuthToken:  authToken,
	}

	// Initialize request handler
	client.handler = NewRequestHandler(accountSID, authToken, "https://api.generatorlabs.com/4.0/")

	return client, nil
}

// RBLClient returns the RBL monitoring API namespace
func (c *Client) RBLClient() *RBL {
	if c.rbl == nil {
		c.rbl = &RBL{handler: c.handler}
	}
	return c.rbl
}

// ContactClient returns the Contact management API namespace
func (c *Client) ContactClient() *Contact {
	if c.contact == nil {
		c.contact = &Contact{handler: c.handler}
	}
	return c.contact
}
