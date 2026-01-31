// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// Cert represents the Certificate monitoring API namespace.
//
// Use this to monitor TLS/SSL certificates for expiration, validity,
// and configuration issues across HTTPS, SMTPS, IMAPS, and other protocols.
type Cert struct {
	handler *RequestHandler
}

// Errors returns certificate error operations.
//
// Use this to retrieve active certificate errors including expiration
// warnings, validation failures, and configuration issues.
//
// Example:
//
//	errors, err := client.Cert().Errors().Get()
func (c *Cert) Errors() *CertErrors {
	return &CertErrors{handler: c.handler}
}

// Monitors returns operations for managing certificate monitors.
//
// Certificate monitors continuously check TLS/SSL certificates on
// specific hosts and ports for expiration and validity issues.
//
// Example:
//
//	monitors, err := client.Cert().Monitors().Get()
func (c *Cert) Monitors() *CertMonitors {
	return &CertMonitors{handler: c.handler}
}

// Profiles returns operations for managing certificate monitoring profiles.
//
// Certificate monitoring profiles define alert thresholds and validation
// rules for certificate checks. Profiles can be assigned to multiple monitors.
//
// Example:
//
//	profiles, err := client.Cert().Profiles().Get()
func (c *Cert) Profiles() *CertProfiles {
	return &CertProfiles{handler: c.handler}
}

// CertErrors handles certificate error operations.
//
// Certificate errors include expiration warnings, validation failures,
// and configuration issues discovered during monitoring checks.
// This resource is read-only (GET operations only).
type CertErrors struct {
	handler *RequestHandler
}

// Get retrieves current certificate errors.
//
// Usage patterns:
//   - Get(): Retrieves all errors (first page)
//   - Get(map[string]interface{}{"status": "active"}): Retrieves filtered errors
//
// For automatic pagination across all pages, use GetAll() instead.
//
// The params map supports filtering options:
//   - status: Filter by error status
//   - monitor_id: Filter by monitor SID
//   - page: Page number (default: 1)
//   - page_size: Results per page (default: 20)
//
// Example:
//
//	// Get all errors
//	errors, err := client.Cert().Errors().Get()
//
//	// Get errors for specific monitor
//	errors, err := client.Cert().Errors().Get(map[string]interface{}{
//	    "monitor_id": "CM62944aeeee2b46d7a28221164f38976a",
//	})
func (e *CertErrors) Get(params ...interface{}) (map[string]interface{}, error) {
	if len(params) == 0 {
		return e.handler.Get("cert/errors", nil)
	}
	// Check if it's a map[string]interface{} (parameters)
	if p, ok := params[0].(map[string]interface{}); ok {
		return e.handler.Get("cert/errors", p)
	}
	return e.handler.Get("cert/errors", nil)
}

// GetAll retrieves all certificate errors with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of certificate errors without manual pagination handling.
//
// Example:
//
//	errors, err := client.Cert().Errors().GetAll(nil)
//	fmt.Printf("Total errors: %d\n", len(errors))
func (e *CertErrors) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := e.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract errors from response
		errors := []interface{}{}
		if errorList, ok := response["errors"].([]interface{}); ok {
			errors = errorList
		}

		allItems = append(allItems, errors...)

		// Check if there are more pages
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(errors) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// CertMonitors handles certificate monitor management operations.
//
// Certificate monitors continuously check TLS/SSL certificates on specific
// hosts and ports. Supports HTTPS, SMTPS, IMAPS, LDAPS, and other protocols.
// Provides full CRUD operations plus pause/resume functionality.
type CertMonitors struct {
	handler *RequestHandler
}

