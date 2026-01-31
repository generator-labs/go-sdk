// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// Contact represents the Contact management API namespace
type Contact struct {
	handler *RequestHandler
}

// Contacts returns operations for managing contacts
func (c *Contact) Contacts() *Contacts {
	return &Contacts{handler: c.handler}
}

// Groups returns operations for managing contact groups
func (c *Contact) Groups() *Groups {
	return &Groups{handler: c.handler}
}

// Contacts handles contact management operations
type Contacts struct {
	handler *RequestHandler
}

// Get retrieves contacts (all, by ID, or by array of IDs)
func (c *Contacts) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return c.handler.Get("contact/contacts", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return c.handler.Get("contact/contacts", params)
		}
		return c.handler.Get(fmt.Sprintf("contact/contacts/%v", ids[0]), nil)
	}
	// Multiple IDs
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return c.handler.Get("contact/contacts", params)
}

// GetAll retrieves all contacts with automatic pagination
func (c *Contacts) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := c.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract contacts from response
		contacts := []interface{}{}
		if contactList, ok := response["contacts"].([]interface{}); ok {
			contacts = contactList
		}

		allItems = append(allItems, contacts...)

		// Check if there are more pages
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(contacts) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new contact
func (c *Contacts) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Post("contact/contacts", params)
}

// Update updates a contact
func (c *Contacts) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Put(fmt.Sprintf("contact/contacts/%v", id), params)
}

// Delete deletes a contact
func (c *Contacts) Delete(id interface{}) (map[string]interface{}, error) {
	return c.handler.Delete(fmt.Sprintf("contact/contacts/%v", id))
}

// Groups handles contact group management operations
type Groups struct {
	handler *RequestHandler
}

// Get retrieves groups (all, by ID, or by array of IDs)
func (g *Groups) Get(ids ...interface{}) (map[string]interface{}, error) {
	if len(ids) == 0 {
		return g.handler.Get("contact/groups", nil)
	}
	if len(ids) == 1 {
		// Check if it's a map[string]interface{} (parameters)
		if params, ok := ids[0].(map[string]interface{}); ok {
			return g.handler.Get("contact/groups", params)
		}
		return g.handler.Get(fmt.Sprintf("contact/groups/%v", ids[0]), nil)
	}
	// Multiple IDs
	idStr := fmt.Sprintf("%v", ids[0])
	for i := 1; i < len(ids); i++ {
		idStr = fmt.Sprintf("%s,%v", idStr, ids[i])
	}
	params := map[string]interface{}{"ids": idStr}
	return g.handler.Get("contact/groups", params)
}

// GetAll retrieves all groups with automatic pagination
func (g *Groups) GetAll(params map[string]interface{}) ([]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	allItems := []interface{}{}
	page := 1

	for {
		params["page"] = page
		response, err := g.Get(params)
		if err != nil {
			return nil, err
		}

		// Extract groups from response
		groups := []interface{}{}
		if groupList, ok := response["groups"].([]interface{}); ok {
			groups = groupList
		}

		allItems = append(allItems, groups...)

		// Check if there are more pages
		hasMore := false
		if hm, ok := response["has_more"].(bool); ok {
			hasMore = hm
		}

		if !hasMore || len(groups) == 0 {
			break
		}

		page++
	}

	return allItems, nil
}

// Create creates a new contact group
func (g *Groups) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return g.handler.Post("contact/groups", params)
}

// Update updates a contact group
func (g *Groups) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return g.handler.Put(fmt.Sprintf("contact/groups/%v", id), params)
}

// Delete deletes a contact group
func (g *Groups) Delete(id interface{}) (map[string]interface{}, error) {
	return g.handler.Delete(fmt.Sprintf("contact/groups/%v", id))
}
