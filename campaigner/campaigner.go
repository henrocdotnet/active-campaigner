package campaigner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// Campaigner is a library for interacting with ActiveCampaign.
type Campaigner struct {
	APIToken string
	BaseURL  string
}

// CheckConfig checks that API Token and BaseURL have been defined.
func (c *Campaigner) CheckConfig() error {
	if len(c.APIToken) == 0 {
		return CustomError{}.SetMessage("campaigner API token not set")
	} else if len(c.BaseURL) == 0 {
		return CustomError{}.SetMessage("campaigner base URL not set")
	}

	return nil
}

// GenerateURL returns a full API URL using the configured BaseURL and a suffix (API endpoint).
func (c *Campaigner) GenerateURL(url string) string {
	if strings.HasPrefix(url, "/") {
		url = strings.Replace(url, "/", "", 1)
	}

	url = fmt.Sprintf("%s/%s", c.BaseURL, url)

	return url
}

func (c *Campaigner) delete(url string) (gorequest.Response, []byte, error) {
	// Locals.
	var (
		r    gorequest.Response
		b    []byte
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	r, b, errs = gorequest.New().
		Delete(c.GenerateURL(url)).
		Set("Api-Token", c.APIToken).
		EndBytes()

	if errs != nil {
		return r, b, CustomError{Message: "could not perform HTTP DELETE request", HTTPErrors: errs}
	}

	return r, b, nil
}

func (c *Campaigner) get(url string) (gorequest.Response, []byte, error) {
	// Locals.
	var (
		r    gorequest.Response
		b    []byte
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	url = c.GenerateURL(url)

	r, b, errs = gorequest.New().
		Get(url).
		Set("Api-Token", c.APIToken).
		EndBytes()

	if errs != nil {
		return r, b, CustomError{Message: "could not perform HTTP GET request", HTTPErrors: errs}
	}

	// TODO(questions): Not sure if output should be indented by default.  Makes dev easier and indentation should never break things but still smells bad.
	var pretty bytes.Buffer
	err := json.Indent(&pretty, b, "", "\t")
	if err != nil {
		panic(err)
	}

	return r, b, nil
}

// Send a POST request to the Active Campaign API.
func (c *Campaigner) post(url string, i interface{}) (r gorequest.Response, b []byte, err error) {
	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	// Generate URL and JSON.
	url = c.GenerateURL(url)
	j, err := json.Marshal(i)
	if err != nil {
		return r, b, fmt.Errorf("Could not marshall json for interface: %s\n", err)
	}

	// Send POST request.
	r, b, errs := gorequest.New().
		Post(url).
		Send(string(j)).
		Set("Api-Token", c.APIToken).
		EndBytes()

	// Error check.
	if errs != nil {
		return r, b, CustomError{Message: "could not perform HTTP POST request", HTTPErrors: errs}
	}

	return r, b, nil
}

// Send a PUT request to the Active Campaign API.
func (c *Campaigner) put(url string, i interface{}) (r gorequest.Response, b []byte, err error) {
	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	url = c.GenerateURL(url)

	j, err := json.Marshal(i)
	if err != nil {
		return r, b, fmt.Errorf("Could not marshall json for interface: %s\n", err)
	}

	r, b, errs := gorequest.New().
		Put(url).
		Send(string(j)).
		Set("Api-Token", c.APIToken).
		EndBytes()

	// Error check.
	if errs != nil {
		return r, b, CustomError{Message: "could not perform HTTP POST request", HTTPErrors: errs}
	}

	return r, b, nil
}
