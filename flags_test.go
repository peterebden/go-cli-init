package cli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	opts := struct {
		D Duration `short:"d"`
	}{}
	_, extraArgs, err := ParseFlags("test", &opts, []string{"test", "-d=3h"}, 0, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Hour, opts.D)

	_, extraArgs, err = ParseFlags("test", &opts, []string{"test", "-d=3"}, 0, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Second, opts.D)
}

func TestDurationDefault(t *testing.T) {
	opts := struct {
		D Duration `short:"d" default:"3h"`
	}{}
	_, extraArgs, err := ParseFlags("test", &opts, []string{"test"}, 0, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(extraArgs))
	assert.EqualValues(t, 3*time.Hour, opts.D)
}
