package campaigner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TagCreate creates a tag.
func (c *Campaigner) TagCreate(tag Tag) (ResponseTagCreate, error) {
	// Setup.
	var (
		target   = "/api/3/tags"
		response ResponseTagCreate
	)

	var data = map[string]interface{}{
		"tag": map[string]string{
			"tag":         tag.Name,
			"description": tag.Description,
			"tagType":     tag.Type,
		},
	}

	// Tag check.
	if len(tag.Name) == 0 {
		return response, fmt.Errorf("tag creation failed, name is empty")
	}
	if len(tag.Description) == 0 {
		return response, fmt.Errorf("tag creation failed, description is empty")
	}
	if len(tag.Type) == 0 {
		return response, fmt.Errorf("tag creation failed, type is empty")
	} else if tag.Type != "contact" && tag.Type != "template" {
		return response, fmt.Errorf("tag creation failed, type %s is invalid, possible values are template, contact", tag.Type)
	}

	// POST request.
	r, body, err := c.post(target, data)
	if err != nil {
		return response, fmt.Errorf("tag creation failed, HTTP error: %s", err)
	}

	// Success.
	if r.StatusCode == http.StatusCreated {
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("tag creation failed, JSON error: %s", err)
		}

		return response, nil
	}

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("tag creation failed, unspecified error (%d): %s", r.StatusCode, string(body))
}

// TagDelete deletes a tag.
//
// TODO(api): This method is missing in the ActiveCampaign documentation.
func (c *Campaigner) TagDelete(id int64) error {
	// Setup.
	var (
		target = fmt.Sprintf("/api/3/tags/%d", id)
	)

	// Send DELETE request.
	r, body, err := c.delete(target)
	if err != nil {
		return fmt.Errorf("tag deletion failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("tag deletion failed, unspecified error: %s", string(body))
	}
}

// TagFind searches for a tag by name.  The list of all available filters is not complete.
//
// TODO(error-checking: Add HTTP status code checking.
//
// TODO(api): Query parameters aren't officially documented.
func (c *Campaigner) TagFind(n string) (response ResponseTagList, err error) {
	// Setup.
	var (
		qs       = fmt.Sprintf("%s=%s", url.QueryEscape("filters[tag]"), url.QueryEscape(n))
		target   = fmt.Sprintf("/api/3/tags/?%s", qs)
	)

	// Error check.
	if len(strings.TrimSpace(n)) == 0 {
		return response, fmt.Errorf("tag find failed, name is empty")
	}

	// Send GET request.
	_, b, err := c.get(target)
	if err != nil {
		return response, fmt.Errorf("tag find failed, HTTP error: %s", err)
	}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("tag list failed, JSON error: %s", err)
	}

	return response, nil
}

// TagList lists all tags.
func (c *Campaigner) TagList() (ResponseTagList, error) {
	// Setup.
	var (
		target   = "/api/3/tags?limit=100"
		response ResponseTagList
	)

	// GET request.
	r, body, err := c.get(target)
	if err != nil {
		return response, fmt.Errorf("tag list failed, HTTP error: %s", err)
	}

	// Success.
	if r.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("tag list failed, JSON error: %s", err)
		}

		return response, nil
	}

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("tag list failed, unspecified error: %s", string(body))
}

// TagRead reads a tag by it's ID.
//
// TODO(api): This endpoint is not documented.
func (c *Campaigner) TagRead(id int64) (response ResponseTagRead, err error) {
	// Setup.
	var target = fmt.Sprintf("/api/3/tags/%d", id)

	// Get request.
	r, body, err := c.get(target)
	if err != nil {
		return response, fmt.Errorf("tag read failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("tag read failed, JSON error: %s", err)
		}

		return response, nil

	case http.StatusNotFound:
		e := CustomErrorNotFound{ CustomError{Message: fmt.Sprintf("tag with id %d not found", id) }}
		// return response, CustomErrorNotFound{}.SetMessage("tag with id %d not found", id)
		return response, e

	default:
		return response, fmt.Errorf("tag read failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// Tag holds a JSON compatible tag as it exists in the API.
type Tag struct {
	ID              int64    `json:"id,string"`
	Type            string   `json:"tagType"`
	Name            string   `json:"tag"`
	Description     string   `json:"description"`
	CreationDate    string   `json:"cdate"`
	SubscriberCount int64    `json:"subscriber_count,string"`
	Links           TagLinks `json:"links"`
}

// TagLinks holds a JSON compatible list of links (nested structure).
type TagLinks struct {
	ContactGoalTags string `json:"contactGoalTags"`
}

// ResponseTagCreate holds a JSON compatible response for creating tags.
type ResponseTagCreate struct {
	Tag Tag `json:"tag"`
}

// ResponseTagList holds a JSON compatible response for listing tags.
type ResponseTagList struct {
	Tags []Tag `json:"tags"`
	Meta struct {
		Total string `json:"total"`
	} `json:"meta"`
}

// ResponseTagRead holds a JSON compatible response for reading tags.
type ResponseTagRead struct {
	Tag Tag `json:"tag"`
}
