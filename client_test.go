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

func TestNew(t *testing.T) {
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
			client, err := New(tt.accountSID, tt.authToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("New() returned nil client")
			}
		})
	}
}

func TestClientNamespaces(t *testing.T) {
	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("RBL", func(t *testing.T) {
		rbl := client.RBL()
		if rbl == nil {
			t.Error("RBL() returned nil")
		}
	})

	t.Run("Contact", func(t *testing.T) {
		contact := client.Contact()
		if contact == nil {
			t.Error("Contact() returned nil")
		}
	})

	t.Run("Cert", func(t *testing.T) {
		cert := client.Cert()
		if cert == nil {
			t.Error("Cert() returned nil")
		}
	})

	t.Run("Cert Errors", func(t *testing.T) {
		errors := client.Cert().Errors()
		if errors == nil {
			t.Error("Cert().Errors() returned nil")
		}
	})

	t.Run("Cert Monitors", func(t *testing.T) {
		monitors := client.Cert().Monitors()
		if monitors == nil {
			t.Error("Cert().Monitors() returned nil")
		}
	})

	t.Run("Cert Profiles", func(t *testing.T) {
		profiles := client.Cert().Profiles()
		if profiles == nil {
			t.Error("Cert().Profiles() returned nil")
		}
	})
}

func TestVersion(t *testing.T) {
	if VERSION != "2.0.1" {
		t.Errorf("VERSION = %s, want 2.0.1", VERSION)
	}
}
