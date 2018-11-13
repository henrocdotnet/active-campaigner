package campaigner

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

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

	log.Printf("ContactList: body:\n%s", string(pb.Bytes()))

	return nil
}


func (c *Campaigner) ContactCreate(contact Contact) (ResponseContactCreate, error) {
	// Setup.
	var (
		result ResponseContactCreate
		url = "/api/3/contacts"
		data = map[string]interface{}{
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

		writeIndentedJson(path.Join(os.Getenv("TEMP"), "contact_create_response.json"), []byte(body))
		logFormattedJson("ContactCreate Result:", result)

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
	} else {
		return result, fmt.Errorf("contact creation failed, unspecified error: %s", body)
	}
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
func (c *Campaigner) ContactRead(id int) (ResponseContactRead, error) {
	// Locals.
	var (
		response ResponseContactRead
		url = fmt.Sprintf("/api/3/contacts/%d", id)
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

	writeIndentedJson(path.Join(os.Getenv("TEMP"), "contact_read_response.json"), b)

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("could not unmarshal JSON for contact %d: %s", id, err)
	}

	return response, nil
}
