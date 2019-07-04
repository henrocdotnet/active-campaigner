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

// Organization holds a JSON compatible organization as it exists in the API.
type Organization struct {
	// TODO(api): Contact and deal counts should probably not be strings.
	Name         string                 `json:"name"`
	Links        map[string]interface{} `json:"links"`
	ID           int64                  `json:"id,string"`
	ContactCount string                 `json:"contactCount"`
	DealCount    string                 `json:"dealCount"`
}

// ResponseOrganizationCreate holds a JSON compatible response for creating organizations.
type ResponseOrganizationCreate struct {
	Organization struct {
		Name  string                 `json:"name"`
		Links map[string]interface{} `json:"links"`
		ID    int64                  `json:"id,string"`
	} `json:"organization"`
}

// ResponseOrganizationUpdate holds a JSON compatible response for updating organizations.
type ResponseOrganizationUpdate struct {
	Organization struct {
		ID               int64                  `json:"id,string"`
		Name             string                 `json:"name"`
		Links            map[string]interface{} `json:"links"`
		CreatedTimestamp string                 `json:"created_timestamp"`
		UpdatedTimestamp string                 `json:"updated_timestamp"`
	} `json:"organization"`
}

// ResponseOrganizationRead holds a JSON compatible response for reading organizations.
type ResponseOrganizationRead struct {
	Organization Organization `json:"organization"`
}

// ResponseOrganizationList holds a JSON compatible response for listing organizations.
type ResponseOrganizationList struct {
	Organizations []Organization `json:"organizations"`
	Meta          struct {
		Total int `json:"total,string"`
	} `json:"meta"`
}

// OrganizationCreate creates an organization.
func (c *Campaigner) OrganizationCreate(org Organization) (result ResponseOrganizationCreate, err error) {
	var (
		uri  = "/api/3/organizations"
		data = map[string]interface{}{
			"organization": org,
		}
	)

	r, body, err := c.post(uri, data)
	if err != nil {
		return result, fmt.Errorf("organization creation failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusCreated:
		err = json.Unmarshal(body, &result)
		if err != nil {
			return result, fmt.Errorf("organization creation failed, JSON error: %s", err)
		}

		return result, nil
	case http.StatusUnprocessableEntity:
		var apiError ActiveCampaignError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return result, fmt.Errorf("organization creation failed, API error unmarshall error: %s", err)
		}

		return result, apiError
	default:
		return result, fmt.Errorf("organization creation failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// OrganizationDelete deletes an organization by it's ID.
//
// TODO(error-checking): Are there other HTTP status codes to check for?
func (c *Campaigner) OrganizationDelete(id int64) error {
	// Setup.
	var (
		uri = fmt.Sprintf("/api/3/organizations/%d", id)
	)

	// Send DELETE request.
	r, b, err := c.delete(uri)
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
		return fmt.Errorf("organization delete failed, unspecified error (%d): %s", r.StatusCode, string(b))
	}
}

// OrganizationFind finds an organization by it's name.  If there are no matches the response contains a list with zero length.
//
// TODO(API): Figure out if more than one name can be searched (wildcard?  partial name?).
func (c *Campaigner) OrganizationFind(n string) (ResponseOrganizationList, error) {
	// TODO(error-checking): Add status code checking.
	// Setup.
	var (
		qs       = fmt.Sprintf("%s=%s", url.QueryEscape("filters[name]"), url.QueryEscape(n))
		u        = fmt.Sprintf("/api/3/organizations/?%s", qs)
		response ResponseOrganizationList
	)

	// Error check.
	if len(strings.TrimSpace(n)) == 0 {
		return response, fmt.Errorf("organization find failed, name is empty")
	}

	// Send GET request.
	r, body, err := c.get(u)
	if err != nil {
		return response, fmt.Errorf("organization find failed. HTTP failure: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("organization list failed, JSON failure: %s", err)
		}

		return response, nil
	}

	return response, fmt.Errorf("organization find failed, unspecified error (%d); %s", r.StatusCode, string(body))
}

// OrganizationList lists all organizations.
func (c *Campaigner) OrganizationList(limit int, offset int) (response ResponseOrganizationList, err error) {
	// Setup.
	qs := url.Values{}
	qs.Set("limit", strconv.Itoa(limit))
	qs.Set("offset", strconv.Itoa(offset))
	u := url.URL{Path: "/api/3/organizations", RawQuery: qs.Encode()}

	// GET request.
	r, body, err := c.get(u.String())
	if err != nil {
		return response, fmt.Errorf("organization list failed, HTTP failure: %s", err)
	}

	// Success.
	// TODO(doc-mismatch): 200 != 201
	if r.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("organization list failed, JSON failure: %s", err)
		}

		return response, nil
	}

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("organization list failed, unspecified error (%d): %s", r.StatusCode, string(body))
}

// OrganizationRead reads an organization by it's ID.
//
// TODO(api): Possible bug in that contactCount and dealCount are not present in the JSON returned by read.
func (c *Campaigner) OrganizationRead(id int64) (response ResponseOrganizationRead, err error) {
	// TODO(error-checking): Should probably return a CustomErrorNotFound here.
	// Setup.
	u := fmt.Sprintf("/api/3/organizations/%d", id)

	// Send GET request.
	r, body, err := c.get(u)
	if err != nil {
		return response, err
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("organization read failed, JSON failure: %s", err)
		}
	case http.StatusNotFound:
		return response, fmt.Errorf("organization read failed, ID %d not found", id)
	default:
		return response, fmt.Errorf("organization read failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}

	return response, nil
}

// OrganizationUpdate updates an organization.
func (c *Campaigner) OrganizationUpdate(id int64, request RequestOrganizationUpdate) (response ResponseOrganizationUpdate, err error) {
	u := url.URL{Path: fmt.Sprintf("/api/3/organizations/%d", id)}
	d := map[string]interface{}{"organization": request}

	r, body, err := c.put(u.String(), d)
	if err != nil {
		return response, fmt.Errorf("organization updated failed, HTTP error: %s", err)
	}

	switch r.StatusCode {
	case http.StatusOK:
		log.Println(string(body))
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("organization update failed, JSON error: %s", err)
		}

	default:
		return response, fmt.Errorf("organization update failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}

	return response, nil
}
