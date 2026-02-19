// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// RBL represents the RBL monitoring API namespace.
//
// Use this to access RBL monitoring operations including host management,
// profile configuration, source management, and real-time checks.
type RBL struct {
	handler *RequestHandler
}

// Hosts returns operations for managing monitored hosts.
//
// Hosts are IP addresses or domains that are continuously monitored
// against configured RBL sources. Use this to create, retrieve, update,
// and delete monitored hosts.
//
// Example:
//
//	hosts, err := client.RBL().Hosts().Get()
func (r *RBL) Hosts() *RBLHosts {
	return &RBLHosts{handler: r.handler}
}

// Profiles returns operations for managing monitoring profiles.
//
// Monitoring profiles define which RBL sources are checked for hosts.
// Profiles can be assigned to multiple hosts to control their monitoring behavior.
//
// Example:
//
//	profiles, err := client.RBL().Profiles().Get()
func (r *RBL) Profiles() *RBLProfiles {
	return &RBLProfiles{handler: r.handler}
}

// Sources returns operations for managing RBL sources.
//
// RBL sources are the blacklist databases that hosts are checked against.
// This includes both system-provided sources and custom sources you've added.
//
// Example:
//
//	sources, err := client.RBL().Sources().Get()
func (r *RBL) Sources() *RBLSources {
	return &RBLSources{handler: r.handler}
}

// Check returns operations for performing manual RBL checks.
//
// Manual checks allow you to check an IP address or host against RBL sources
// without requiring it to be added as a monitored host.
//
// Example:
//
//	result, err := client.RBL().Check().Start(map[string]interface{}{"host": "1.2.3.4"})
func (r *RBL) Check() *RBLCheck {
	return &RBLCheck{handler: r.handler}
}

// Listings retrieves current RBL listings for all monitored hosts.
//
// Returns a list of all hosts that are currently listed on one or more RBLs.
//
// Example:
//
//	listings, err := client.RBL().Listings()
func (r *RBL) Listings() (map[string]interface{}, error) {
	return r.handler.Get("rbl/listings", nil)
}

// RBLCheck handles manual RBL check operations.
//
// Provides methods to start a new check and retrieve check status/results.
type RBLCheck struct {
	handler *RequestHandler
}

// Start initiates a new manual RBL check.
//
// Required parameters:
//   - host: IP address or hostname to check
//
// Optional parameters:
//   - callback: URL to receive results via webhook
//   - details: Set to true for detailed results
//
// Example:
//
//	result, err := client.RBL().Check().Start(map[string]interface{}{
//	    "host": "1.2.3.4",
//	    "details": true,
//	})
func (c *RBLCheck) Start(params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Post("rbl/check/start", params)
}

// Status retrieves the status of a manual RBL check.
//
// The id parameter should be the check ID returned from Start().
//
// Optional parameters:
//   - details: Set to true for detailed results
//
// Example:
//
//	status, err := client.RBL().Check().Status("PP1234567890abcdef...", nil)
func (c *RBLCheck) Status(id string, params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Get(fmt.Sprintf("rbl/check/status/%s", id), params)
}

// RBLHosts handles host management operations.
//
// Provides CRUD operations for monitored hosts including create, retrieve,
// update, and delete. Also provides GetAll() for automatic pagination.
type RBLHosts struct {
	handler *RequestHandler
}

// Get retrieves hosts with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all hosts (first page)
//   - Get(map[string]interface{}{"status": "active"}): Retrieves filtered hosts
//   - Get("HTxxxxxxxxxx"): Retrieves a single host by ID
//   - Get("HTxxxxxxxxxx", "HTyyyyyyyyyy"): Retrieves multiple hosts by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// The params map supports filtering options like:
//   - status: Filter by status ("active", "paused", etc.)
//   - page: Page number (default: 1)
//   - page_size: Results per page (default: 20)
//
// Example:
//
//	// Get all hosts (first page)
//	hosts, err := client.RBL().Hosts().Get()
//
//	// Get filtered hosts
//	hosts, err := client.RBL().Hosts().Get(map[string]interface{}{"status": "active"})
//
//	// Get single host
//	host, err := client.RBL().Hosts().Get("HT1a2b3c4d5e6f7890abcdef1234567890")
//
//	// Get multiple hosts
//	hosts, err := client.RBL().Hosts().Get("HTxxxxxxxx", "HTyyyyyyyy")
func (h *RBLHosts) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return h.handler.Get("rbl/hosts", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return h.handler.Get("rbl/hosts", params)
		}
		return h.handler.Get(fmt.Sprintf("rbl/hosts/%v", ids[0]), nil)
	}
	// Multiple IDs - convert to comma-separated string
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return h.handler.Get("rbl/hosts", params)
}

