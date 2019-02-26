package campaigner

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testListID = int64(1)
)

func TestListList_Success(t *testing.T) {
	_, err := C.ListList()
	assert.Nil(t, err)
}


func TestListContactAdd_Success(t *testing.T) {
	contact := Contact{
		FirstName:    "Test",
		LastName:     "User " + NOW,
		EmailAddress: "test_" + NOW + "@user.com",
		PhoneNumber:  config.UnitTestPhone,
	}

	c, err := C.ContactCreate(contact)
	require.Nil(t, err)

	// Test list addition for brand new contact.
	_, err = C.ListContactAdd(testListID, c.Contact.ID)
	assert.Nil(t, err)

	// Test list addition for contact already in a list.
	_, err = C.ListContactAdd(testListID, c.Contact.ID)
	assert.Nil(t, err)

	err = C.ContactDelete(c.Contact.ID)
	assert.Nil(t, err)
}