// Get retrieves certificate monitors with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all monitors (first page)
//   - Get(map[string]interface{}{"status": "active"}): Retrieves filtered monitors
//   - Get("CMxxxxxxxxxx"): Retrieves a single monitor by ID
//   - Get("CMxxxxxxxxxx", "CMyyyyyyyyyy"): Retrieves multiple monitors by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// The params map supports filtering options:
//   - status: Filter by status ("active", "paused", etc.)
//   - protocol: Filter by protocol ("https", "smtps", etc.)
//   - page: Page number (default: 1)
//   - page_size: Results per page (default: 20)
//
// Example:
//
//	// Get all monitors
//	monitors, err := client.Cert().Monitors().Get()
//
//	// Get active HTTPS monitors
//	monitors, err := client.Cert().Monitors().Get(map[string]interface{}{
//	    "status": "active",
//	    "protocol": "https",
//	})
//
//	// Get single monitor
//	monitor, err := client.Cert().Monitors().Get("CM62944aeeee2b46d7a28221164f38976a")
func (m *CertMonitors) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return m.handler.Get("cert/monitors", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return m.handler.Get("cert/monitors", params)
		}
		return m.handler.Get(fmt.Sprintf("cert/monitors/%v", ids[0]), nil)
	}
	// Multiple IDs - convert to comma-separated string
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return m.handler.Get("cert/monitors", params)
}

// GetAll retrieves all certificate monitors with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of certificate monitors without manual pagination handling.
//
// Example:
//
//	monitors, err := client.Cert().Monitors().GetAll(nil)
//	fmt.Printf("Total monitors: %d\n", len(monitors))
func (m *CertMonitors) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := m.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract monitors from response
		monitors := []interface{}{}
		if monitorList, ok := response["monitors"].([]interface{}); ok {
			monitors = monitorList
		}

		allItems = append(allItems, monitors...)

		// Check if there are more pages
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(monitors) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new certificate monitor.
//
// Required parameters:
//   - hostname: Domain name or IP address to monitor
//   - protocol: Protocol to use ("https", "smtps", "imaps", "ldaps", "mysql", "postgresql", "sips")
//   - port: Port number (e.g., 443 for HTTPS, 587 for SMTPS)
//
// Optional parameters:
//   - name: Descriptive name for the monitor
//   - cert_profile: Certificate monitoring profile SID (e.g., "CP79b597e61a984a35b5eb7dcdbc3de53c")
//   - contact_group: Contact group SID for alerts (e.g., "CG4f3e2d1c0b9a8776655443322110fedc")
//
// Example:
//
//	monitor, err := client.Cert().Monitors().Create(map[string]interface{}{
//	    "name": "Production HTTPS",
//	    "hostname": "www.example.com",
//	    "protocol": "https",
//	    "port": 443,
//	    "cert_profile": "CP79b597e61a984a35b5eb7dcdbc3de53c",
//	})
func (m *CertMonitors) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return m.handler.Post("cert/monitors", params)
}

// Update updates a certificate monitor.
//
// The id parameter should be the monitor SID (e.g., "CM62944aeeee2b46d7a28221164f38976a").
//
// Updatable fields:
//   - name: Update the monitor name
//   - hostname: Update the hostname
//   - port: Update the port number
//   - cert_profile: Change the monitoring profile
//   - contact_group: Change the contact group
//
// Example:
//
//	_, err := client.Cert().Monitors().Update("CM62944aeeee2b46d7a28221164f38976a", map[string]interface{}{
//	    "name": "Updated Monitor Name",
//	})
func (m *CertMonitors) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return m.handler.Put(fmt.Sprintf("cert/monitors/%v", id), params)
}

// Delete deletes a certificate monitor.
//
// The id parameter should be the monitor SID (e.g., "CM62944aeeee2b46d7a28221164f38976a").
// This permanently removes the monitor and stops certificate checking.
//
// Example:
//
//	_, err := client.Cert().Monitors().Delete("CM62944aeeee2b46d7a28221164f38976a")
func (m *CertMonitors) Delete(id interface{}) (map[string]interface{}, error) {
	return m.handler.Delete(fmt.Sprintf("cert/monitors/%v", id))
}

// Pause temporarily pauses monitoring for a certificate.
//
// The id parameter should be the monitor SID (e.g., "CM62944aeeee2b46d7a28221164f38976a").
// Pausing stops checks without deleting the monitor. Use Resume() to restart.
//
// Example:
//
//	_, err := client.Cert().Monitors().Pause("CM62944aeeee2b46d7a28221164f38976a")
func (m *CertMonitors) Pause(id interface{}) (map[string]interface{}, error) {
	return m.handler.Post(fmt.Sprintf("cert/monitors/%v/pause", id), nil)
}

