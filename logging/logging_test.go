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
	logging.SetFormatter(logFormatter(nil, true, false, false))
	log := logging.MustGetLogger("test_module")
	log.Infof("hello %s", "world")
	line := backend.Head().Record.Formatted(0)
	assert.Equal(t, `{"file":"logging/logging_test.go:27","func":"github.com/peterebden/go-cli-init/logging.TestJSONFormatter","level":"info","module":"test_module","msg":"hello world","time":"1970-01-01T00:00:00.000Z"}`+"\n", line)
}

func TestLogAppend(t *testing.T) {
	_, err := InitLoggingOptions(&Options{})
	assert.NoError(t, err)
}

func TestStructTagsEquivalence(t *testing.T) {
	// Testing this as a user might, to customise some of the flags.
	_, err := InitLoggingOptionsLike(&struct {
		Verbosity     Verbosity `short:"v" long:"verbosity" description:"Verbosity of output (error, warning, notice, info, debug)" default:"warning"`
		File          string    `long:"log_file" description:"File to echo full logging output to" default:"plz-out/log/build.log"`
		FileVerbosity Verbosity `long:"log_file_level" description:"Log level for file output" default:"debug"`
		Append        bool      `long:"log_append" description:"Append log to existing file instead of overwriting its content. If not set, a new file will be chosen if the existing one is already open."`
		Colour        bool      `long:"colour" description:"Forces coloured output."`
		NoColour      bool      `long:"nocolour" description:"Forces colourless output."`
		Structured    bool      `long:"structured_logs" env:"STRUCTURED_LOGS" description:"Output logs in structured (JSON) format"`
	}{})
	assert.NoError(t, err)
}
