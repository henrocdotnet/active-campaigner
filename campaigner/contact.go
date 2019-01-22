package campaigner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ContactList lists contacts.
//
// TODO(api): Is there support for searching/filtering?
func (c *Campaigner) ContactList() (response ResponseContactList, err error) {
	// TODO(API): Need to add return type here (and return it).
	// Setup.
	url := "/api/3/contacts"

	// Send GET request.
	r, body, err := c.get(url)
	if err != nil {
		return response, fmt.Errorf("contact list failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("contact list failed, JSON error: %s", err)
		}

		return response, nil

	default:
		return response, fmt.Errorf("contact list failed, unspecified error: %s", body)
	}
}

// ContactCreate creates a contact.
func (c *Campaigner) ContactCreate(contact Contact) (result ResponseContactCreate, err error) {
	// Setup.
	var (
		url    = "/api/3/contacts"
		data   = map[string]interface{}{
			"contact": contact,
		}
	)

	// Send POST request.
	r, body, err := c.post(url, data)
	if err != nil {
		return result, fmt.Errorf("contact creation failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusCreated: // Success.
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			return result, fmt.Errorf("contact creation failed, json error: %s", err)
		}

		return result, nil

	case http.StatusUnprocessableEntity:
		var apiError ActiveCampaignError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return result, fmt.Errorf("could not unmarshal API error json: %s", err)
		}

		return result, apiError

	default:
		return result, fmt.Errorf("contact creation failed, unspecified error: %s", body)
	}
}

// ContactDelete deletes a contact.
func (c *Campaigner) ContactDelete(id int64) error {
	// TODO(error-checking): Are there specific HTTP codes that can be checked for?
	// Send DELETE request.
	r, b, err := c.delete(fmt.Sprintf("/api/3/contacts/%d", id))
	if err != nil {
		return fmt.Errorf("contact deletion failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK: // Success.
		return nil
	default:
		return fmt.Errorf("could not delete contact: %s", b)
	}
}

// ContactRead reads a contact.
func (c *Campaigner) ContactRead(id int64) (response ResponseContactRead, err error) {
	// TODO(response-parsing): Quite a bit of extra data is being returned besides the contact itself.  Not sure if this should be parsed and wrapped back into the main contact struct.
	// Setup.
	var url = fmt.Sprintf("/api/3/contacts/%d", id)

	// Send GET request.
	r, body, err := c.get(url)
	if err != nil {
		return response, err
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK: // Success.
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("contact read failed, JSON error: %s", err)
		}

		return response, nil
	case http.StatusNotFound:
		return response, fmt.Errorf("contact read failed: ID %d not found", id)
	default:
		return response, fmt.Errorf("contact read failed: unspecified error: %s", string(body))
	}
}

// ContactTagCreate links a tag to a contact.
//
// TODO(API): The API return JSON also includes a contacts[] entry with one contact in it.  Tested this on a tag I know is attached to more than one contact.
//
// TODO(API): The API returns different JSON for a request with a bogus ID in it.  The contact and tag ID are returned as strings instead of ints.
func (c *Campaigner) ContactTagCreate(request RequestContactTagCreate) (response ResponseContactTagCreate, err error) {
	// TODO(error-checking): Is it possible to check for a not found error specifically?
	// TODO(error-checking): Nonexistent contact or tag should return a CustomErrorNotFound error.
	// Setup.
	var (
		url  = "/api/3/contactTags"
		data = map[string]interface{}{
			"contactTag": request,
		}
	)

	// API doesn't appear to do much validation on tag or contact IDs passed.
	if request.TagID < 1 {
		return response, fmt.Errorf("contact tagging failed, task has invalid tag ID")
	} else if request.ContactID < 1 {
		return response, fmt.Errorf("contact tagging failed, task has invalid contact ID")
	}

	// Check that contact exists.
	_, err = c.ContactRead(request.ContactID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find contact: %s", err)
	}

	// Check that tag exists.
	_, err = c.TagRead(request.TagID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find tag: %s", err)
	}

	// Send POST request.
	r, b, err := c.post(url, data)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, HTTP error: %s", err)
	}

	// Response check.
	// TODO(API): I am getting both 200 and 201 on this step.  It looks like 201 is returned for a brand new association and 200 if it already exists.
	switch r.StatusCode {
	case http.StatusOK: // Association already exists.  Response JSON does not include the contacts:[] bit.
		err := json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("contact tagging failed, JSON error: %s", err)
		}
	case http.StatusCreated: // Brand new association.  Response JSON includes the contact, bleh.
		err := json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("contact tagging failed, JSON error: %s", err)
		}
	default:
		return response, fmt.Errorf("contact tagging failed, unspecified error (%d): %s", r.StatusCode, string(b))
	}

	return response, nil

}

// ContactTagDelete removes a tag from a contact.  This removes the "link" and not the tag itself.
func (c *Campaigner) ContactTagDelete(id int64) error {
	// Setup.
	var (
		url  = fmt.Sprintf("/api/3/contactTags/%d", id)
	)

	// Send DELETE request.
	r, b, err := c.delete(url)
	if err != nil {
		return fmt.Errorf("contact tag deletion failed, HTTP failure: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK: // Success.
		return nil
	case http.StatusNotFound:
		e := new(CustomErrorNotFound)
		e.Message = fmt.Sprintf("contact tag deletion failed, ID `%d` not found", id)
		return e
	default:
		return fmt.Errorf("contact tag deletion failed, unspecified error: %s", b)
	}
}

// ContactTagReadByContactID reads assigned tags for a contact by it's ID.
func (c *Campaigner) ContactTagReadByContactID(id int64) (response ResponseContactTagRead, err error) {
	// Setup.
	var url  = fmt.Sprintf("/api/3/contacts/%d/contactTags", id)

	// Send GET request.
	r, body, err := c.get(url)
	if err != nil {
		return response, fmt.Errorf("contact tags read failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK: // Success.
		if err = json.Unmarshal(body, &response); err != nil {
			return response, fmt.Errorf("contact tags read failed, JSON error: %s", err)
		}

		return response, nil

	case http.StatusNotFound:
		e := new(CustomErrorNotFound)
		e.Message = fmt.Sprintf("contact tags read failed, ID `%d` not found", id)
		return response, e

	default:
		log.Printf("response? %#v\n", r)
		return response, fmt.Errorf("contact tags read failed, unspecified error: %s", string(body))
	}
}

// ContactTag holds a JSON compatible contact tag.
type ContactTag struct {
	DateCreated string          `json:"cdate"`
	ContactID   int64json       `json:"contact"`
	ID          int64json       `json:"id"`
	TagID       int64json       `json:"tag"`
	Links       ContactTagLinks `json:"links"`
}

// ContactTagLinks holds a JSON compatible list of contact tag links (nested structure, see ContactTags).
type ContactTagLinks struct {
	// TODO(rename): Possibly rename.
	Contact string `json:"contact"`
	Tag     string `json:"tag"`
}

// ResponseContactTagCreate holds a JSON compatible response for creating contacts.
type ResponseContactTagCreate struct {
	ContactTag ContactTag `json:"contactTag"`
}

// ResponseContactTagRead holds a JSON compatible response for reading contacts.
type ResponseContactTagRead struct {
	ContactTags []ContactTag `json:"contactTags"`
}
