package campaigner

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)


// TODO(cleanup): Remove references to globals key and baseURL

var (
	baseURL          = "https://247waiter.api-us1.com"
	baseURLInvalid   = "http://localhost:9000"
	key              = "a014448a70f26e40eaf0af44810e33d3d0a31d83056bd69cdc1c004be05d9f2af636f00b"
	contactCreatedID int64

	i     = 4
	users = map[int]TestContact{
		2: {ID: 2, FirstName: "Test", LastName: "User 00001", EmailAddress: "247actestuser00001@henroc.net", PhoneNumber: "3015413441"},
		3: {ID: 3, FirstName: "Test", LastName: "User 00002", EmailAddress: "247actestuser00002@henroc.net", PhoneNumber: "3015413441"},
		4: {ID: 4, FirstName: "Test", LastName: "User 00003", EmailAddress: "247actestuser00003@henroc.net", PhoneNumber: "3015413441"},
		5: {ID: 5, FirstName: "Test", LastName: "User 00004", EmailAddress: "247actestuser00004@henroc.net", PhoneNumber: "3015413441"},
	}

)

type TestContact struct {
	ID           int
	FirstName    string
	LastName     string
	EmailAddress string
	PhoneNumber  string
}
func TestContactList(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	c.ContactList()
}

func TestContactRead(t *testing.T) {
	log.Printf("config key? %s\n", config.APIToken)

	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	r, err := c.ContactRead(2)

	assert.NotNil(t, r)
	assert.Nil(t, err)

	logFormattedJSON("TestContactRead: ", r)

	// writeIndentedJson("contact_read_response.json", data)
}

func TestContactCreate(t *testing.T) {
	id := 5
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	contact := Contact{FirstName: users[id].FirstName, LastName: users[id].LastName, EmailAddress: users[id].EmailAddress, PhoneNumber: users[id].PhoneNumber}

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

	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	contact := Contact{FirstName: "Light", LastName: "Saber", EmailAddress: "lightsabervc@gmail.com", PhoneNumber: ""}

	c.ContactCreate(contact)
}

// TODO(unit-test): Test this test.
func TestContactCreateFailureInvalidURL(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: "http://invalid"}
	contact := Contact{FirstName: "Henry", LastName: "Rivera", EmailAddress: "h@247waiter.com", PhoneNumber: "3015413441"}

	r, err := c.ContactCreate(contact)

	assert.Empty(t, r.Contact.ID)
	assert.NotNil(t, err)

	log.Printf("TestContactCreateFailureInvalidURL: %s", err)
}

func TestContactDelete(t *testing.T) {
	log.Printf("config? %s %s\n", config.BaseURL, config.APIToken)
	id := 2
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	err := c.ContactDelete(id)
	assert.Nil(t, err)
}

// NOTE(API): Their docs say a 404 or 422 can be returned as errors but all I've gotten so far are 201 and 500.
func TestContactTag_FailureNotFound(t *testing.T) {
	set := [][]int64{
		{ 0, 0 }, // This goes right along and returns a valid json response with invalid links.
		{ -1, -1 }, // Returning a 500 error now.
		{ 1, 2147483647 }, // Asking for trouble but I doubt these will ever exist.
		{ 2147483647, 2147483647 }, // Asking for trouble but I doubt these will ever exist.
	}

	for _, i := range set {
		task := TaskContactTag{ ContactID: i[0], TagID: i[1] }

		c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
		r, err := c.ContactTag(task)
		assert.NotNil(t, err)

		log.Println("DUMP ERROR?")
		dump(err)

		log.Println("DUMP RESPONSE?")
		dump(r)
	}
}

func TestContactTag_Success(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	task := TaskContactTag{ ContactID: 3, TagID: 14 } // Leads, cold
	r, err := c.ContactTag(task)
	assert.Nil(t, err)

	dump(r)
}

func TestContactTag_ParseJsonStingsInNumbers(t *testing.T) {
	var (
		response ResponseContactTag
		js = []string{
			`{"contacts":[{"cdate":"2018-11-06T11:01:00-06:00","email":"247actestuser00002@henroc.net","phone":"3015413441","firstName":"Test","lastName":"User 00002","orgid":"0","segmentio_id":"","bounced_hard":"0","bounced_soft":"0","bounced_date":"0000-00-00","ip":"0","ua":"","hash":"33e218bc9a6cf37cbdfbc2da73237dc4","socialdata_lastcheck":"0000-00-00 00:00:00","email_local":"","email_domain":"","sentcnt":"0","rating_tstamp":"0000-00-00","gravatar":"1","deleted":"0","anonymized":"0","adate":"2018-11-29T16:52:30-06:00","udate":"2018-11-29T16:53:17-06:00","deleted_at":"0000-00-00 00:00:00","created_utc_timestamp":"2018-11-06 11:01:00","updated_utc_timestamp":"2018-11-29 16:53:17","links":{"bounceLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/bounceLogs","contactAutomations":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactAutomations","contactData":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactData","contactGoals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactGoals","contactLists":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLists","contactLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLogs","contactTags":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactTags","contactDeals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactDeals","deals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/deals","fieldValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/fieldValues","geoIps":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/geoIps","notes":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/notes","organization":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/organization","plusAppend":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/plusAppend","trackingLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/trackingLogs","scoreValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/scoreValues"},"id":"3","organization":null}],"contactTag":{"contact":3,"tag":13,"cdate":"2018-11-29T16:53:17-06:00","links":{"tag":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/tag","contact":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/contact"},"id":"22"}}`,
			`{"contacts":[{"cdate":"2018-11-06T11:01:00-06:00","email":"247actestuser00002@henroc.net","phone":"3015413441","firstName":"Test","lastName":"User 00002","orgid":"0","segmentio_id":"","bounced_hard":"0","bounced_soft":"0","bounced_date":"0000-00-00","ip":"0","ua":"","hash":"33e218bc9a6cf37cbdfbc2da73237dc4","socialdata_lastcheck":"0000-00-00 00:00:00","email_local":"","email_domain":"","sentcnt":"0","rating_tstamp":"0000-00-00","gravatar":"1","deleted":"0","anonymized":"0","adate":"2018-11-29T16:52:30-06:00","udate":"2018-11-29T16:53:17-06:00","deleted_at":"0000-00-00 00:00:00","created_utc_timestamp":"2018-11-06 11:01:00","updated_utc_timestamp":"2018-11-29 16:53:17","links":{"bounceLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/bounceLogs","contactAutomations":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactAutomations","contactData":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactData","contactGoals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactGoals","contactLists":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLists","contactLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLogs","contactTags":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactTags","contactDeals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactDeals","deals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/deals","fieldValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/fieldValues","geoIps":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/geoIps","notes":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/notes","organization":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/organization","plusAppend":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/plusAppend","trackingLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/trackingLogs","scoreValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/scoreValues"},"id":"3","organization":null}],"contactTag":{"contact":3,"tag":13,"cdate":"2018-11-29T16:53:17-06:00","links":{"tag":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/tag","contact":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/contact"},"id":22}}`,
		}
	)

	for x, s := range js {
		log.Println(x)
		err := json.Unmarshal([]byte(s), &response)
		assert.Nil(t, err)
		assert.IsType(t, int64(1), int64(response.ContactTag.ID))
	}
}
