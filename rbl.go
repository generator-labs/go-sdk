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
