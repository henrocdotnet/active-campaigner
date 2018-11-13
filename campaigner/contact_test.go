package campaigner

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type configSetup struct {
	ApiKey string `envconfig:"api_token"`
	BaseURL string `envconfig:"base_url"`
}

// TODO(cleanup): Remove references to globals key and baseUrl

var (
	baseUrl = "https://247waiter.api-us1.com"
	baseUrlInvalid = "http://localhost:9000"
	key = "a014448a70f26e40eaf0af44810e33d3d0a31d83056bd69cdc1c004be05d9f2af636f00b"
	contactCreatedID = ""

	i = 4
	users = map[int]TestContact{
		2: { ID: 2, FirstName: "Test", LastName: "User 00001", EmailAddress: "247actestuser00001@henroc.net", PhoneNumber: "3015413441" },
		3: { ID: 3, FirstName: "Test", LastName: "User 00002", EmailAddress: "247actestuser00002@henroc.net", PhoneNumber: "3015413441" },
		4: { ID: 4, FirstName: "Test", LastName: "User 00003", EmailAddress: "247actestuser00003@henroc.net", PhoneNumber: "3015413441" },
		5: { ID: 5, FirstName: "Test", LastName: "User 00004", EmailAddress: "247actestuser00004@henroc.net", PhoneNumber: "3015413441" },
	}

	config configSetup
)


// TODO: Move flag parsing into TestMain at some point.
func init() {
	err := envconfig.Process("ac", &config)
	if err != nil {
		log.Fatal(err)
	}
}


type TestContact struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	PhoneNumber  string
}

func TestContactList(t *testing.T) {
	c := Campaigner{ ApiToken: config.ApiKey, BaseURL: config.BaseURL }
	c.ContactList()
}


func TestContactRead(t *testing.T) {
	log.Printf("config key? %s\n", config.ApiKey)

	c := Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}
	r, err := c.ContactRead(2)

	assert.NotNil(t, r)
	assert.Nil(t, err)

	logFormattedJson("TestContactRead: ", r)

	// writeIndentedJson("contact_read_response.json", data)
}


func TestContactCreate(t *testing.T) {
	id := 5
	c := Campaigner{ ApiToken: config.ApiKey, BaseURL: config.BaseURL }
	contact := Contact{ FirstName: users[id].FirstName, LastName: users[id].LastName, EmailAddress: users[id].EmailAddress, PhoneNumber: users[id].PhoneNumber }

	r, err := c.ContactCreate(contact)

	// log.Printf("TestContactCreate: error? %s\n", err)
	// log.Printf("TestContactCreate: error? %#v\n", err)

	assert.NotNil(t, r)
	assert.Nil(t, err, "could not create contact: %s", err)
	contactCreatedID = r.Contact.ID
}


// TODO(unit-test): Complete test.
func TestContactCreateFailureExists(t *testing.T) {
	return

	c := Campaigner{ ApiToken: config.ApiKey, BaseURL: config.BaseURL }

	contact := Contact{ FirstName: "Light", LastName: "Saber", EmailAddress: "lightsabervc@gmail.com", PhoneNumber: "" }

	c.ContactCreate(contact)
}


// TODO(unit-test): Test this test.
func TestContactCreateFailureInvalidURL(t *testing.T) {
	c := Campaigner{ ApiToken: config.ApiKey, BaseURL: "http://invalid" }
	contact := Contact{ FirstName: "Henry", LastName: "Rivera", EmailAddress: "h@247waiter.com", PhoneNumber: "3015413441" }

	r, err := c.ContactCreate(contact)

	assert.Empty(t, r.Contact.ID)
	assert.NotNil(t, err)

	log.Printf("TestContactCreateFailureInvalidURL: %s", err)
}


func TestContactDelete(t *testing.T) {
	log.Printf("config? %s %s\n", config.BaseURL, config.ApiKey)
	id := 2
	c := Campaigner{ ApiToken: config.ApiKey, BaseURL: config.BaseURL }
	err := c.ContactDelete(id)
	assert.Nil(t, err)
}
