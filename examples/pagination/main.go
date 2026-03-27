// Example: Paginate through large result sets
package main

import (
	"fmt"
	"log"
	"os"

	generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
	accountSid := os.Getenv("GENERATOR_LABS_ACCOUNT_SID")
	authToken := os.Getenv("GENERATOR_LABS_AUTH_TOKEN")

	if accountSid == "" || authToken == "" {
		log.Fatal("Error: Set GENERATOR_LABS_ACCOUNT_SID and GENERATOR_LABS_AUTH_TOKEN environment variables")
	}

	client, err := generatorlabs.New(accountSid, authToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== Fetching all hosts with pagination ===")

	allHosts := []interface{}{}
	page := 1
	pageSize := 50

	for {
		fmt.Printf("Fetching page %d...\n", page)

		response, err := client.RBL().Hosts().Get(map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
		})
		if err != nil {
			log.Fatalf("API Error: %v", err)
		}

		hosts := []interface{}{}
		if hostList, ok := response.Data["hosts"].([]interface{}); ok {
			hosts = hostList
		}

		allHosts = append(allHosts, hosts...)
		fmt.Printf("  Retrieved %d hosts\n", len(hosts))

		// Check if there are more pages
		hasMore := false
		if hm, ok := response.Data["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore {
			break
		}

		page++
	}

	fmt.Printf("\nTotal hosts retrieved: %d\n", len(allHosts))

	// Alternative: Use the built-in pagination helper
	fmt.Println("\n=== Using pagination helper ===")

	allHostsHelper, err := client.RBL().Hosts().GetAll(map[string]interface{}{
		"page_size": 50,
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Total hosts via helper: %d\n", len(allHostsHelper))
}