// GetAll retrieves all hosts with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of hosts. Use this when you need to retrieve all hosts without manually
// handling pagination.
//
// The params map supports the same filtering options as Get():
//   - status: Filter by status ("active", "paused", etc.)
//   - page_size: Results per page (default: 20, affects API call count)
//
// Example:
//
//	hosts, err := client.RBL().Hosts().GetAll(map[string]interface{}{"status": "active"})
//	fmt.Printf("Total hosts: %d\n", len(hosts))
func (h *RBLHosts) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := h.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract hosts from response
		hosts := []interface{}{}
		if hostList, ok := response["hosts"].([]interface{}); ok {
			hosts = hostList
		}

		allItems = append(allItems, hosts...)

		// Check if there are more pages
		totalPages := 1.0
		if tp, ok := response["total_pages"].(float64); ok {
			totalPages = tp
		}

		if float64(page) >= totalPages || len(hosts) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new monitored host.
//
// Required parameters:
//   - name: Descriptive name for the host
//   - host: IP address, domain, or CIDR range
//   - type: Host type ("rbl", "uribl", or "ip_range")
//
// Optional parameters:
//   - profile: Monitoring profile ID (e.g., "RP9f8e7d6c5b4a3210fedcba0987654321")
//   - contact_group: Contact group ID (e.g., "CG4f3e2d1c0b9a8776655443322110fedc")
//
// Example:
//
//	host, err := client.RBL().Hosts().Create(map[string]interface{}{
//	    "name": "Production Mail Server",
//	    "host": "192.168.1.100",
//	    "type": "rbl",
//	    "profile": "RP9f8e7d6c5b4a3210fedcba0987654321",
//	})
func (h *RBLHosts) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return h.handler.Post("rbl/hosts", params)
}

// Update updates a monitored host.
//
// The id parameter should be the host SID (e.g., "HT1a2b3c4d5e6f7890abcdef1234567890").
//
// Updatable fields:
//   - name: Update the host name
//   - profile: Change the monitoring profile
//   - contact_group: Change the contact group
//
// Example:
//
//	_, err := client.RBL().Hosts().Update("HT1a2b3c4d5e6f7890abcdef1234567890", map[string]interface{}{
//	    "name": "Updated Server Name",
//	})
func (h *RBLHosts) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return h.handler.Put(fmt.Sprintf("rbl/hosts/%v", id), params)
}

// Delete deletes a monitored host.
//
// The id parameter should be the host SID (e.g., "HT1a2b3c4d5e6f7890abcdef1234567890").
// This permanently removes the host from monitoring.
//
// Example:
//
//	_, err := client.RBL().Hosts().Delete("HT1a2b3c4d5e6f7890abcdef1234567890")
func (h *RBLHosts) Delete(id interface{}) (map[string]interface{}, error) {
	return h.handler.Delete(fmt.Sprintf("rbl/hosts/%v", id))
}

// Pause temporarily pauses monitoring for a host.
//
// The id parameter should be the host SID (e.g., "HT1a2b3c4d5e6f7890abcdef1234567890").
// Pausing stops checks without deleting the host. Use Resume() to restart.
//
// Example:
//
//	_, err := client.RBL().Hosts().Pause("HT1a2b3c4d5e6f7890abcdef1234567890")
func (h *RBLHosts) Pause(id interface{}) (map[string]interface{}, error) {
	return h.handler.Post(fmt.Sprintf("rbl/hosts/%v/pause", id), nil)
}

