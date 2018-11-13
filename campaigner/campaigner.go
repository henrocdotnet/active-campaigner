package campaigner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"strings"
)

type Campaigner struct {
	ApiToken string
	BaseURL  string
}


func (c *Campaigner) Bleh() {

}

func (c *Campaigner) CheckConfig() error {
	if len(c.ApiToken) == 0 {
		return CustomError{}.SetMessage("campaigner API token not set")
	} else if len(c.BaseURL) == 0 {
		return CustomError{}.SetMessage("campaigner base URL not set")
	}

	return nil
}

func (c *Campaigner) GenerateURL(url string) string {
	if strings.HasPrefix(url, "/") {
		url = strings.Replace(url, "/", "", 1)
	}

	url = fmt.Sprintf("%s/%s", c.BaseURL, url)

	return url
}

func (c *Campaigner) Delete(url string) (gorequest.Response, string, error){
	// Locals.
	var (
		r gorequest.Response
		b string
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	r, b, errs = gorequest.New().
		Delete(c.GenerateURL(url)).
		Set("Api-Token", c.ApiToken).
		End()

	if errs != nil {
		return r, b, CustomError{ Message: "could not perform HTTP DELETE request", HttpErrors: errs }
	}

	return r, b, nil
}

func (c *Campaigner) get(url string) (gorequest.Response, []byte, error) {
	// Locals.
	var (
		r gorequest.Response
		b []byte
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	url = c.GenerateURL(url)

	log.Printf("Campaigner.Get: Using url %s\n", url)

	r, b, errs = gorequest.New().
		Get(url).
		Set("Api-Token", c.ApiToken).
		EndBytes()

	if errs != nil {
		return r, b, CustomError{ Message: "could not perform HTTP GET request", HttpErrors: errs }
	}

	log.Printf("RESPONSE:\n%#v\n", r)
	log.Printf("BODY:\n%#v\n", b)

	var pretty bytes.Buffer
	err := json.Indent(&pretty, b, "", "\t")
	if err != nil {
		panic(err)
	}

	log.Printf("BODY(string):\n %s\n", string(pretty.Bytes()))
	log.Printf("ERRORS:\n%#v\n", errs)

	return r, b, nil
}


// Send a POST request to the Active Campaign API.
// TODO(error-check): Should check that base URL and API key are at least non-empty.
func (c *Campaigner) post(url string, i interface{}) (gorequest.Response, []byte, error) {
	// Locals.
	var (
		r gorequest.Response
		b []byte
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	// Generate URL and JSON.
	url = c.GenerateURL(url)
	j, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("Could not marshal json for interface: %s\n", err)
	}

	// Send POST request.
	r, b, errs = gorequest.New().
		Post(url).
		Send(string(j)).
		Set("Api-Token", c.ApiToken).
		EndBytes()

	// Error check.
	if errs != nil {
		return r, b, CustomError{ Message: "could not perform HTTP POST request", HttpErrors: errs }
	}

	return r, b, nil
}

// TODO(cleanup): Not being used just yet.
func (c *Campaigner) Put(url string, i interface{}) (gorequest.Response, string, error) {
	// Locals.
	var (
		r gorequest.Response
		b string
		errs []error
	)

	// Check API config.
	if err := c.CheckConfig(); err != nil {
		return r, b, err
	}

	url = c.GenerateURL(url)

	j, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("Could not marshall json for interface: %s\n", err)
	}

	r, b, errs = gorequest.New().
		Post(url).
		Send(string(j)).
		Set("Api-Token", c.ApiToken).
		End()

	// log.Printf("j: %#v", j)
	// log.Printf("r: %#v", r)
	// log.Printf("b: %#v", b)
	// log.Printf("errs: %#v", errs)

	log.Printf("\nDATA SENT:\n %#v", j)
	log.Printf("RESPONSE:\n%#v\n", r)
	log.Printf("BODY:\n%#v\n", b)

	/*
	var m map[string]interface{}
	err = json.Unmarshal([]byte(b), &m)
	if err != nil {
		panic(err)
	}
	*/

	var pretty bytes.Buffer
	err = json.Indent(&pretty, []byte(b), "", "\t")
	if err != nil {
		panic(err)
	}

	log.Printf("BODY(string):\n %s\n", string(pretty.Bytes()))
	log.Printf("ERRORS:\n%#v\n", errs)

	return r, b, nil
}