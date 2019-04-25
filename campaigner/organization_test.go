package campaigner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testOrganizationID = int64(1)
	testOrganizationName = "lightsaber pizza"
)

// This is run by TestMain.  Tests are called in order.
func TestOrganizationSuite(t *testing.T) {
	runTestWithPackagePath(t, TestOrganizationList_Success)
	runTestWithPackagePath(t, TestOrganizationCreate_FailureEmpty)
	runTestWithPackagePath(t, TestOrganizationCreate_Success)
	runTestWithPackagePath(t, TestOrganizationRead_Failure)
	runTestWithPackagePath(t, TestOrganizationRead_Success)
	runTestWithPackagePath(t, TestOrganizationFind_FailureNameEmpty)
	runTestWithPackagePath(t, TestOrganizationFind_Success)
	runTestWithPackagePath(t, TestOrganizationDelete_FailureNotFound)
	runTestWithPackagePath(t, TestOrganizationDelete_Success)
}

func TestOrganizationCreate_FailureEmpty(t *testing.T) {
	o := Organization{}
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	resp, err := c.OrganizationCreate(o)

	assert.NotNil(t, err)
	assert.Empty(t, resp.Organization.ID)
}

func TestOrganizationCreate_Success(t *testing.T) {
	// Setup.
	testOrganizationName = fmt.Sprintf("Test Organization %s", NOW)
	org := Organization{ Name: testOrganizationName }
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	resp, err := c.OrganizationCreate(org)

	testOrganizationID = resp.Organization.ID

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, int64(1), resp.Organization.ID)
}

func TestOrganizationDelete_FailureNotFound(t *testing.T) {
	var (
		c         = Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
		invalidID = int64(0)
	)

	err := c.OrganizationDelete(invalidID)

	assert.NotNil(t, err) // Should get an error back.
	assert.IsType(t, new(CustomErrorNotFound), err, err)
}

func TestOrganizationDelete_Success(t *testing.T) {
	var (
		c          = Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
		unexpected string
	)

	// Create a test organization if running test standalone.
	if testOrganizationID == 0 {
		testOrganizationName = fmt.Sprintf("Test Organization %s", NOW)
		org := Organization{ Name: testOrganizationName }
		resp, err := c.OrganizationCreate(org)
		assert.Nil(t, err, "could not create organization for one-off test")
		if err != nil {
			return
		}
		testOrganizationID = resp.Organization.ID
	}

	err := c.OrganizationDelete(testOrganizationID)

	if err != nil {
		unexpected = err.Error()
	}

	assert.Nil(t, err, unexpected)
}

func TestOrganizationFind_FailureNameEmpty(t *testing.T) {
	names := []string{ "", " "}

	for _, n := range names {
		_, err := C.OrganizationFind(n)
		e := fmt.Sprintf("should have gotten an error for name `%s`", n)
		assert.NotNil(t, err, e)
		assert.Equal(t, "organization find failed, name is empty", err.Error(), e)
	}
}

func TestOrganizationFind_Success(t *testing.T) {
	n := testOrganizationName

	r, err := C.OrganizationFind(n)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, r.Meta.Total)
	assert.Equal(t, 1, len(r.Organizations))
	assert.Equal(t, n, r.Organizations[0].Name)

}

func TestOrganizationList_Success(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	offset := 0
	limit :=  20
	_, err := c.OrganizationList(limit, offset)

	assert.Nil(t, err)
}

func TestOrganizationRead_Failure(t *testing.T) {
	badID := int64(-1)

	_, err := C.OrganizationRead(badID)

	assert.NotNil(t, err)
}

func TestOrganizationRead_Success(t *testing.T) {
	r, err := C.OrganizationRead(testOrganizationID)
	assert.Nil(t, err)
	assert.Equal(t, testOrganizationID, r.Organization.ID)
	assert.NotEmpty(t, r.Organization.Name)
}

func TestOrganizationUpdate_Failure(t *testing.T) {
	badID := int64(-1)
	request := RequestOrganizationUpdate{ Name: "" }

	_, err := C.OrganizationUpdate(badID, request)
	assert.NotNil(t, err)
}

func TestOrganizationUpdate_Success(t *testing.T) {
	id := int64(63)
	request := RequestOrganizationUpdate{ Name: "Org " + NOW }

	_, err := C.OrganizationUpdate(id, request)
	assert.Nil(t, err)
}
