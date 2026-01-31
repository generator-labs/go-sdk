// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// RBL represents the RBL monitoring API namespace
type RBL struct {
	handler *RequestHandler
}

// Hosts returns operations for managing monitored hosts
func (r *RBL) Hosts() *RBLHosts {
	return &RBLHosts{handler: r.handler}
}

// Profiles returns operations for managing monitoring profiles
func (r *RBL) Profiles() *RBLProfiles {
	return &RBLProfiles{handler: r.handler}
}

// Sources returns operations for managing RBL sources
func (r *RBL) Sources() *RBLSources {
	return &RBLSources{handler: r.handler}
}

// Check checks if an IP address is listed on any RBLs
func (r *RBL) Check(ip string) (map[string]interface{}, error) {
	params := map[string]interface{}{"ip": ip}
	return r.handler.Get("rbl/check", params)
}

// Listings returns current RBL listings for monitored hosts
func (r *RBL) Listings() (map[string]interface{}, error) {
	return r.handler.Get("rbl/listings", nil)
}

// RBLHosts handles host management operations
type RBLHosts struct {
	handler *RequestHandler
}

// Get retrieves hosts (all, by ID, or by array of IDs)
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

// GetAll retrieves all hosts with automatic pagination
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
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(hosts) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new monitored host
func (h *RBLHosts) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return h.handler.Post("rbl/hosts", params)
}

// Update updates a monitored host
func (h *RBLHosts) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return h.handler.Put(fmt.Sprintf("rbl/hosts/%v", id), params)
}

// Delete deletes a monitored host
func (h *RBLHosts) Delete(id interface{}) (map[string]interface{}, error) {
	return h.handler.Delete(fmt.Sprintf("rbl/hosts/%v", id))
}

// RBLProfiles handles monitoring profile operations
type RBLProfiles struct {
	handler *RequestHandler
}

// Get retrieves profiles (all, by ID, or by array of IDs)
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

// GetAll retrieves all profiles with automatic pagination
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

// Create creates a new monitoring profile
func (p *RBLProfiles) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Post("rbl/profiles", params)
}

// Update updates a monitoring profile
func (p *RBLProfiles) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Put(fmt.Sprintf("rbl/profiles/%v", id), params)
}

// Delete deletes a monitoring profile
func (p *RBLProfiles) Delete(id interface{}) (map[string]interface{}, error) {
	return p.handler.Delete(fmt.Sprintf("rbl/profiles/%v", id))
}

// RBLSources handles RBL source operations
type RBLSources struct {
	handler *RequestHandler
}

// Get retrieves RBL sources (all, by ID, or by array of IDs)
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

// GetAll retrieves all sources with automatic pagination
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
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(sources) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new RBL source
func (s *RBLSources) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return s.handler.Post("rbl/sources", params)
}

// Update updates an RBL source
func (s *RBLSources) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return s.handler.Put(fmt.Sprintf("rbl/sources/%v", id), params)
}

// Delete deletes an RBL source
func (s *RBLSources) Delete(id interface{}) (map[string]interface{}, error) {
	return s.handler.Delete(fmt.Sprintf("rbl/sources/%v", id))
}
