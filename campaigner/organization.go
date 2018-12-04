package campaigner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

// Represents an organization as it exists in the API.
type Organization struct {
	Name         string        `json:"name"`
	Links        []interface{} `json:"links"`
	ID           int64         `json:"id,string"`
	ContactCount string        `json:"contactCount"`
	DealCount    string        `json:"dealCount"`
}

func (c *Campaigner) OrganizationCreate(org Organization) (ResponseOrganizationCreate, error) {
	var (
		url  = "/api/3/organizations"
		data = map[string]interface{}{
			"organization": org,
		}
		result ResponseOrganizationCreate
	)

	r, b, err := c.post(url, data)
	if err != nil {
		err.(CustomError).WriteToLog()
		return result, fmt.Errorf("could not creation organization, HTTP failure: %s", err)
	}

	if r.StatusCode == http.StatusCreated {
		err = json.Unmarshal(b, &result)
		if err != nil {
			return result, fmt.Errorf("could not create organization, JSON failure: %s", err)
		}

		return result, nil
	}

	if r.StatusCode == http.StatusUnprocessableEntity {
		var apiError ActiveCampaignError
		err = json.Unmarshal(b, &apiError)
		if err != nil {
			return result, fmt.Errorf("could not unmarshal API error json: %s", err)
		}

		return result, apiError
	}

	log.Printf("response: %#v\n", r)
	log.Printf("body: %s\n", string(b))

	return result, nil
}

// Delete an organization by it's ID.
//
// TODO(error-checking): Are there other HTTP status codes to check for?
func (c *Campaigner) OrganizationDelete(id int64) error {
	// Setup.
	var (
		url = fmt.Sprintf("/api/3/organizations/%d", id)
	)

	// Send DELETE request.
	r, b, err := c.Delete(url)
	if err != nil {
		return fmt.Errorf("organization delete failed, HTTP failure: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusNotFound:
		e := new(CustomErrorNotFound)
		e.Message = fmt.Sprintf("organization delete failed, ID `%d` not found", id)
		return e
	case http.StatusOK:
		return nil
	default:
		log.Printf("response? %#v\n", r)
		return fmt.Errorf("organization delete failed, unspecified error: %s", b)
	}
}

// Find an organization by it's name.
//
// TODO(API): Figure out if more than one name can be searched (wildcard?  partial name?).
func (c *Campaigner) OrganizationFind(n string) (ResponseOrganizationList, error) {
	// Setup.
	var (
		qs = fmt.Sprintf("%s=%s", url2.QueryEscape("filters[name]"), url2.QueryEscape(n))
		url = fmt.Sprintf("/api/3/organizations/?%s", qs)
		response ResponseOrganizationList
	)

	// Error check.
	if len(strings.TrimSpace(n)) == 0 {
		return response, fmt.Errorf("organization find failed, name is empty")
	}

	log.Printf("url? %s\n", url)

	// Send GET request.
	r, b, err := c.get(url)
	if err != nil {
		return response, fmt.Errorf("organization find failed. HTTP failure: %s", err)
	}

	log.Printf("OrganizationFind: status code: %d\n", r.StatusCode)

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("organization list failed, JSON failure: %s", err)
	}

	return response, nil
}

// List all organizations.
func (c *Campaigner) OrganizationList() (ResponseOrganizationList, error) {
	// Setup.
	var (
		url      = "/api/3/organizations"
		response ResponseOrganizationList
	)

	// GET request.
	r, b, err := c.get(url)
	if err != nil {
		err.(CustomError).WriteToLog()
		return response, fmt.Errorf("organization list failed, HTTP failure: %s", err)
	}

	// Success.
	// TODO(doc-mismatch): 200 != 201
	if r.StatusCode == http.StatusOK {
		err = json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("organization list failed, JSON failure: %s", err)
		}

		return response, nil
	}

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("organization list failed, unspecified error: %s", b)
}

