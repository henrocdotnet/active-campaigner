package campaigner

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

// TODO(API): Need to add return type here (and return it).
// TODO(error-checking): Need to add call to CheckConfig before calling API.
// TODO(error-checking): Needs to return error of some kind.
func (c *Campaigner) ContactList() error {
	// Perform the query.
	url := "/api/3/contacts"
	_, body, err := c.get(url)
	if err != nil {
		err.(CustomError).WriteToLog()
		return err
	}

	// Process the result.
	var j map[string]interface{}
	err = json.Unmarshal(body, &j)
	if err != nil {
		log.Fatalf("Could not unmarshal body: %s", err)
	}

	/*

	// Debug: Write the result to a local file.
	// TODO: This should use writeIndentedJson function.
	var pb bytes.Buffer
	f, err := os.Create("sandbox.local/contact-list.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json.Indent(&pb, body, "", "\t")
	_, err = f.Write(pb.Bytes())
	if err != nil {
		panic(err)
	}
	*/

	// log.Printf("ContactList: body:\n%s", string(pb.Bytes()))
	logFormattedJSON("contact list", j)
	//log.Printf("ContactList: body:\n%s", string(j))

	return nil
}

func (c *Campaigner) ContactCreate(contact Contact) (ResponseContactCreate, error) {
	// Setup.
	var (
		result ResponseContactCreate
		url    = "/api/3/contacts"
		data   = map[string]interface{}{
			"contact": contact,
		}
	)

	// POST request.
	r, body, err := c.post(url, data)
	if err != nil {
		return result, fmt.Errorf("contact creation failed, HTTP error: %s", err)
	}

	// Success.
	if r.StatusCode == http.StatusCreated {
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			return result, fmt.Errorf("contact creation failed, json error: %s", err)
		}

		writeIndentedJSON(path.Join(os.Getenv("TEMP"), "contact_create_response.json"), []byte(body))
		logFormattedJSON("ContactCreate Result:", result)

		return result, nil

	}

	// Failure.
	if r.StatusCode == http.StatusUnprocessableEntity {
		var apiError ActiveCampaignError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return result, fmt.Errorf("could not unmarshal API error json: %s", err)
		}

		return result, apiError
	}
	return result, fmt.Errorf("contact creation failed, unspecified error: %s", body)
}

func (c *Campaigner) ContactDelete(id int) error {
	r, b, err := c.Delete(fmt.Sprintf("/api/3/contacts/%d", id))

	if err != nil {
		log.Printf("ERROR CONTACT DELETE: %s\n", err)
		err.(CustomError).WriteToLog()

	}

	if r.StatusCode != 200 {
		msg := fmt.Sprintf("could not delete contact: %s", b)
		return errors.New(msg)
	}

	return nil
}

// TODO: Some other junk is being returned besides the contact itself.  Not sure if this should be parsed and wrapped back into the main contact struct.
func (c *Campaigner) ContactRead(id int64) (ResponseContactRead, error) {
	// Locals.
	var (
		response ResponseContactRead
		url      = fmt.Sprintf("/api/3/contacts/%d", id)
	)

	// Perform query.
	r, b, err := c.get(url)
	if err != nil {
		return response, err
	}

	// Error check.
	if r.StatusCode == http.StatusNotFound { // Contact not found.
		return response, fmt.Errorf("read failed: ID %d not found", id)
	} else if r.StatusCode != http.StatusOK { // Other error.
		return response, fmt.Errorf("read failed: ID %d general error", id)
	}

	writeIndentedJSON(path.Join(os.Getenv("TEMP"), "contact_read_response.json"), b)

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("could not unmarshal JSON for contact %d: %s", id, err)
	}

	return response, nil
}

// TODO(API): The API return JSON also includes a contacts[] entry with one contact in it.  Tested this on a tag I know is attached to more than one contact.
// TODO(API): The API returns different JSON for a request with a bogus ID in it.  contact and tag ID are returned as strings instead of ints.
func (c *Campaigner) ContactTag(task TaskContactTag) (response ResponseContactTag, err error) {
	// Setup.
	var (
		url  = "/api/3/contactTags"
		data = map[string]interface{}{
			"contactTag": task,
		}
	)

	// API doesn't appear to do much validation on tag or contact IDs passed.
	if task.TagID < 1 {
		return response, fmt.Errorf("contact tagging failed, task has invalid tag ID")
	} else if task.ContactID < 1 {
		return response, fmt.Errorf("contact tagging failed, task has invalid contact ID")
	}

	// Check that contact exists.
	_, err = c.ContactRead(task.ContactID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find contact: %s", err)
	}

	// Check that tag exists.
	_, err = c.TagRead(task.TagID)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, could not find tag: %s", err)
	}

	log.Println("DUMPING TASK?")
	dump(task)

	// Send POST request.
	r, b, err := c.post(url, data)
	if err != nil {
		return response, fmt.Errorf("contact tagging failed, HTTP error: %s", err)
	}

	log.Printf("WHERE IS MY BODY? (%d)\n", r.StatusCode)
	log.Println(string(b))

	// Response check.
	// TODO(API): I am getting both 200 and 201 on this step.  It looks like 201 is returned for a brand new association and 200 if it already existed.
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

type ContactTag struct {
	DateCreated string          `json:"cdate"`
	ContactID   int64json       `json:"contact"`
	ID          int64json       `json:"id"`
	TagID       int64json       `json:"tag"`
	Links       ContactTagLinks `json:"links"`
}

type ContactTagLinks struct {
	Contact string `json:"contact"`
	Tag     string `json:"tag"`
}

type ResponseContactTag struct {
	ContactTag ContactTag `json:"contactTag"`
}