// Resume resumes monitoring for a paused host.
//
// The id parameter should be the host SID (e.g., "HT1a2b3c4d5e6f7890abcdef1234567890").
// This restarts checks on a previously paused host.
//
// Example:
//
//	_, err := client.RBL().Hosts().Resume("HT1a2b3c4d5e6f7890abcdef1234567890")
func (h *RBLHosts) Resume(id interface{}) (map[string]interface{}, error) {
	return h.handler.Post(fmt.Sprintf("rbl/hosts/%v/resume", id), nil)
}

// RBLProfiles handles monitoring profile operations.
//
// Monitoring profiles define which RBL sources are checked for hosts.
// Provides CRUD operations including create, retrieve, update, and delete.
type RBLProfiles struct {
	handler *RequestHandler
}

// Get retrieves monitoring profiles with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all profiles (first page)
//   - Get(map[string]interface{}{"page": 2}): Retrieves with pagination
//   - Get("RPxxxxxxxxxx"): Retrieves a single profile by ID
//   - Get("RPxxxxxxxxxx", "RPyyyyyyyyyy"): Retrieves multiple profiles by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// Example:
//
//	// Get all profiles
//	profiles, err := client.RBL().Profiles().Get()
//
//	// Get single profile
//	profile, err := client.RBL().Profiles().Get("RP9f8e7d6c5b4a3210fedcba0987654321")
func (p *RBLProfiles) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return p.handler.Get("rbl/profiles", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return p.handler.Get("rbl/profiles", params)
		}
		return p.handler.Get(fmt.Sprintf("rbl/profiles/%v", ids[0]), nil)
	}
	// Multiple IDs
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return p.handler.Get("rbl/profiles", params)
}

// GetAll retrieves all monitoring profiles with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of profiles without manual pagination handling.
//
// Example:
//
//	profiles, err := client.RBL().Profiles().GetAll(nil)
//	fmt.Printf("Total profiles: %d\n", len(profiles))
func (p *RBLProfiles) GetAll(params map[string]interface{}) ([]interface{}, error) {
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
		totalPages := 1.0
		if tp, ok := response["total_pages"].(float64); ok {
			totalPages = tp
		}

		if float64(page) >= totalPages || len(profiles) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new monitoring profile.
//
// A monitoring profile defines which RBL sources are checked for hosts.
// Profiles can be assigned to multiple hosts.
//
// Required parameters:
//   - name: Descriptive name for the profile
//   - entries: Array of RBL source IDs to include in checks
//
// Example:
//
//	profile, err := client.RBL().Profiles().Create(map[string]interface{}{
//	    "name": "Standard RBL Check",
//	    "entries": []string{"RB1234567890", "RB0987654321"},
//	})
func (p *RBLProfiles) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Post("rbl/profiles", params)
}

// Update updates a monitoring profile.
//
// The id parameter should be the profile SID (e.g., "RP9f8e7d6c5b4a3210fedcba0987654321").
//
// Updatable fields:
//   - name: Update the profile name
//   - entries: Update the list of RBL sources
//
// Example:
//
//	_, err := client.RBL().Profiles().Update("RP9f8e7d6c5b4a3210fedcba0987654321", map[string]interface{}{
//	    "name": "Updated Profile Name",
//	})
func (p *RBLProfiles) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Put(fmt.Sprintf("rbl/profiles/%v", id), params)
}

// Delete deletes a monitoring profile.
//
// The id parameter should be the profile SID (e.g., "RP9f8e7d6c5b4a3210fedcba0987654321").
// Note: You cannot delete a profile that is currently assigned to hosts.
//
// Example:
//
//	_, err := client.RBL().Profiles().Delete("RP9f8e7d6c5b4a3210fedcba0987654321")
func (p *RBLProfiles) Delete(id interface{}) (map[string]interface{}, error) {
	return p.handler.Delete(fmt.Sprintf("rbl/profiles/%v", id))
}

// RBLSources handles RBL source operations.
//
// RBL sources are the blacklist databases that hosts are checked against.
// This includes both system-provided sources and custom sources.
// Provides CRUD operations and pause/resume functionality.
type RBLSources struct {
	handler *RequestHandler
}

