// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArrayParamsConvertedToCommaSeparated(t *testing.T) {
	var capturedBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		capturedBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status_code":    200,
			"status_message": "OK",
		})
	}))
	defer server.Close()

	client, err := New(
		"AC0123456789abcdef0123456789abcdef",
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		&Config{BaseURL: server.URL + "/"},
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("string slice in POST", func(t *testing.T) {
		_, err := client.RBL().Hosts().Create(map[string]interface{}{
			"name": "Test Host",
			"host": "1.2.3.4",
			"contact_group": []string{
				"CG11111111111111111111111111111111",
				"CG22222222222222222222222222222222",
			},
		})
		if err != nil {
			t.Fatalf("Create() error: %v", err)
		}

		if !containsParam(capturedBody, "contact_group", "CG11111111111111111111111111111111,CG22222222222222222222222222222222") {
			t.Errorf("Expected contact_group as comma-separated string, got body: %s", capturedBody)
		}
	})

	t.Run("interface slice in POST", func(t *testing.T) {
		_, err := client.RBL().Hosts().Create(map[string]interface{}{
			"name": "Test Host",
			"host": "1.2.3.4",
			"contact_group": []interface{}{
				"CG11111111111111111111111111111111",
				"CG22222222222222222222222222222222",
			},
		})
		if err != nil {
			t.Fatalf("Create() error: %v", err)
		}

		if !containsParam(capturedBody, "contact_group", "CG11111111111111111111111111111111,CG22222222222222222222222222222222") {
			t.Errorf("Expected contact_group as comma-separated string, got body: %s", capturedBody)
		}
	})

	t.Run("string value unchanged", func(t *testing.T) {
		_, err := client.RBL().Hosts().Create(map[string]interface{}{
			"name":          "Test Host",
			"host":          "1.2.3.4",
			"contact_group": "CG11111111111111111111111111111111",
		})
		if err != nil {
			t.Fatalf("Create() error: %v", err)
		}

		if !containsParam(capturedBody, "contact_group", "CG11111111111111111111111111111111") {
			t.Errorf("Expected contact_group as single string, got body: %s", capturedBody)
		}
	})
}

// containsParam checks if a URL-encoded body contains a specific key=value pair.
func containsParam(body, key, value string) bool {
	expected := key + "=" + value
	// URL-encode commas
	encoded := key + "=" + urlEncode(value)
	return body == expected || contains(body, expected) || contains(body, encoded)
}

func urlEncode(s string) string {
	result := ""
	for _, c := range s {
		if c == ',' {
			result += "%2C"
		} else {
			result += string(c)
		}
	}
	return result
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
