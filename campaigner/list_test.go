package campaigner

import (
	"github.com/stretchr/testify/assert"
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
	_, err := C.ListContactAdd(testListID, int64(44))
	assert.Nil(t, err)
}


