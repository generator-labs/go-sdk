// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper to build a paginated JSON response
func paginatedResponse(items []interface{}, key string, page, totalPages, total int) map[string]interface{} {
	return map[string]interface{}{
		"status_code":    200,
		"status_message": "OK",
		key:              items,
		"total":          float64(total),
		"page":           float64(page),
		"total_pages":    float64(totalPages),
		"page_size":      float64(2),
	}
}

func makeHosts(count, offset int) []interface{} {
	items := make([]interface{}, count)
	for i := 0; i < count; i++ {
		items[i] = map[string]interface{}{
			"name": fmt.Sprintf("host_%d", i+offset+1),
		}
	}
	return items
}

func TestGetAllSinglePage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := paginatedResponse(makeHosts(3, 0), "hosts", 1, 1, 3)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", &Config{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	hosts, err := client.RBL().Hosts().GetAll(nil)
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	if len(hosts) != 3 {
		t.Errorf("GetAll() returned %d items, want 3", len(hosts))
	}
}

func TestGetAllMultiplePages(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		page := callCount

		var items []interface{}
		if page <= 2 {
			items = makeHosts(2, (page-1)*2)
		} else {
			items = makeHosts(1, 4)
		}

		resp := paginatedResponse(items, "hosts", page, 3, 5)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", &Config{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	hosts, err := client.RBL().Hosts().GetAll(map[string]interface{}{"page_size": 2})
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	if len(hosts) != 5 {
		t.Errorf("GetAll() returned %d items, want 5", len(hosts))
	}

	if callCount != 3 {
		t.Errorf("GetAll() made %d API calls, want 3", callCount)
	}
}

func TestGetAllEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := paginatedResponse([]interface{}{}, "hosts", 1, 1, 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", &Config{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	hosts, err := client.RBL().Hosts().GetAll(nil)
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	if len(hosts) != 0 {
		t.Errorf("GetAll() returned %d items, want 0", len(hosts))
	}
}

func TestGetAllCustomPageSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageSize := r.URL.Query().Get("page_size")
		if pageSize != "50" {
			t.Errorf("Expected page_size=50, got %s", pageSize)
		}

		resp := paginatedResponse(makeHosts(1, 0), "hosts", 1, 1, 1)
		resp["page_size"] = float64(50)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", &Config{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	hosts, err := client.RBL().Hosts().GetAll(map[string]interface{}{"page_size": 50})
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	if len(hosts) != 1 {
		t.Errorf("GetAll() returned %d items, want 1", len(hosts))
	}
}

func TestGetAllContactsResourceKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := paginatedResponse(
			[]interface{}{
				map[string]interface{}{"name": "contact_a"},
				map[string]interface{}{"name": "contact_b"},
			},
			"contacts", 1, 1, 2,
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := New("AC0123456789abcdef0123456789abcdef", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", &Config{BaseURL: server.URL + "/"})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	contacts, err := client.Contact().Contacts().GetAll(nil)
	if err != nil {
		t.Fatalf("GetAll() error: %v", err)
	}

	if len(contacts) != 2 {
		t.Errorf("GetAll() returned %d items, want 2", len(contacts))
	}
}
