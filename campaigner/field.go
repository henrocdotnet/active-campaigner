package campaigner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Field struct {
	Title        string        `json:"title"`
	Description  string        `json:"descript"`
	Type         string        `json:"type"`
	IsRequired   string        `json:"isrequired"`
	Perstag      string        `json:"perstag"`
	DefaultValue string        `json:"defval"`
	ShowInList   string        `json:"show_in_list"`
	Rows         string        `json:"rows"`
	Columns      string        `json:"cols"`
	IsVisible    string        `json:"visible"`
	Service      string        `json:"service"`
	OrderNumber  string        `json:"ordernum"`
	DateCreated  string        `json:"cdate"`
	DateUpdated  string        `json:"udate"`
	Options      []interface{} `json:"options"`
	Relations    []string      `json:"relations"`
	Links        struct {
		Options   string `json:"options"`
		Relations string `json:"relations"`
	} `json:"links"`
	ID int64 `json:"id,string"`
}

type ResponseFieldList struct {
	FieldOptions       []interface{}                `json:"fieldOptions"`
	FieldRelationships []ResponseFieldRelationships `json:"fieldRels"`
	Fields             []Field                      `json:"fields"`
	Meta               struct {
		Total string `json:"total"`
	} `json:"meta"`
}

type ResponseFieldRead struct {
	FieldOptions       []interface{}                `json:"fieldOptions"`
	FieldRelationships []ResponseFieldRelationships `json:"fieldRels"`
	Field              Field                        `json:"field"`
}

type ResponseFieldRelationships struct {
	Field        string        `json:"field"`
	RelationID   string        `json:"relid"`
	DisplayOrder string        `json:"dorder"`
	DateCreated  string        `json:"cdate"`
	Links        []interface{} `json:"links"`
	ID           string        `json:"id"`
}

// FieldList lists custom fields.
func (c *Campaigner) FieldList() (response ResponseFieldList, err error) {
	// Setup.
	u := "/api/3/fields"

	// Send GET request.
	r, body, err := c.get(u)
	if err != nil {
		return response, fmt.Errorf("field list failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		log.Println(string(body))
		err := json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("field list failed, JSON error: %s", err)
		}

		//logFormattedJSON("field list", response)

		return response, nil
	}

	return response, fmt.Errorf("field list failed, unspecified error (%d): %s", r.StatusCode, string(body))
}

// FieldRead reads a custom field.
func (c *Campaigner) FieldRead(id int64) (response ResponseFieldRead, err error) {
	// Setup.
	u := fmt.Sprintf("/api/3/fields/%d", id)
	log.Println(u)

	// Send GET request.
	r, body, err := c.get(u)
	if err != nil {
		return response, fmt.Errorf("field read failed, HTTP error: %s", err)
	}

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(body, &response)
		if err != nil {
			return response, fmt.Errorf("field read failed, JSON error: %s", err)
		}
		//log.Println(string(body))
		//logFormattedJSON("field read", response)
		//dump(response)

		return response, nil
	}

	return response, fmt.Errorf("field read failed, unspecified error (%d): %s", r.StatusCode, string(body))
}
