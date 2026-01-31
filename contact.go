// This file is part of the Generator Labs Go SDK package.
//
// (c) Generator Labs <support@generatorlabs.com>
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package generatorlabs

import "fmt"

// Contact represents the Contact management API namespace.
//
// Use this to manage contacts and contact groups for receiving alerts
// and notifications from RBL and certificate monitoring.
type Contact struct {
	handler *RequestHandler
}

// Contacts returns operations for managing individual contacts.
//
// Contacts are email addresses or phone numbers that receive alerts
// when monitoring events occur. Contacts can be organized into groups.
//
// Example:
//
//	contacts, err := client.Contact().Contacts().Get()
func (c *Contact) Contacts() *Contacts {
	return &Contacts{handler: c.handler}
}

// Groups returns operations for managing contact groups.
//
// Contact groups organize multiple contacts and can be assigned to
// hosts, monitors, and profiles for alert distribution.
//
// Example:
//
//	groups, err := client.Contact().Groups().Get()
func (c *Contact) Groups() *Groups {
	return &Groups{handler: c.handler}
}

// Contacts handles contact management operations.
//
// Provides CRUD operations for individual contacts including create, retrieve,
// update, and delete. Also provides GetAll() for automatic pagination.
type Contacts struct {
	handler *RequestHandler
}

// Get retrieves contacts with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all contacts (first page)
//   - Get(map[string]interface{}{"status": "active"}): Retrieves filtered contacts
//   - Get("COxxxxxxxxxx"): Retrieves a single contact by ID
//   - Get("COxxxxxxxxxx", "COyyyyyyyyyy"): Retrieves multiple contacts by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// Example:
//
//	// Get all contacts
//	contacts, err := client.Contact().Contacts().Get()
//
//	// Get single contact
//	contact, err := client.Contact().Contacts().Get("CO1234567890abcdef1234567890abcdef")
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

// GetAll retrieves all contacts with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of contacts without manual pagination handling.
//
// Example:
//
//	contacts, err := client.Contact().Contacts().GetAll(nil)
//	fmt.Printf("Total contacts: %d\n", len(contacts))
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

// Create creates a new contact.
//
// Required parameters:
//   - type: Contact type ("email" or "sms")
//   - value: Email address or phone number (E.164 format for SMS)
//
// Optional parameters:
//   - name: Descriptive name for the contact
//
// Example:
//
//	contact, err := client.Contact().Contacts().Create(map[string]interface{}{
//	    "type": "email",
//	    "value": "alerts@example.com",
//	    "name": "Operations Team",
//	})
func (c *Contacts) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Post("contact/contacts", params)
}

// Update updates a contact.
//
// The id parameter should be the contact SID (e.g., "CO1234567890abcdef1234567890abcdef").
//
// Updatable fields:
//   - name: Update the contact name
//   - value: Update the email or phone number
//
// Example:
//
//	_, err := client.Contact().Contacts().Update("CO1234567890abcdef1234567890abcdef", map[string]interface{}{
//	    "name": "Updated Contact Name",
//	})
func (c *Contacts) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return c.handler.Put(fmt.Sprintf("contact/contacts/%v", id), params)
}

// Delete deletes a contact.
//
// The id parameter should be the contact SID (e.g., "CO1234567890abcdef1234567890abcdef").
// Note: Cannot delete a contact that is currently assigned to a contact group.
//
// Example:
//
//	_, err := client.Contact().Contacts().Delete("CO1234567890abcdef1234567890abcdef")
func (c *Contacts) Delete(id interface{}) (map[string]interface{}, error) {
	return c.handler.Delete(fmt.Sprintf("contact/contacts/%v", id))
}

// Groups handles contact group management operations.
//
// Contact groups organize multiple contacts and can be assigned to hosts,
// monitors, and profiles for alert distribution. Provides CRUD operations.
type Groups struct {
	handler *RequestHandler
}

// Get retrieves contact groups with flexible input options.
//
// Usage patterns:
//   - Get(): Retrieves all groups (first page)
//   - Get(map[string]interface{}{"page": 2}): Retrieves with pagination
//   - Get("CGxxxxxxxxxx"): Retrieves a single group by ID
//   - Get("CGxxxxxxxxxx", "CGyyyyyyyyyy"): Retrieves multiple groups by IDs
//
// For automatic pagination across all pages, use GetAll() instead.
//
// Example:
//
//	// Get all groups
//	groups, err := client.Contact().Groups().Get()
//
//	// Get single group
//	group, err := client.Contact().Groups().Get("CG4f3e2d1c0b9a8776655443322110fedc")
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

// GetAll retrieves all contact groups with automatic pagination.
//
// This method automatically iterates through all pages and returns a complete
// list of contact groups without manual pagination handling.
//
// Example:
//
//	groups, err := client.Contact().Groups().GetAll(nil)
//	fmt.Printf("Total groups: %d\n", len(groups))
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

// Create creates a new contact group.
//
// Required parameters:
//   - name: Descriptive name for the group
//   - contacts: Array of contact SIDs to include in the group
//
// Example:
//
//	group, err := client.Contact().Groups().Create(map[string]interface{}{
//	    "name": "Operations Team",
//	    "contacts": []string{"CO1234567890", "CO0987654321"},
//	})
func (g *Groups) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return g.handler.Post("contact/groups", params)
}

// Update updates a contact group.
//
// The id parameter should be the group SID (e.g., "CG4f3e2d1c0b9a8776655443322110fedc").
//
// Updatable fields:
//   - name: Update the group name
//   - contacts: Update the list of contact SIDs in the group
//
// Example:
//
//	_, err := client.Contact().Groups().Update("CG4f3e2d1c0b9a8776655443322110fedc", map[string]interface{}{
//	    "name": "Updated Group Name",
//	})
func (g *Groups) Update(id interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	return g.handler.Put(fmt.Sprintf("contact/groups/%v", id), params)
}

// Delete deletes a contact group.
//
// The id parameter should be the group SID (e.g., "CG4f3e2d1c0b9a8776655443322110fedc").
// Note: Cannot delete a group that is currently assigned to hosts, monitors, or profiles.
//
// Example:
//
//	_, err := client.Contact().Groups().Delete("CG4f3e2d1c0b9a8776655443322110fedc")
func (g *Groups) Delete(id interface{}) (map[string]interface{}, error) {
	return g.handler.Delete(fmt.Sprintf("contact/groups/%v", id))
}
