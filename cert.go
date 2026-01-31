// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// Cert represents the Certificate monitoring API namespace
type Cert struct {
	handler *RequestHandler
}

// Errors returns current certificate errors
func (c *Cert) Errors() *CertErrors {
	return &CertErrors{handler: c.handler}
}

// Monitors returns operations for managing certificate monitors
func (c *Cert) Monitors() *CertMonitors {
	return &CertMonitors{handler: c.handler}
}

// Profiles returns operations for managing certificate monitoring profiles
func (c *Cert) Profiles() *CertProfiles {
	return &CertProfiles{handler: c.handler}
}

// CertErrors handles certificate error operations
type CertErrors struct {
	handler *RequestHandler
}

// Get retrieves certificate errors
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

// GetAll retrieves all errors with automatic pagination
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

// CertMonitors handles certificate monitor management operations
type CertMonitors struct {
	handler *RequestHandler
}

// Get retrieves monitors (all, by ID, or by array of IDs)
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

// GetAll retrieves all monitors with automatic pagination
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

// Create creates a new certificate monitor
func (m *CertMonitors) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return m.handler.Post("cert/monitors", params)
}

// Update updates a certificate monitor
func (m *CertMonitors) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return m.handler.Put(fmt.Sprintf("cert/monitors/%v", id), params)
}

// Delete deletes a certificate monitor
func (m *CertMonitors) Delete(id interface{}) (map[string]interface{}, error) {
	return m.handler.Delete(fmt.Sprintf("cert/monitors/%v", id))
}

// Pause pauses monitoring for a certificate
func (m *CertMonitors) Pause(id interface{}) (map[string]interface{}, error) {
	return m.handler.Post(fmt.Sprintf("cert/monitors/%v/pause", id), nil)
}

// Resume resumes monitoring for a certificate
func (m *CertMonitors) Resume(id interface{}) (map[string]interface{}, error) {
	return m.handler.Post(fmt.Sprintf("cert/monitors/%v/resume", id), nil)
}

// CertProfiles handles certificate monitoring profile operations
type CertProfiles struct {
	handler *RequestHandler
}

// Get retrieves profiles (all, by ID, or by array of IDs)
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

// GetAll retrieves all profiles with automatic pagination
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

// Create creates a new certificate monitoring profile
func (p *CertProfiles) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Post("cert/profiles", params)
}

// Update updates a certificate monitoring profile
func (p *CertProfiles) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return p.handler.Put(fmt.Sprintf("cert/profiles/%v", id), params)
}

// Delete deletes a certificate monitoring profile
func (p *CertProfiles) Delete(id interface{}) (map[string]interface{}, error) {
	return p.handler.Delete(fmt.Sprintf("cert/profiles/%v", id))
}
