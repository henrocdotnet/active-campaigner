package campaigner

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// Not sure how useful this map will be.  IIRC contact IDs are reused, organization IDs are not.
	testMap = map[int64]Organization{
		1: {ID: 1, Name: "lightsaber pizza"},
		2: {ID: 2, Name: "Test Organization 00001"},
		3: {ID: 3, Name: "Henroc Test 00019"},
	}

	testOrganizationID int64
)

// This is run by TestMain.  Tests are called in order.
func TestOrganizations(t *testing.T) {
	runTestWithName(t, TestOrganizationList_Success)
	runTestWithName(t, TestOrganizationCreate_FailureEmpty)
	runTestWithName(t, TestOrganizationCreate_Success)
	runTestWithName(t, TestOrganizationDelete_FailureNotFound)
	runTestWithName(t, TestOrganizationDelete_Success)
}

func TestOrganizationList_Success(t *testing.T) {
	c := Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}

	_, err := c.OrganizationList()
	if err != nil {
		log.Println(err)
	}

	assert.Nil(t, err)
}

func TestOrganizationCreate_FailureEmpty(t *testing.T) {
	o := Organization{}
	c := Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}

	resp, err := c.OrganizationCreate(o)
	if err != nil {
		log.Printf("Found expected error: %s\n", err)
	}

	assert.NotNil(t, err)
	assert.Empty(t, resp.Organization.ID)
}

func TestOrganizationCreate_Success(t *testing.T) {
	o := testMap[2]
	log.Printf("%s\n", o.Name)
	c := Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}

	resp, err := c.OrganizationCreate(o)
	if err != nil {
		log.Printf("TEST ORG CREATE ERROR: %s\n", err)
	}

	testOrganizationID = resp.Organization.ID

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, int64(1), resp.Organization.ID)
}

func TestOrganizationDelete_Success(t *testing.T) {
	var (
		c          = Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}
		err        = c.OrganizationDelete(testOrganizationID)
		unexpected string
	)

	if err != nil {
		unexpected = err.Error()
	}

	assert.Nil(t, err, unexpected)
}

func TestOrganizationDelete_FailureNotFound(t *testing.T) {
	var (
		c         = Campaigner{ApiToken: config.ApiKey, BaseURL: config.BaseURL}
		invalidID = int64(0)
	)

	err := c.OrganizationDelete(invalidID)

	assert.NotNil(t, err) // Should get an error back.
	assert.IsType(t, new(CustomErrorNotFound), err, err.Error())
}
