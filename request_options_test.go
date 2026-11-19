// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// recorded is one request the test API received
type recorded struct {
	method string
	path   string
	query  url.Values
	body   url.Values
	header http.Header
}

// reply is one answer of the test API
type reply struct {
	status  int
	headers map[string]string
	body    string
}

const okBody = `{"status_code":200,"status_message":"OK"}`

// testAPI starts a local API that records each request and answers with the given
// replies in order, then with 200 OK
func testAPI(t *testing.T, replies ...reply) (*Client, *[]recorded) {
	t.Helper()

	requests := &[]recorded{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(body))
		*requests = append(*requests, recorded{
			method: r.Method,
			path:   r.URL.Path,
			query:  r.URL.Query(),
			body:   form,
			header: r.Header,
		})

		answer := reply{status: 200, body: okBody}
		if len(replies) > 0 {
			answer = replies[0]
			replies = replies[1:]
		}
		if answer.status == 0 {
			answer.status = 200
		}
		if answer.body == "" && answer.status != http.StatusNotModified {
			answer.body = okBody
		}

		for name, value := range answer.headers {
			w.Header().Set(name, value)
		}
		w.WriteHeader(answer.status)
		_, _ = w.Write([]byte(answer.body))
	}))
	t.Cleanup(server.Close)

	config := DefaultConfig()
	config.BaseURL = server.URL + "/4.0/"

	client, err := New(
		"AC0123456789abcdef0123456789abcdef",
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		config,
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client, requests
}

func errorBody(status string) string {
	return `{"status_code":` + status + `,"status_message":"Request failed."}`
}

func TestIdempotencyKeySameOnEveryRetry(t *testing.T) {
	client, requests := testAPI(t, reply{status: 500, body: errorBody("500")}, reply{})

	if _, err := client.RBL().Hosts().Create(map[string]interface{}{"name": "Test", "host": "192.0.2.1"}); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if len(*requests) != 2 {
		t.Fatalf("got %d requests, want 2", len(*requests))
	}
	key := (*requests)[0].header.Get("Idempotency-Key")
	if len(key) != 32 {
		t.Errorf("Idempotency-Key = %q, want 32 characters", key)
	}
	if retried := (*requests)[1].header.Get("Idempotency-Key"); retried != key {
		t.Errorf("retry sent Idempotency-Key %q, want %q", retried, key)
	}
}

func TestIdempotencyKeyNewForEachCall(t *testing.T) {
	client, requests := testAPI(t)

	_, _ = client.RBL().Hosts().Create(map[string]interface{}{"name": "One"})
	_, _ = client.RBL().Hosts().Create(map[string]interface{}{"name": "Two"})

	if (*requests)[0].header.Get("Idempotency-Key") == (*requests)[1].header.Get("Idempotency-Key") {
		t.Errorf("two calls sent the same Idempotency-Key")
	}
}

func TestIdempotencyKeyFromCaller(t *testing.T) {
	client, requests := testAPI(t)

	_, _ = client.RBL().Check().Start(map[string]interface{}{"host": "192.0.2.1"}, WithIdempotencyKey("my-key"))

	if key := (*requests)[0].header.Get("Idempotency-Key"); key != "my-key" {
		t.Errorf("Idempotency-Key = %q, want my-key", key)
	}
}

func TestIdempotencyKeyOnlyWhereAccepted(t *testing.T) {
	client, requests := testAPI(t)

	_, _ = client.Contact().Contacts().Resend("COx")
	_, _ = client.Contact().Contacts().Confirm("COx", map[string]interface{}{"authcode": "1234"})
	_, _ = client.Cert().Monitors().Delete("CMx")
	_, _ = client.Cert().Monitors().Check("CMx")
	_, _ = client.RBL().Hosts().Update("HTx", map[string]interface{}{"name": "Test"})
	_, _ = client.RBL().Hosts().Pause("HTx")
	_, _ = client.RBL().Hosts().Delete("HTx")

	if len(*requests) != 7 {
		t.Fatalf("got %d requests, want 7", len(*requests))
	}
	for i, request := range *requests {
		key := request.header.Get("Idempotency-Key")
		if i < 4 && key == "" {
			t.Errorf("%s %s sent no Idempotency-Key", request.method, request.path)
		}
		if i >= 4 && key != "" {
			t.Errorf("%s %s sent an Idempotency-Key", request.method, request.path)
		}
	}
}

func TestReplayed(t *testing.T) {
	client, _ := testAPI(t, reply{headers: map[string]string{"Idempotent-Replayed": "true"}}, reply{})

	first, _ := client.RBL().Hosts().Create(map[string]interface{}{"name": "Test"})
	second, _ := client.RBL().Hosts().Create(map[string]interface{}{"name": "Test"})

	if !first.Replayed || second.Replayed {
		t.Errorf("Replayed = %v, %v, want true, false", first.Replayed, second.Replayed)
	}
}

