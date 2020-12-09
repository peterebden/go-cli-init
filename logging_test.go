package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/op/go-logging.v1"
)

func TestParseVerbosity(t *testing.T) {
	var v Verbosity
	assert.NoError(t, v.UnmarshalFlag("error"))
	assert.EqualValues(t, logging.ERROR, v)
	assert.NoError(t, v.UnmarshalFlag("1"))
	assert.EqualValues(t, logging.WARNING, v)
	assert.NoError(t, v.UnmarshalFlag("v"))

	assert.EqualValues(t, logging.NOTICE, v)
	assert.Error(t, v.UnmarshalFlag("blah"))
}
