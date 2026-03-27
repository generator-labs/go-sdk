// Example: Proper error handling and configuration
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
	accountSid := os.Getenv("GENERATOR_LABS_ACCOUNT_SID")
	authToken := os.Getenv("GENERATOR_LABS_AUTH_TOKEN")

	if accountSid == "" || authToken == "" {
		log.Fatal("Error: Set GENERATOR_LABS_ACCOUNT_SID and GENERATOR_LABS_AUTH_TOKEN environment variables")
	}

	// Initialize client with custom configuration
	config := &generatorlabs.Config{
		Timeout:        45 * time.Second,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     5,
		RetryBackoff:   2.0,
	}
	client, err := generatorlabs.New(accountSid, authToken, config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== Example 1: Handling API errors ===")
	_, err = client.RBL().Hosts().Get("999999")
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
		fmt.Println("This is expected for a non-existent resource\n")
	}

	fmt.Println("=== Example 2: Invalid credentials ===")
	_, err = generatorlabs.New("INVALID", authToken)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
		fmt.Println("Credential validation works!\n")
	}

	fmt.Println("=== Example 3: Network resilience ===")
	// The SDK automatically retries on:
	// - Connection errors
	// - 5xx server errors
	// - 429 rate limit errors
	// With exponential backoff

	_, err = client.RBL().Check().Start(map[string]interface{}{"host": "1.1.1.1"})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Println("Request succeeded (with automatic retries if needed)")

	fmt.Println("\n=== Example 4: Graceful degradation ===")
	hosts, err := client.RBL().Hosts().Get()
	if err != nil {
		// Log error and continue with cached/default data
		fmt.Printf("API error: %v\n", err)
		fmt.Println("Using cached data due to API error")
		hosts = &generatorlabs.Response{Data: map[string]interface{}{"hosts": []interface{}{}}}
	} else {
		if hostList, ok := hosts.Data["hosts"].([]interface{}); ok {
			fmt.Printf("Successfully retrieved %d hosts\n", len(hostList))
		}
	}

	fmt.Println("\nAll examples completed!")
}
