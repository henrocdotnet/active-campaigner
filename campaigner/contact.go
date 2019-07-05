package campaigner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ContactList lists contacts.
func (c *Campaigner) ContactList(limit int, offset int) (response ResponseContactList, err error) {
	// Setup.
	qs := url.Values{}
	qs.Set("limit", strconv.Itoa(limit))
	qs.Set("offset", strconv.Itoa(offset))
	u := url.URL{ Path: "/api/3/contacts", RawQuery: qs.Encode() }

	fmt.Printf("contact list: url = %s\n", u.String())

	//qs := fmt.Sprintf("%s=%s", url2.Query)
	//url := "/api/3/contacts"

	// Send GET request.
	r, body, err := c.get(u.String())
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
	// TODO(api): The struct used in the request has to be rebuilt as AC does not like all of the extra fields in a "real contact".  This
	//            might be caused by sending the organization ID in the request which I don't think the unit tests currently cover.
	// Setup.
	var (
		uri  = "/api/3/contacts"
		data = map[string]interface{}{
			"contact": struct {
				EmailAddress   string    `json:"email"`
				PhoneNumber    string    `json:"phone"`
				FirstName      string    `json:"firstName"`
				LastName       string    `json:"lastName"`
				OrganizationID Int64json `json:"orgid"`
			}{EmailAddress: contact.EmailAddress, PhoneNumber: contact.PhoneNumber, FirstName: contact.FirstName, LastName: contact.LastName, OrganizationID: contact.OrganizationID},
		}
	)

	// Send POST request.
	r, body, err := c.post(uri, data)
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


// ContactFind searches for a contact by email.
//
// Partial emails are not supported by the API.
func (c *Campaigner) ContactFind(email string) (response ResponseContactList, err error) {
	// Setup.
	var (
		qs       = fmt.Sprintf("%s=%s", url.QueryEscape("filters[email]"), url.QueryEscape(email))
		u        = fmt.Sprintf("/api/3/contacts/?%s", qs)
	)

	// Error check.
	if len(strings.TrimSpace(email)) == 0 {
		return response, fmt.Errorf("contact find failed, email is empty")
	}

	// Send GET request.
	r, body, err := c.get(u)
	if err != nil {
		return response, fmt.Errorf("contact find failed, HTTP failure: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("contact find failed, JSON failure: %s", err)
		}

		log.Printf("contact find response? %#v\n", response)
		return response, nil
	}

	return response, fmt.Errorf("contact find failed, unspecified error (%d); %s", r.StatusCode, string(body))

}

// ContactRead reads a contact.
func (c *Campaigner) ContactRead(id int64) (response ResponseContactRead, err error) {
	// TODO(response-parsing): Quite a bit of extra data is being returned besides the contact itself.  Not sure if this should be parsed and wrapped back into the main contact struct.
	// Setup.
	var uri = fmt.Sprintf("/api/3/contacts/%d", id)

	// Send GET request.
	r, body, err := c.get(uri)
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
		return response, fmt.Errorf("contact read failed: unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// ContactUpdate updates a contact.
func (c *Campaigner) ContactUpdate(id int64, request RequestContactUpdate) (response ResponseContactUpdate, err error) {
	// Send PUT request.
	u := fmt.Sprintf("/api/3/contact/sync")
	d := map[string]interface{}{"contact": request}
	r, body, err := c.post(u, d)
	if err != nil {
		return response, fmt.Errorf("contact update failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		log.Println("OK BODY FOOL")
		log.Println(string(body))
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("contact update failed, JSON error: %s", err)
		}
		return response, nil
	default:
		return response, fmt.Errorf("contact update failed, unspecified error (%d): %s", r.StatusCode, string(body))
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
		return fmt.Errorf("contact deletion failed, unspecified error: %s", string(b))
	}
}

// ContactFieldDeleteByFieldValueID deletes a contact field value by it's ID.  Note that this is different from the Field::ID.
func (c *Campaigner) ContactFieldDeleteByFieldValueID(id int64) (err error) {
	u := fmt.Sprintf("/api/3/fieldValues/%d", id)
	r, body, err := c.delete(u)
	if err != nil {
		return fmt.Errorf("contact field deletion failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		var r interface{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return fmt.Errorf("contact field deletion failed, JSON error: %s", err)
		}

		return nil
	default:
		return fmt.Errorf("contact field deletion failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// ContactFieldUpdate updates a custom field for a contact.
func (c *Campaigner) ContactFieldUpdate(contactID int64, fieldID int64, value string) (response ResponseContactFieldUpdate, err error) {
	// Check that both the contact and field exist.
	_, err = c.ContactRead(contactID)
	if err != nil {
		return response, fmt.Errorf("contact field update failed, could not find contact: %s", err)
	}
	_, err = c.FieldRead(fieldID)
	if err != nil {
		return response, fmt.Errorf("contact field update failed, could not find field: %s", err)
	}

	// Send POST request.
	req := RequestContactFieldUpdate{ContactID: contactID, FieldID: fieldID, Value: value}
	u := "/api/3/fieldValues"
	r, body, err := c.post(u, map[string]interface{}{"fieldValue": req})
	if err != nil {
		return response, fmt.Errorf("contact field update failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusCreated: // Create.
		fallthrough
	case http.StatusOK: // Update.
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("contact field update failed, JSON error: %s", err)
		}
		//dump(response)
		//logFormattedJSON("response", response)

		return response, nil
	}

	return response, fmt.Errorf("contact field update failed, unspecified error (%d): %s", r.StatusCode, string(body))
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
		uri  = "/api/3/contactTags"
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
	rC, err := c.ContactRead(request.ContactID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find contact: %s", err)
	}

	// Check that tag exists.
	rT, err := c.TagRead(request.TagID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find tag: %s", err)
	}

	// Send POST request.
	r, b, err := c.post(uri, data)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, HTTP error: %s", err)
	}

	// Response check.
	// TODO(API): I am getting both 200 and 201 on this step.  It looks like 201 is returned for a brand new association and 200 if it already exists.
	// Association already exists.  Response JSON does not include the contacts:[] bit.
	// Brand new association.  Response JSON includes the contact, bleh.
	switch r.StatusCode {
		case http.StatusOK, http.StatusCreated:
		if err := json.Unmarshal(b, &response); err != nil {
			return response, fmt.Errorf("contact tagging failed, JSON error: %s", err)
		}

		response.Custom.ContactEmail = rC.Contact.EmailAddress
		response.Custom.TagName = rT.Tag.Name

		return response, nil
	default:
		return response, fmt.Errorf("contact tagging failed, unspecified error (%d): %s", r.StatusCode, string(b))
	}
}

// ContactTagDelete removes a tag from a contact.  This removes the "link" and not the tag itself.
func (c *Campaigner) ContactTagDelete(id int64) error {
	// Setup.
	var (
		uri = fmt.Sprintf("/api/3/contactTags/%d", id)
	)

	// Send DELETE request.
	r, b, err := c.delete(uri)
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
		return fmt.Errorf("contact tag deletion failed, unspecified error: %s", string(b))
	}
}

// ContactTagReadByContactID reads assigned tags for a contact by it's ID.
func (c *Campaigner) ContactTagReadByContactID(id int64) (response ResponseContactTagRead, err error) {
	// Setup.
	var uri = fmt.Sprintf("/api/3/contacts/%d/contactTags", id)

	// Send GET request.
	r, body, err := c.get(uri)
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

// ResponseContactFieldUpdate holds a JSON compatible response for updating contact fields.
type ResponseContactFieldUpdate struct {
	Contacts   []Contact         `json:"contacts"`
	FieldValue ContactFieldValue `json:"fieldValue"`
}

// ContactTag holds a JSON compatible contact tag.
type ContactTag struct {
	DateCreated string          `json:"cdate"`
	ContactID   Int64json       `json:"contact"`
	ID          Int64json       `json:"id"`
	TagID       Int64json       `json:"tag"`
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
	Custom struct {
		ContactEmail string
		TagName      string
	}
	ContactTag ContactTag `json:"contactTag"`
}

// ResponseContactTagRead holds a JSON compatible response for reading contacts.
type ResponseContactTagRead struct {
	ContactTags []ContactTag `json:"contactTags"`
}
