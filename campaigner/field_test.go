package campaigner

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFieldList_Success(t *testing.T) {
	_, err := C.FieldList()
	assert.Nil(t, err)
}

func TestFieldRead_Success(t *testing.T) {
	r, err := C.FieldList()
	assert.Nil(t, err)
	require.NotEmpty(t, r.Fields)

	f := r.Fields[0]

	r2, err := C.FieldRead(f.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, r2.Field.ID)
}