func TestRequestID(t *testing.T) {
	client, requests := testAPI(t,
		reply{headers: map[string]string{"X-Request-Id": "req-1"}},
		reply{status: 404, headers: map[string]string{"X-Request-Id": "req-2"}, body: errorBody("404")},
	)

	response, err := client.RBL().Hosts().Get(WithRequestID("req-1"))
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if sent := (*requests)[0].header.Get("X-Request-Id"); sent != "req-1" {
		t.Errorf("sent X-Request-Id %q, want req-1", sent)
	}
	if response.RequestID != "req-1" {
		t.Errorf("RequestID = %q, want req-1", response.RequestID)
	}

	_, err = client.RBL().Hosts().Get("HTx")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Get() error = %v, want an APIError", err)
	}
	if apiErr.StatusCode != 404 || apiErr.RequestID != "req-2" {
		t.Errorf("APIError = %d, %q, want 404, req-2", apiErr.StatusCode, apiErr.RequestID)
	}
}

func TestETagAndNotModified(t *testing.T) {
	client, requests := testAPI(t,
		reply{headers: map[string]string{"ETag": `W/"abc"`}},
		reply{status: 304, headers: map[string]string{"ETag": `W/"abc"`}},
	)

	first, err := client.RBL().Hosts().Get()
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	second, err := client.RBL().Hosts().Get(WithIfNoneMatch(first.ETag))
	if err != nil {
		t.Fatalf("Get() with If-None-Match error: %v", err)
	}

	if first.ETag != `W/"abc"` || first.NotModified {
		t.Errorf("first = %q, %v, want the ETag and not NotModified", first.ETag, first.NotModified)
	}
	if sent := (*requests)[1].header.Get("If-None-Match"); sent != `W/"abc"` {
		t.Errorf("sent If-None-Match %q", sent)
	}
	if !second.NotModified || len(second.Data) != 0 {
		t.Errorf("second = %v, %v, want NotModified with no data", second.NotModified, second.Data)
	}
}

func TestFields(t *testing.T) {
	client, requests := testAPI(t)

	_, _ = client.RBL().Hosts().Get("HTx", WithFields("id", "host"))
	_, _ = client.RBL().Hosts().Get(map[string]interface{}{"page": 2}, WithFields("id", "host"))
	_, _ = client.RBL().Listings(map[string]interface{}{"page": 2})

	single, list, listings := (*requests)[0], (*requests)[1], (*requests)[2]
	if single.path != "/4.0/rbl/hosts/HTx.json" || single.query.Get("fields") != "id,host" {
		t.Errorf("single get = %s?%s", single.path, single.query.Encode())
	}
	if list.path != "/4.0/rbl/hosts.json" || list.query.Get("fields") != "id,host" || list.query.Get("page") != "2" {
		t.Errorf("list = %s?%s", list.path, list.query.Encode())
	}
	if listings.path != "/4.0/rbl/listings.json" || listings.query.Get("page") != "2" {
		t.Errorf("listings = %s?%s", listings.path, listings.query.Encode())
	}
}

func TestGetTakesOneID(t *testing.T) {
	client, requests := testAPI(t)

	if _, err := client.RBL().Hosts().Get("HTx", "HTy"); err == nil {
		t.Errorf("Get() with two ids returned no error")
	}
	if len(*requests) != 0 {
		t.Errorf("got %d requests, want 0", len(*requests))
	}
}

func TestRetryAfter(t *testing.T) {
	t.Run("a long wait is returned as the error", func(t *testing.T) {
		client, requests := testAPI(t,
			reply{status: 429, headers: map[string]string{"Retry-After": "3600"}, body: errorBody("429")},
		)

		_, err := client.Cert().Monitors().Check("CMx")

		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != 429 {
			t.Fatalf("Check() error = %v, want an APIError 429", err)
		}
		if len(*requests) != 1 {
			t.Errorf("got %d requests, want 1", len(*requests))
		}
	})

	t.Run("a short wait is retried", func(t *testing.T) {
		client, requests := testAPI(t,
			reply{status: 429, headers: map[string]string{"Retry-After": "1"}, body: errorBody("429")},
		)

		start := time.Now()
		if _, err := client.RBL().Hosts().Get(); err != nil {
			t.Fatalf("Get() error: %v", err)
		}
		if len(*requests) != 2 {
			t.Errorf("got %d requests, want 2", len(*requests))
		}
		if waited := time.Since(start); waited < time.Second {
			t.Errorf("retried after %v, want the Retry-After of 1s", waited)
		}
	})
}

