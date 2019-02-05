package campaigner

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TODO(cleanup): Remove references to globals key and baseURL

var (
	testContactID    = int64(1)
	testContactTagID = Int64json(0)
)

// Tests all contact functionality as a group.  The created contact is used by other tests.
func TestContactSuite(t *testing.T) {
	runTestWithPackagePath(t, TestContactList)
	runTestWithPackagePath(t, TestContactCreate_FailureInvalidURL)
	runTestWithPackagePath(t, TestContactCreate_Success)
	runTestWithPackagePath(t, TestContactCreate_FailureExists)
	runTestWithPackagePath(t, TestContactRead)
	runTestWithPackagePath(t, TestContactDelete_Success)
}

// Tests all contact tagging functionality as a group.  This allows the ID generated by create to be used for read and delete.
func TestContactTaggingSuite(t *testing.T) {
	runTestWithPackagePath(t, TestContactTag_ParseJsonStringsInNumbers)
	runTestWithPackagePath(t, TestTagCreate_Success)
	runTestWithPackagePath(t, TestContactCreate_Success)
	runTestWithPackagePath(t, TestContactTagCreate_FailureNotFound)
	runTestWithPackagePath(t, TestContactTagCreate_Success)
	runTestWithPackagePath(t, TestContactTagRead_FailureNotFound)
	runTestWithPackagePath(t, TestContactTagRead_Success)
	runTestWithPackagePath(t, TestContactTagDelete_FailureNotFound)
	runTestWithPackagePath(t, TestContactTagDelete_Success)
	runTestWithPackagePath(t, TestContactDelete_Success)
}

func TestContactList(t *testing.T) {
	// TODO(unit-test): Should check the internals of the ResponseContactList value.
	_, err := C.ContactList()
	assert.Nil(t, err)
}

func TestContactCreate_FailureExists(t *testing.T) {
	contact := Contact{FirstName: "Test", LastName: "User", EmailAddress: config.UnitTestEmail, PhoneNumber: config.UnitTestPhone}
	r, err := C.ContactCreate(contact)

	assert.Empty(t, r.Contact.ID)
	assert.NotNil(t, err)
}

func TestContactCreate_FailureInvalidURL(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: "http://invalid/"}
	contact := Contact{FirstName: "Test", LastName: "User", EmailAddress: config.UnitTestEmail, PhoneNumber: config.UnitTestPhone}

	r, err := c.ContactCreate(contact)

	assert.Empty(t, r.Contact.ID)
	assert.NotNil(t, err)
}

func TestContactCreate_Success(t *testing.T) {
	if len(config.UnitTestEmail) == 0 {
		t.Error("Unit test email address environment variable is empty")
		return
	}

	contact := Contact{
		FirstName:    "Test",
		LastName:     "User " + NOW,
		EmailAddress: config.UnitTestEmail,
		PhoneNumber:  config.UnitTestPhone,
	}

	contact = Contact{
		FirstName:    "Test",
		LastName:     "User",
		EmailAddress: "test@user.com",
		PhoneNumber:  "2125551212",
	}

	r, err := C.ContactCreate(contact)

	assert.NotNil(t, r)
	assert.Nil(t, err, "could not create contact: %s", err)
	testContactID = r.Contact.ID
}

func TestContactDelete_Success(t *testing.T) {
	err := C.ContactDelete(testContactID)
	assert.Nil(t, err)
}

func TestContactRead(t *testing.T) {
	r, err := C.ContactRead(testContactID)

	assert.NotNil(t, r)
	assert.Nil(t, err)
}

// TODO(unit-testing) Complete  this test.
func TestContactUpdate_Success(t *testing.T) {
	id := int64(35)
	con, err := C.ContactRead(id)
	require.Nil(t, err)

	req := RequestContactUpdate{
		ID: con.Contact.ID,
		EmailAddress: con.Contact.EmailAddress,
		FirstName: con.Contact.FirstName,
		LastName: con.Contact.LastName,
		PhoneNumber: con.Contact.PhoneNumber,
		OrganizationID: 1,
	}
	_, err = C.ContactUpdate(id, req)
	assert.Nil(t, err)

	r, err := C.ContactRead(id)
	assert.Nil(t, err)
	dump(r)
}


func TestContactFieldUpdate_Success(t *testing.T) {
	now := time.Now()
	v := fmt.Sprintf("test value %s_%s", now.Format("20060102"), now.Format("220841.000"))

	_, err := C.ContactFieldUpdate(testContactID, 2, v)
	assert.Nil(t, err)
}

// NOTE(api): Their docs say a 404 or 422 can be returned as errors but all I've gotten so far are 200, 201, and 500.
// TODO(unit-test): Need to test for CustomErrorNotFound here.
func TestContactTagCreate_FailureNotFound(t *testing.T) {
	set := [][]int64{
		{0, 0},                   // This goes right along and returns a valid json response with invalid links.
		{-1, -1},                 // Returning a 500 error now.
		{1, 2147483647},          // Asking for trouble but I doubt these will ever exist.
		{2147483647, 2147483647}, // Asking for trouble but I doubt these will ever exist.
	}

	for _, i := range set {
		request := RequestContactTagCreate{ContactID: i[0], TagID: i[1]}
		_, err := C.ContactTagCreate(request)
		assert.NotNil(t, err)
	}
}

