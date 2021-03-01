package logging

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

func TestJSONFormatter(t *testing.T) {
	backend := logging.InitForTesting(logging.DEBUG)
	logging.SetFormatter(logFormatter(nil, true))
	log := logging.MustGetLogger("test_module")
	log.Infof("hello %s", "world")
	line := backend.Head().Record.Formatted(0)
	assert.Equal(t, `{"file":"logging/logging_test.go:27","func":"github.com/peterebden/go-cli-init/logging.TestJSONFormatter","level":"info","module":"test_module","msg":"hello world","time":"1970-01-01T00:00:00.000Z"}`+"\n", line)
}