// Get retrieves RBL sources with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all sources (first page)
//   - Get(map[string]interface{}{"status": "active"}): Retrieves filtered sources
//   - Get("RBxxxxxxxxxx"): Retrieves a single source by ID
//   - Get("RBxxxxxxxxxx", "RByyyyyyyyyy"): Retrieves multiple sources by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// Example:
//
//	// Get all sources
//	sources, err := client.RBL().Sources().Get()
//
//	// Get single source
//	source, err := client.RBL().Sources().Get("RB1234567890abcdef1234567890abcdef")
func (s *RBLSources) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return s.handler.Get("rbl/sources", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return s.handler.Get("rbl/sources", params)
		}
		return s.handler.Get(fmt.Sprintf("rbl/sources/%v", ids[0]), nil)
	}
	// Multiple IDs
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return s.handler.Get("rbl/sources", params)
}

// GetAll retrieves all RBL sources with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of RBL sources without manual pagination handling.
//
// Example:
//
//	sources, err := client.RBL().Sources().GetAll(nil)
//	fmt.Printf("Total sources: %d\n", len(sources))
func (s *RBLSources) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := s.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract sources from response
		sources := []interface{}{}
		if sourceList, ok := response["sources"].([]interface{}); ok {
			sources = sourceList
		}

		allItems = append(allItems, sources...)

		// Check if there are more pages
		totalPages := 1.0
		if tp, ok := response["total_pages"].(float64); ok {
			totalPages = tp
		}

		if float64(page) >= totalPages || len(sources) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new custom RBL source.
//
// Custom RBL sources allow you to add your own private blacklists.
//
// Required parameters:
//   - name: Descriptive name for the source
//   - host: DNS hostname of the RBL (e.g., "dnsbl.example.com")
//
// Optional parameters:
//   - description: Detailed description of the source
//
// Example:
//
//	source, err := client.RBL().Sources().Create(map[string]interface{}{
//	    "name": "Private DNSBL",
//	    "host": "private-dnsbl.example.com",
//	})
func (s *RBLSources) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return s.handler.Post("rbl/sources", params)
}

// Update updates a custom RBL source.
//
// The id parameter should be the source SID (e.g., "RB1234567890abcdef1234567890abcdef").
// Note: Only custom sources can be updated, system sources cannot be modified.
//
// Updatable fields:
//   - name: Update the source name
//   - host: Update the DNS hostname
//   - description: Update the description
//
// Example:
//
//	_, err := client.RBL().Sources().Update("RB1234567890abcdef1234567890abcdef", map[string]interface{}{
//	    "name": "Updated Source Name",
//	})
func (s *RBLSources) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return s.handler.Put(fmt.Sprintf("rbl/sources/%v", id), params)
}

// Delete deletes a custom RBL source.
//
// The id parameter should be the source SID (e.g., "RB1234567890abcdef1234567890abcdef").
// Note: Only custom sources can be deleted, system sources cannot be removed.
// Cannot delete a source that is currently used in monitoring profiles.
//
// Example:
//
//	_, err := client.RBL().Sources().Delete("RB1234567890abcdef1234567890abcdef")
func (s *RBLSources) Delete(id interface{}) (map[string]interface{}, error) {
	return s.handler.Delete(fmt.Sprintf("rbl/sources/%v", id))
}

// Pause temporarily pauses an RBL source.
//
// The id parameter should be the source SID (e.g., "RB1234567890abcdef1234567890abcdef").
// Pausing stops checks against this source. Use Resume() to restart.
//
// Example:
//
//	_, err := client.RBL().Sources().Pause("RB1234567890abcdef1234567890abcdef")
func (s *RBLSources) Pause(id interface{}) (map[string]interface{}, error) {
	return s.handler.Post(fmt.Sprintf("rbl/sources/%v/pause", id), nil)
}

// Resume resumes a paused RBL source.
//
// The id parameter should be the source SID (e.g., "RB1234567890abcdef1234567890abcdef").
// This restarts checks against a previously paused source.
//
// Example:
//
//	_, err := client.RBL().Sources().Resume("RB1234567890abcdef1234567890abcdef")
func (s *RBLSources) Resume(id interface{}) (map[string]interface{}, error) {
	return s.handler.Post(fmt.Sprintf("rbl/sources/%v/resume", id), nil)
}
