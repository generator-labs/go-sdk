// Example: Manage monitored hosts (create, list, update, delete)
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

	// List all hosts
	fmt.Println("=== Listing all monitored hosts ===")
	hosts, err := client.RBL().Hosts().Get()
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	if hostList, ok := hosts.Data["hosts"].([]interface{}); ok {
		fmt.Printf("Total hosts: %d\n\n", len(hostList))
		for _, h := range hostList {
			if host, ok := h.(map[string]interface{}); ok {
				fmt.Printf("ID: %v, IP: %v, Description: %v\n",
					host["id"], host["ip"], host["description"])
			}
		}
	}

	// Create a new host
	fmt.Println("\n=== Creating a new host ===")
	newHost, err := client.RBL().Hosts().Create(map[string]interface{}{
		"ip":          "203.0.113.10",
		"description": "Example host from Go SDK",
		"profile_id":  1, // Use your profile ID
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	var hostID string
	if hostData, ok := newHost.Data["host"].(map[string]interface{}); ok {
		hostID = fmt.Sprintf("%v", hostData["id"])
		fmt.Printf("Created host ID: %s\n", hostID)
	}

	// Get specific host
	fmt.Println("\n=== Getting specific host ===")
	host, err := client.RBL().Hosts().Get(hostID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Println("Host details:")
	fmt.Printf("%+v\n", host)

	// Update host
	fmt.Println("\n=== Updating host ===")
	_, err = client.RBL().Hosts().Update(hostID, map[string]interface{}{
		"description": "Updated description from Go SDK",
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Println("Updated host description")

	// Delete host
	fmt.Println("\n=== Deleting host ===")
	_, err = client.RBL().Hosts().Delete(hostID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Deleted host ID: %s\n", hostID)
}
