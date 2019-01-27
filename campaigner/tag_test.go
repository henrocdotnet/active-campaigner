package campaigner

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

// Tests all tag functionality as a group.  The created contact is used by other tests.
func TestTagSuite(t *testing.T) {
	runTestWithPackagePath(t, TestTagList)
	runTestWithPackagePath(t, TestTagCreate_Failure)
	runTestWithPackagePath(t, TestTagCreate_Success)
	runTestWithPackagePath(t, TestTagFind_Success)
	runTestWithPackagePath(t, TestTagRead_FailureNotFound)
	runTestWithPackagePath(t, TestTagRead_Success)
	runTestWithPackagePath(t, TestTagDelete_Success)
}

func TestTagCreate_Failure(t *testing.T) {
	// Missing name.
	tag := Tag{ Description: "Test Tag Description" }
	_, err := C.TagCreate(tag)
	assert.NotNil(t, err)

	// Missing description.
	tag = Tag{ Name: "Test Tag" }
	_, err = C.TagCreate(tag)
	assert.NotNil(t, err)

	// Missing type.
	tag = Tag{ Name: "Test Tag", Description: "Test Tag Description" }
	_, err = C.TagCreate(tag)
	assert.NotNil(t, err)

	// Incorrect type.
	tag = Tag{ Name: "Test Tag", Description: "Test Tag Description", Type: "invalid" }
	_, err = C.TagCreate(tag)
	assert.NotNil(t, err)
}

func TestTagCreate_Success(t *testing.T) {
	now := time.Now()
	timestamp := fmt.Sprintf("Timestamp: %s_%s", now.Format("20060102"), now.Format("220841.000"))
	tag := Tag{ Name: fmt.Sprintf("Test Tag %s", timestamp), Description: "Test Tag Description", Type: "contact" }

	r, err := C.TagCreate(tag)

	testContactTagID = int64json(r.Tag.ID)

	assert.Nil(t, err)
}


// TODO(testing): ID needs to come from somewhere else.
func TestTagDelete_Success(t *testing.T) {
	id := testContactTagID

	err := C.TagDelete(int64(id))

	assert.Nil(t, err)
}

func TestTagFind_Success(t *testing.T) {
	//noinspection SpellCheckingInspection
	n := "Onboarding - Account Activated"

	r, err := C.TagFind(n)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, r.Meta.Total)
	require.Equal(t, 1, len(r.Tags))
	assert.Equal(t, n, r.Tags[0].Name)
}

func TestTagList(t *testing.T) {
	r, err := C.TagList()
	if err != nil {
		log.Println(err)
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, r.Tags)
}

func TestTagRead_FailureNotFound(t *testing.T) {
	_, err := C.TagRead(2147483647)
	assert.NotNil(t, err)
	assert.IsType(t, CustomErrorNotFound{}, err, err.Error())
}

func TestTagRead_Success(t *testing.T) {
	_, err := C.TagRead(int64(testContactTagID))
	assert.Nil(t, err)
}
