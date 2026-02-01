// Example: Check if an IP address is listed on any RBLs
package main

import (
	"fmt"
	"log"
	"os"

	generatorlabs "github.com/generator-labs/go-sdk"
)

func main() {
	// Get credentials from environment variables
	accountSid := os.Getenv("GENERATOR_LABS_ACCOUNT_SID")
	authToken := os.Getenv("GENERATOR_LABS_AUTH_TOKEN")

	if accountSid == "" || authToken == "" {
		log.Fatal("Error: Set GENERATOR_LABS_ACCOUNT_SID and GENERATOR_LABS_AUTH_TOKEN environment variables")
	}

	// Initialize client
	client, err := generatorlabs.New(accountSid, authToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Check a single IP address
	ip := "8.8.8.8"
	fmt.Printf("Checking IP: %s\n", ip)

	result, err := client.RBL().Check(ip)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	fmt.Println("Results:")
	fmt.Printf("%+v\n", result)

	// Check if IP is listed
	if listed, ok := result["listed"].(bool); ok && listed {
		fmt.Printf("\nWARNING: IP %s is listed on one or more RBLs!\n", ip)
		if listings, ok := result["listings"].([]interface{}); ok {
			fmt.Printf("Listed on: %d RBL(s)\n", len(listings))
		}
	} else {
		fmt.Printf("\nIP %s is clean - not listed on any RBLs\n", ip)
	}
}
