package flags

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	opts := struct {
		D Duration `short:"d"`
	}{}
	_, extraArgs, err := ParseFlags("test", &opts, []string{"test", "-d=3h"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Hour, opts.D)

	_, extraArgs, err = ParseFlags("test", &opts, []string{"test", "-d=3"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Second, opts.D)
}

func TestDurationDefault(t *testing.T) {
	opts := struct {
		D Duration `short:"d" default:"3h"`
	}{}
	_, extraArgs, err := ParseFlags("test", &opts, []string{"test"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Hour, opts.D)
}

func TestByteSize(t *testing.T) {
	opts := struct {
		S ByteSize `short:"s"`
	}{}
	_, extraArgs, err := ParseFlags("test", &opts, []string{"test", "-s=2M"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 2000000, opts.S)
}

func TestBadFlagsErrors(t *testing.T) {
	opts := struct {
		S1 ByteSize `short:"s"`
		S2 ByteSize `short:"s"`
	}{}
	_, _, err := ParseFlags("test", &opts, []string{"test"}, 0, nil, nil)
	assert.Error(t, err)
}

func TestActiveCommand(t *testing.T) {
	opts := struct {
		Build struct {
			Target string
		} `command:"build"`
		Query struct {
			Deps struct {
				Target string
			} `command:"deps"`
		} `command:"query"`
	}{}

	parser, _, err := ParseFlags("App Name", &opts, []string{"/path/to/exe", "build", "//:target"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", ActiveCommand(parser.Command))

	parser, _, err = ParseFlags("App Name", &opts, []string{"/path/to/exe", "query", "deps", "//:target"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "deps", ActiveCommand(parser.Command))
}

func TestFullActiveCommand(t *testing.T) {
	opts := struct {
		Build struct {
			Target string
		} `command:"build"`
		Query struct {
			Deps struct {
				Target string
			} `command:"deps"`
		} `command:"query"`
	}{}

	parser, _, err := ParseFlags("App Name", &opts, []string{"/path/to/exe", "build", "//:target"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", ActiveFullCommand(parser.Command))

	parser, _, err = ParseFlags("App Name", &opts, []string{"/path/to/exe", "query", "deps", "//:target"}, 0, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "query.deps", ActiveFullCommand(parser.Command))
}
