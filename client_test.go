// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		accountSID string
		authToken  string
		wantErr    bool
	}{
		{
			name:       "valid credentials",
			accountSID: "AC0123456789abcdef0123456789abcdef",
			authToken:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr:    false,
		},
		{
			name:       "invalid account SID format",
			accountSID: "invalid",
			authToken:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			wantErr:    true,
		},
		{
			name:       "invalid auth token format",
			accountSID: "AC0123456789abcdef0123456789abcdef",
			authToken:  "invalid",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.accountSID, tt.authToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClientNamespaces(t *testing.T) {
	client, err := NewClient("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("RBLClient", func(t *testing.T) {
		rbl := client.RBLClient()
		if rbl == nil {
			t.Error("RBLClient() returned nil")
		}
	})

	t.Run("ContactClient", func(t *testing.T) {
		contact := client.ContactClient()
		if contact == nil {
			t.Error("ContactClient() returned nil")
		}
	})
}

func TestVersion(t *testing.T) {
	if VERSION != "2.0.0" {
		t.Errorf("VERSION = %s, want 2.0.0", VERSION)
	}
}
