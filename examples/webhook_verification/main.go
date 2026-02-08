// Example: Verifying webhook signatures
//
// This example shows how to verify incoming webhook requests from Generator Labs
// using the SDK's built-in signature verification helper.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
	// Your webhook's signing secret, available in the Edit Webhook panel of the Portal.
	// Store this securely (e.g., environment variable), never hard-code it.
	signingSecret := os.Getenv("GENERATOR_LABS_WEBHOOK_SECRET")
	if signingSecret == "" {
		log.Fatal("Error: Set GENERATOR_LABS_WEBHOOK_SECRET environment variable")
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		header := r.Header.Get("X-Webhook-Signature")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		//
		// Example 1: Basic verification
		//
		// Verify the signature with the default 5-minute tolerance window.
		// On success, returns the decoded JSON payload as a map.
		// Returns an error if verification fails.
		//
		payload, err := generatorlabs.VerifyWebhook(string(body), header, signingSecret, generatorlabs.DefaultWebhookTolerance)
		if err != nil {
			log.Printf("Verification failed: %v", err)
			http.Error(w, `{"error":"Invalid signature"}`, http.StatusForbidden)
			return
		}

		fmt.Println("Webhook verified successfully!")

		// Process the event
		event, _ := payload["event"].(string)
		fmt.Printf("Event: %s\n", event)

		switch event {
		case "rbl.host.listed":
			// Handle host listed event
		case "rbl.host.delisted":
			// Handle host delisted event
		case "billing.balance.alert":
			// Handle low balance alert
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	//
	// Example 2: Custom tolerance
	//
	// Set a custom tolerance window (in seconds) for timestamp validation.
	// Use 0 to disable timestamp checking entirely.
	//
	// payload, err := generatorlabs.VerifyWebhook(string(body), header, signingSecret, 600)  // 10-minute tolerance
	// payload, err := generatorlabs.VerifyWebhook(string(body), header, signingSecret, 0)    // disable timestamp check

	fmt.Println("Webhook endpoint listening on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
