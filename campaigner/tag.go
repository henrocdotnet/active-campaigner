package campaigner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

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
	dump(data)
	logFormattedJSON("tag create data?", data)

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

	// return response, nil

	// POST request.
	r, b, err := c.post(target, data)
	if err != nil {
		err.(CustomError).WriteToLog()
		return response, fmt.Errorf("tag creation failed, HTTP error: %s", err)
	}

	// Success.
	if r.StatusCode == http.StatusCreated {
		err = json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("tag creation failed, JSON error: %s", err)
		}

		return response, nil
	}

	dump(r)

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("tag creation failed, unspecified error: %s", err)
}

// TODO(API): This method is not documented.
func (c *Campaigner) TagDelete(id int64) error {
	// Setup.
	var (
		target = fmt.Sprintf("/api/3/tags/%d", id)
	)

	// Send DELETE request.
	r, b, err := c.Delete(target)
	if err != nil {
		return fmt.Errorf("tag deletion failed, HTTP error: %s", err)
	}

	dump(r)
	dump(b)

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("tag deletion failed, unspecified error: %s", b)
	}
}

// TODO(API): Query parameters aren't documented.
func (c *Campaigner) TagFind(n string) (ResponseTagList, error) {
	// Setup.
	var (
		qs       = fmt.Sprintf("%s=%s", url.QueryEscape("filters[tag]"), url.QueryEscape(n))
		target   = fmt.Sprintf("/api/3/tags/?%s", qs)
		response ResponseTagList
	)
	log.Println(target)

	// Error check.
	if len(strings.TrimSpace(n)) == 0 {
		return response, fmt.Errorf("tag find failed, name is empty")
	}

	log.Printf("url? %s\n", target)

	// Send GET request.
	r, b, err := c.get(target)
	if err != nil {
		return response, fmt.Errorf("tag find failed, HTTP error: %s", err)
	}

	log.Printf("TagFind: status code: %d\n", r.StatusCode)

	err = json.Unmarshal(b, &response)
	if err != nil {
		return response, fmt.Errorf("tag list failed, JSON error: %s", err)
	}

	return response, nil
}

func (c *Campaigner) TagList() (ResponseTagList, error) {
	// Setup.
	var (
		target   = "/api/3/tags?limit=100"
		response ResponseTagList
	)

	// GET request.
	r, b, err := c.get(target)
	if err != nil {
		err.(CustomError).WriteToLog()
		return response, fmt.Errorf("tag list failed, HTTP error: %s", err)
	}

	// Success.
	if r.StatusCode == http.StatusOK {
		err = json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("tag list failed, JSON error: %s", err)
		}

		return response, nil
	}

	// Failure (API docs are not clear about errors here).
	return response, fmt.Errorf("tag list failed, unspecified error: %s", err)
}

// TODO(API): This is not documented.
func (c *Campaigner) TagRead(id int64) (response ResponseTagRead, err error) {
	// Setup.
	var target = fmt.Sprintf("/api/3/tags/%d", id)

	// Get request.
	r, b, err := c.get(target)
	if err != nil {
		return response, fmt.Errorf("tag read failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(b, &response)
		if err != nil {
			return response, fmt.Errorf("tag read failed, JSON error: %s", err)
		}

		return response, nil

	case http.StatusNotFound:
		e := CustomErrorNotFound{ CustomError{Message: fmt.Sprintf("tag with id %d not found", id) }}
		// return response, CustomErrorNotFound{}.SetMessage("tag with id %d not found", id)
		return response, e

	default:
		return response, fmt.Errorf("tag read failed, unspecified error (%d): %s", r.StatusCode, string(b))
	}
}

type Tag struct {
	ID              int64    `json:"id,string"`
	Type            string   `json:"tagType"`
	Name            string   `json:"tag"`
	Description     string   `json:"description"`
	CreationDate    string   `json:"cdate"`
	SubscriberCount int64    `json:"subscriber_count,string"`
	Links           TagLinks `json:"links"`
}

type TagLinks struct {
	ContactGoalTags string `json:"contactGoalTags"`
}

type ResponseTagCreate struct {
	Tag Tag `json:"tag"`
}

type ResponseTagList struct {
	Tags []Tag `json:"tags"`
	Meta struct {
		Total string `json:"total"`
	} `json:"meta"`
}

type ResponseTagRead struct {
	Tag Tag `json:"tag"`
}
