// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// DefaultWebhookTolerance is the default maximum age of a webhook
// request in seconds (5 minutes).
const DefaultWebhookTolerance = 300

// VerifyWebhook verifies a webhook signature and returns the decoded payload.
//
// Parameters:
//   - body: The raw request body string
//   - header: The X-Webhook-Signature header value
//   - secret: Your webhook's signing secret
//   - tolerance: Maximum age in seconds (0 to disable)
//
// Returns the decoded JSON payload as a map, or an error if verification fails.
//
// Example:
//
//	body, _ := io.ReadAll(r.Body)
//	header := r.Header.Get("X-Webhook-Signature")
//
//	payload, err := generatorlabs.VerifyWebhook(string(body), header, signingSecret, 0)
func VerifyWebhook(body, header, secret string, tolerance int) (map[string]interface{}, error) {
	if header == "" {
		return nil, fmt.Errorf("missing X-Webhook-Signature header")
	}

	// Parse the header: t=timestamp,v1=signature
	parts := make(map[string]string)
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}

	t, tOk := parts["t"]
	v1, v1Ok := parts["v1"]
	if !tOk || !v1Ok {
		return nil, fmt.Errorf("invalid X-Webhook-Signature header format")
	}

	// Check timestamp tolerance
	if tolerance > 0 {
		ts, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in X-Webhook-Signature header")
		}
		if math.Abs(float64(time.Now().Unix()-ts)) > float64(tolerance) {
			return nil, fmt.Errorf("webhook timestamp is outside the tolerance window")
		}
	}

	// Compute and compare the signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(t + "." + body))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(v1)) {
		return nil, fmt.Errorf("webhook signature verification failed")
	}

	// Decode and return the payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return nil, fmt.Errorf("failed to decode webhook JSON payload: %w", err)
	}

	return payload, nil
}
