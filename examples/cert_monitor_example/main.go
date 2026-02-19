// Example: Certificate monitoring - list errors, manage monitors and profiles
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

	// ===================================================================
	// Certificate Errors
	// ===================================================================
	fmt.Println("=== Listing Certificate Errors ===")
	errors, err := client.Cert().Errors().Get()
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	if errorList, ok := errors["errors"].([]interface{}); ok {
		fmt.Printf("Total errors: %d\n\n", len(errorList))
		for _, e := range errorList {
			if errorData, ok := e.(map[string]interface{}); ok {
				fmt.Printf("Error ID: %v\n", errorData["id"])
				fmt.Printf("  Monitor: %v\n", errorData["monitor_name"])
				fmt.Printf("  Type: %v\n", errorData["error_type"])
				fmt.Printf("  Message: %v\n\n", errorData["message"])
			}
		}
	}

	// ===================================================================
	// Certificate Profiles
	// ===================================================================
	fmt.Println("=== Managing Certificate Profiles ===")

	// List all profiles
	profiles, err := client.Cert().Profiles().Get()
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	if profileList, ok := profiles["profiles"].([]interface{}); ok {
		fmt.Printf("Total profiles: %d\n", len(profileList))
	}

	// Create a new profile
	fmt.Println("\n=== Creating a new certificate profile ===")
	newProfile, err := client.Cert().Profiles().Create(map[string]interface{}{
		"name":                     "Example Certificate Profile",
		"expiration_warning_days":  30,
		"expiration_critical_days": 7,
		"check_self_signed":        true,
		"check_hostname_mismatch":  true,
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	var profileID string
	if profileData, ok := newProfile["profile"].(map[string]interface{}); ok {
		profileID = fmt.Sprintf("%v", profileData["id"])
		fmt.Printf("Created profile ID: %s\n", profileID)
	}

	// Get specific profile
	fmt.Println("\n=== Getting specific profile ===")
	profile, err := client.Cert().Profiles().Get(profileID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	if profileData, ok := profile["profile"].(map[string]interface{}); ok {
		fmt.Printf("Profile name: %v\n", profileData["name"])
		fmt.Printf("Expiration warning days: %v\n", profileData["expiration_warning_days"])
	}

	// Update profile
	fmt.Println("\n=== Updating profile ===")
	_, err = client.Cert().Profiles().Update(profileID, map[string]interface{}{
		"expiration_warning_days": 45,
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Println("Updated profile warning days to 45")

	// ===================================================================
	// Certificate Monitors
	// ===================================================================
	fmt.Println("\n=== Managing Certificate Monitors ===")

	// List all monitors
	monitors, err := client.Cert().Monitors().Get()
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	if monitorList, ok := monitors["monitors"].([]interface{}); ok {
		fmt.Printf("Total monitors: %d\n", len(monitorList))
	}

	// Create a new HTTPS monitor
	fmt.Println("\n=== Creating HTTPS certificate monitor ===")
	httpsMonitor, err := client.Cert().Monitors().Create(map[string]interface{}{
		"name":          "Example HTTPS Monitor",
		"hostname":      "example.com",
		"protocol":      "https",
		"profile":       profileID,
		"contact_group": "CG4f3e2d1c0b9a8776655443322110fed", // Use your contact group ID
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	var httpsMonitorID string
	if monitorData, ok := httpsMonitor["monitor"].(map[string]interface{}); ok {
		httpsMonitorID = fmt.Sprintf("%v", monitorData["id"])
		fmt.Printf("Created HTTPS monitor ID: %s\n", httpsMonitorID)
	}

	// Create a mail server monitor (SMTPS)
	fmt.Println("\n=== Creating SMTPS certificate monitor ===")
	smtpsMonitor, err := client.Cert().Monitors().Create(map[string]interface{}{
		"name":          "Example Mail Server Monitor",
		"hostname":      "mail.example.com",
		"protocol":      "smtps",
		"profile":       profileID,
		"contact_group": "CG4f3e2d1c0b9a8776655443322110fed",
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}

	var smtpsMonitorID string
	if monitorData, ok := smtpsMonitor["monitor"].(map[string]interface{}); ok {
		smtpsMonitorID = fmt.Sprintf("%v", monitorData["id"])
		fmt.Printf("Created SMTPS monitor ID: %s\n", smtpsMonitorID)
	}

	// Get specific monitor
	fmt.Println("\n=== Getting specific monitor ===")
	monitor, err := client.Cert().Monitors().Get(httpsMonitorID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	if monitorData, ok := monitor["monitor"].(map[string]interface{}); ok {
		fmt.Printf("Monitor name: %v\n", monitorData["name"])
		fmt.Printf("Hostname: %v\n", monitorData["hostname"])
		fmt.Printf("Protocol: %v\n", monitorData["protocol"])
		fmt.Printf("Status: %v\n", monitorData["status"])
	}

	// Update monitor
	fmt.Println("\n=== Updating monitor ===")
	_, err = client.Cert().Monitors().Update(httpsMonitorID, map[string]interface{}{
		"name": "Updated HTTPS Monitor Name",
	})
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Println("Updated monitor name")

	// Pause monitoring
	fmt.Println("\n=== Pausing monitor ===")
	_, err = client.Cert().Monitors().Pause(httpsMonitorID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Paused monitor ID: %s\n", httpsMonitorID)

	// Resume monitoring
	fmt.Println("\n=== Resuming monitor ===")
	_, err = client.Cert().Monitors().Resume(httpsMonitorID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Resumed monitor ID: %s\n", httpsMonitorID)

	// ===================================================================
	// Cleanup
	// ===================================================================
	fmt.Println("\n=== Cleaning up - Deleting created resources ===")

	// Delete monitors
	_, err = client.Cert().Monitors().Delete(httpsMonitorID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Deleted HTTPS monitor ID: %s\n", httpsMonitorID)

	_, err = client.Cert().Monitors().Delete(smtpsMonitorID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Deleted SMTPS monitor ID: %s\n", smtpsMonitorID)

	// Delete profile
	_, err = client.Cert().Profiles().Delete(profileID)
	if err != nil {
		log.Fatalf("API Error: %v", err)
	}
	fmt.Printf("Deleted profile ID: %s\n", profileID)

	fmt.Println("\n=== Certificate Monitoring Example Complete ===")
}
