package campaigner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestTagCreate_Failure(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	// Missing name.
	tag := Tag{ Description: "Test Tag Description" }
	_, err := c.TagCreate(tag)
	assert.NotNil(t, err)

	// Missing description.
	tag = Tag{ Name: "Test Tag" }
	_, err = c.TagCreate(tag)
	assert.NotNil(t, err)

	// Missing type.
	tag = Tag{ Name: "Test Tag", Description: "Test Tag Description" }
	_, err = c.TagCreate(tag)
	assert.NotNil(t, err)

	// Incorrect type.
	tag = Tag{ Name: "Test Tag", Description: "Test Tag Description", Type: "invalid" }
	_, err = c.TagCreate(tag)
	assert.NotNil(t, err)
}

func TestTagCreate_Success(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("Timestamp: %s_%s", now.Format("20060102"), now.Format("220841"))

	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	tag := Tag{ Name: fmt.Sprintf("Test Tag %s", timestamp), Description: "Test Tag Description", Type: "contact" }

	r, err := c.TagCreate(tag)
	dump(r)

	testContactTagID = int64json(r.Tag.ID)

	log.Printf("testContactTagID: %d\n", testContactID)

	assert.Nil(t, err)
}


// TODO(testing): ID needs to come from somewhere else.
func TestTagDelete_Success(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	id := int64(25)

	err := c.TagDelete(id)

	assert.Nil(t, err)
}

func TestTagFind_Success(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
	// n := "DELIVERY SERVICE"
	n := "Onboarding - Account Activated"

	r, err := c.TagFind(n)
	dump(r)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, r.Meta.Total)
	assert.Equal(t, 1, len(r.Tags))
	assert.Equal(t, n, r.Tags[0].Name)
}

func TestTagList(t *testing.T) {
	c := Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	r, err := c.TagList()
	if err != nil {
		log.Println(err)
	}

	dump(r)

	assert.Nil(t, err)
}

func TestTagRead_FailureNotFound(t *testing.T) {
	_, err := C.TagRead(2147483647)
	assert.NotNil(t, err)
	assert.IsType(t, CustomErrorNotFound{}, err, err.Error())
}

func TestTagRead_Success(t *testing.T) {
	_, err := C.TagRead(1)
	assert.Nil(t, err)
}