func TestRetriesRunOut(t *testing.T) {
	failed := reply{status: 503, headers: map[string]string{"X-Request-Id": "req-3"}, body: errorBody("503")}
	client, requests := testAPI(t, failed, failed)
	client.handler.client.RetryMax = 1

	_, err := client.RBL().Hosts().Get()

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Get() error = %v, want an APIError", err)
	}
	if apiErr.StatusCode != 503 || apiErr.StatusMessage != "Request failed." || apiErr.RequestID != "req-3" {
		t.Errorf("APIError = %+v, want 503, the API message and the request id", apiErr)
	}
	if len(*requests) != 2 {
		t.Errorf("got %d requests, want 2", len(*requests))
	}
}

func TestAccount(t *testing.T) {
	client, requests := testAPI(t,
		reply{body: `{"status_code":200,"status_message":"OK","data":{"balance":80.1226}}`},
		reply{body: "80.1226\n"},
		reply{},
		reply{},
		reply{},
		reply{status: 401, body: "-1\n"},
	)

	account, err := client.Account().Get()
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if data, _ := account.Data["data"].(map[string]interface{}); data["balance"] != 80.1226 {
		t.Errorf("account = %v", account.Data)
	}

	balance, err := client.Account().Balance()
	if err != nil || balance != 80.1226 {
		t.Errorf("Balance() = %v, %v, want 80.1226", balance, err)
	}

	_, _ = client.Account().Usage(map[string]interface{}{"months": 3})
	_, _ = client.Account().PauseChecks()
	_, _ = client.Account().ResumeChecks()

	_, err = client.Account().Balance()
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 401 {
		t.Errorf("Balance() error = %v, want an APIError 401", err)
	}

	want := []string{
		"GET /4.0/account.json",
		"GET /4.0/account/balance.json",
		"GET /4.0/account/usage.json",
		"POST /4.0/account/checks/pause.json",
		"POST /4.0/account/checks/resume.json",
		"GET /4.0/account/balance.json",
	}
	for i, request := range *requests {
		if got := request.method + " " + request.path; got != want[i] {
			t.Errorf("request %d = %s, want %s", i, got, want[i])
		}
	}
	if months := (*requests)[2].query.Get("months"); months != "3" {
		t.Errorf("usage months = %q, want 3", months)
	}
}

func TestCertFindings(t *testing.T) {
	client, requests := testAPI(t,
		reply{},
		reply{},
		reply{},
		reply{body: `{"status_code":200,"status_message":"OK","total_pages":2,"data":[{"id":"CM1"}]}`},
		reply{body: `{"status_code":200,"status_message":"OK","total_pages":2,"data":[{"id":"CM2"}]}`},
	)

	_, _ = client.Cert().Findings().Get(map[string]interface{}{"dismissed": 1})
	_, _ = client.Cert().Findings().Dismiss(map[string]interface{}{
		"finding_id": []string{"CERT_CAA01", "CERT_CAA02"},
		"monitors":   []string{"CM1", "CM2"},
	})
	_, _ = client.Cert().Findings().Restore(map[string]interface{}{"finding_id": "CERT_CAA01"})

	all, err := client.Cert().Findings().GetAll(nil)
	if err != nil || len(all) != 2 {
		t.Errorf("GetAll() = %v, %v, want 2 monitors", all, err)
	}

	list, dismiss, restore := (*requests)[0], (*requests)[1], (*requests)[2]
	if list.path != "/4.0/cert/findings.json" || list.query.Get("dismissed") != "1" {
		t.Errorf("list = %s?%s", list.path, list.query.Encode())
	}
	if dismiss.method != "POST" || dismiss.path != "/4.0/cert/findings/dismiss.json" ||
		dismiss.body.Get("finding_id") != "CERT_CAA01,CERT_CAA02" || dismiss.body.Get("monitors") != "CM1,CM2" {
		t.Errorf("dismiss = %s %s %s", dismiss.method, dismiss.path, dismiss.body.Encode())
	}
	if restore.path != "/4.0/cert/findings/restore.json" || restore.body.Get("finding_id") != "CERT_CAA01" {
		t.Errorf("restore = %s %s", restore.path, restore.body.Encode())
	}
}

func TestCertCheckAndCount(t *testing.T) {
	client, requests := testAPI(t, reply{}, reply{body: "3\n"})

	if _, err := client.Cert().Monitors().Check("CMx"); err != nil {
		t.Fatalf("Check() error: %v", err)
	}
	count, err := client.Cert().Errors().Count()
	if err != nil || count != 3 {
		t.Errorf("Count() = %v, %v, want 3", count, err)
	}

	if got := (*requests)[0].method + " " + (*requests)[0].path; got != "POST /4.0/cert/monitors/CMx/check.json" {
		t.Errorf("check = %s", got)
	}
	if (*requests)[1].path != "/4.0/cert/errors/count.json" {
		t.Errorf("count = %s", (*requests)[1].path)
	}
}