// Resume resumes monitoring for a paused certificate.
//
// The id parameter should be the monitor SID (e.g., "CM62944aeeee2b46d7a28221164f38976a").
// This restarts checks on a previously paused monitor.
//
// Example:
//
//	_, err := client.Cert().Monitors().Resume("CM62944aeeee2b46d7a28221164f38976a")
func (m *CertMonitors) Resume(id interface{}) (map[string]interface{}, error) {
	return m.handler.Post(fmt.Sprintf("cert/monitors/%v/resume", id), nil)
}

// CertProfiles handles certificate monitoring profile operations.
//
// Certificate monitoring profiles define alert thresholds (days until expiration)
// and validation rules for certificate checks. Profiles can be assigned to
// multiple monitors to control their monitoring behavior.
type CertProfiles struct {
	handler *RequestHandler
}

// Get retrieves certificate monitoring profiles with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all profiles (first page)
//   - Get(map[string]interface{}{"page": 2}): Retrieves with pagination
//   - Get("CPxxxxxxxxxx"): Retrieves a single profile by ID
//   - Get("CPxxxxxxxxxx", "CPyyyyyyyyyy"): Retrieves multiple profiles by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// Example:
//
//	// Get all profiles
//	profiles, err := client.Cert().Profiles().Get()
//
//	// Get single profile
//	profile, err := client.Cert().Profiles().Get("CP79b597e61a984a35b5eb7dcdbc3de53c")
func (p *CertProfiles) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return p.handler.Get("cert/profiles", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return p.handler.Get("cert/profiles", params)
		}
		return p.handler.Get(fmt.Sprintf("cert/profiles/%v", ids[0]), nil)
	}
	// Multiple IDs
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return p.handler.Get("cert/profiles", params)
}

// GetAll retrieves all certificate monitoring profiles with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of certificate monitoring profiles without manual pagination handling.
//
// Example:
//
//	profiles, err := client.Cert().Profiles().GetAll(nil)
//	fmt.Printf("Total profiles: %d\n", len(profiles))
func (p *CertProfiles) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := p.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract profiles from response
		profiles := []interface{}{}
		if profileList, ok := response["profiles"].([]interface{}); ok {
			profiles = profileList
		}

		allItems = append(allItems, profiles...)

		// Check if there are more pages
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(profiles) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new certificate monitoring profile.
//
// Required parameters:
//   - name: Descriptive name for the profile
//   - expiration_thresholds: Array of day thresholds for alerts (e.g., [30, 14, 7])
//
// Optional parameters:
//   - alert_on_invalid: Alert on validation failures (boolean, default: true)
//   - alert_on_revoked: Alert on certificate revocation (boolean, default: true)
//
// Example:
//
//	profile, err := client.Cert().Profiles().Create(map[string]interface{}{
//	    "name": "Standard Certificate Monitoring",
//	    "expiration_thresholds": []int{30, 14, 7, 1},
//	    "alert_on_invalid": true,
//	})
func (p *CertProfiles) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Post("cert/profiles", params)
}

// Update updates a certificate monitoring profile.
//
// The id parameter should be the profile SID (e.g., "CP79b597e61a984a35b5eb7dcdbc3de53c").
//
// Updatable fields:
//   - name: Update the profile name
//   - expiration_thresholds: Update the day thresholds for alerts
//   - alert_on_invalid: Update invalid certificate alert setting
//   - alert_on_revoked: Update revoked certificate alert setting
//
// Example:
//
//	_, err := client.Cert().Profiles().Update("CP79b597e61a984a35b5eb7dcdbc3de53c", map[string]interface{}{
//	    "name": "Updated Profile Name",
//	    "expiration_thresholds": []int{60, 30, 14},
//	})
func (p *CertProfiles) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Put(fmt.Sprintf("cert/profiles/%v", id), params)
}

// Delete deletes a certificate monitoring profile.
//
// The id parameter should be the profile SID (e.g., "CP79b597e61a984a35b5eb7dcdbc3de53c").
// Note: Cannot delete a profile that is currently assigned to monitors.
//
// Example:
//
//	_, err := client.Cert().Profiles().Delete("CP79b597e61a984a35b5eb7dcdbc3de53c")
func (p *CertProfiles) Delete(id interface{}) (map[string]interface{}, error) {
	return p.handler.Delete(fmt.Sprintf("cert/profiles/%v", id))
}