func TestContactTagCreate_Success(t *testing.T) {
	require.NotEmpty(t, testContactTagID)
	request := RequestContactTagCreate{ContactID: testContactID, TagID: int64(testContactTagID)} // Leads, cold
	r, err := C.ContactTagCreate(request)
	assert.Nil(t, err)

	testContactTagID = r.ContactTag.ID
}

func TestContactTagDelete_FailureNotFound(t *testing.T) {
	id := int64(0)
	err := C.ContactTagDelete(id)
	assert.NotNil(t, err)
	assert.IsType(t, new(CustomErrorNotFound), err, err.Error())
}

// TODO(unit-test): Need to feed a known existing tag-relationship ID here.
func TestContactTagDelete_Success(t *testing.T) {
	assert.NotEmpty(t, testContactTagID)
	id := int64(testContactTagID)
	err := C.ContactTagDelete(id)
	assert.Nil(t, err)
}

func TestContactTagRead_FailureNotFound(t *testing.T) {
	id := int64(0)
	_, err := C.ContactTagReadByContactID(id)
	assert.NotNil(t, err)
	assert.IsType(t, new(CustomErrorNotFound), err, err.Error())
}

func TestContactTagRead_Success(t *testing.T) {
	require.NotEmpty(t, testContactTagID)

	id := int64(testContactID)
	r, err := C.ContactTagReadByContactID(id)

	assert.Nil(t, err)
	require.NotEmpty(t, r.ContactTags)
	require.NotEmpty(t, r.ContactTags[0])
	assert.NotEmpty(t, r.ContactTags[0].ID)
	assert.NotEmpty(t, r.ContactTags[0].TagID)
	assert.NotEmpty(t, r.ContactTags[0].ContactID)
}

// Tests that custom type int64json is parsing correctly.
func TestContactTag_ParseJsonStringsInNumbers(t *testing.T) {
	//noinspection SpellCheckingInspection
	var (
		response ResponseContactTagCreate
		js       = []string{
			`{"contacts":[{"cdate":"2018-11-06T11:01:00-06:00","email":"247actestuser00002@henroc.net","phone":"3015413441","firstName":"Test","lastName":"User 00002","orgid":"0","segmentio_id":"","bounced_hard":"0","bounced_soft":"0","bounced_date":"0000-00-00","ip":"0","ua":"","hash":"33e218bc9a6cf37cbdfbc2da73237dc4","socialdata_lastcheck":"0000-00-00 00:00:00","email_local":"","email_domain":"","sentcnt":"0","rating_tstamp":"0000-00-00","gravatar":"1","deleted":"0","anonymized":"0","adate":"2018-11-29T16:52:30-06:00","udate":"2018-11-29T16:53:17-06:00","deleted_at":"0000-00-00 00:00:00","created_utc_timestamp":"2018-11-06 11:01:00","updated_utc_timestamp":"2018-11-29 16:53:17","links":{"bounceLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/bounceLogs","contactAutomations":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactAutomations","contactData":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactData","contactGoals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactGoals","contactLists":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLists","contactLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLogs","contactTags":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactTags","contactDeals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactDeals","deals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/deals","fieldValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/fieldValues","geoIps":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/geoIps","notes":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/notes","organization":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/organization","plusAppend":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/plusAppend","trackingLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/trackingLogs","scoreValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/scoreValues"},"id":"3","organization":null}],"contactTag":{"contact":3,"tag":13,"cdate":"2018-11-29T16:53:17-06:00","links":{"tag":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/tag","contact":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/contact"},"id":"22"}}`,
			`{"contacts":[{"cdate":"2018-11-06T11:01:00-06:00","email":"247actestuser00002@henroc.net","phone":"3015413441","firstName":"Test","lastName":"User 00002","orgid":"0","segmentio_id":"","bounced_hard":"0","bounced_soft":"0","bounced_date":"0000-00-00","ip":"0","ua":"","hash":"33e218bc9a6cf37cbdfbc2da73237dc4","socialdata_lastcheck":"0000-00-00 00:00:00","email_local":"","email_domain":"","sentcnt":"0","rating_tstamp":"0000-00-00","gravatar":"1","deleted":"0","anonymized":"0","adate":"2018-11-29T16:52:30-06:00","udate":"2018-11-29T16:53:17-06:00","deleted_at":"0000-00-00 00:00:00","created_utc_timestamp":"2018-11-06 11:01:00","updated_utc_timestamp":"2018-11-29 16:53:17","links":{"bounceLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/bounceLogs","contactAutomations":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactAutomations","contactData":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactData","contactGoals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactGoals","contactLists":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLists","contactLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactLogs","contactTags":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactTags","contactDeals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/contactDeals","deals":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/deals","fieldValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/fieldValues","geoIps":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/geoIps","notes":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/notes","organization":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/organization","plusAppend":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/plusAppend","trackingLogs":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/trackingLogs","scoreValues":"https:\/\/247waiter.api-us1.com\/api\/3\/contacts\/3\/scoreValues"},"id":"3","organization":null}],"contactTag":{"contact":3,"tag":13,"cdate":"2018-11-29T16:53:17-06:00","links":{"tag":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/tag","contact":"https:\/\/247waiter.api-us1.com\/api\/3\/contactTags\/22\/contact"},"id":22}}`,
		}
	)

	for _, s := range js {
		err := json.Unmarshal([]byte(s), &response)
		assert.Nil(t, err)
		assert.IsType(t, int64(1), int64(response.ContactTag.ID))
	}
}
