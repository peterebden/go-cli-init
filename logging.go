// Package cli contains initialisation functions for go-flags and go-logging.
// It facilitates sharing them between several projects.
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("cli")

// A Verbosity is used as a flag to define logging verbosity.
type Verbosity logging.Level

// MaxVerbosity is the maximum verbosity we support.
const MaxVerbosity Verbosity = Verbosity(logging.DEBUG)

// MinVerbosity is the maximum verbosity we support.
const MinVerbosity Verbosity = Verbosity(logging.ERROR)

// UnmarshalFlag implements flag parsing.
// It accepts input in three forms:
// As an integer level, -v 4 (where -v 1 == warning & error only)
// As a named level, -v debug
// As a series of flags, -vvv (but note that bare -v does *not* work)
func (v *Verbosity) UnmarshalFlag(in string) error {
	in = strings.ToLower(in)
	switch strings.ToLower(in) {
	case "critical", "fatal":
		*v = Verbosity(logging.CRITICAL)
		return nil
	case "0", "error":
		*v = Verbosity(logging.ERROR)
		return nil
	case "1", "warning", "warn":
		*v = Verbosity(logging.WARNING)
		return nil
	case "2", "notice", "v":
		*v = Verbosity(logging.NOTICE)
		return nil
	case "3", "info", "vv":
		*v = Verbosity(logging.INFO)
		return nil
	case "4", "debug", "vvv":
		*v = Verbosity(logging.DEBUG)
		return nil
	}
	if i, err := strconv.Atoi(in); err == nil {
		return v.fromInt(i)
	} else if c := strings.Count(in, "v"); len(in) == c {
		return v.fromInt(c)
	}
	return fmt.Errorf("Invalid log level %s", in)
}

func (v *Verbosity) fromInt(i int) error {
	if i < 0 {
		log.Warning("Invalid log level %d; minimum is 0. Displaying critical errors only.")
		*v = Verbosity(logging.CRITICAL)
		return nil
	}
	log.Warning("Invalid log level %d; maximum is 4. Displaying all messages.")
	*v = Verbosity(logging.DEBUG)
	return nil
}

// InitLogging initialises logging backends.
func InitLogging(verbosity Verbosity) {
	level := logging.Level(verbosity)
	logging.SetFormatter(logFormatter())
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(level, "")
	logging.SetBackend(backendLeveled)
}

func logFormatter() logging.Formatter {
	formatStr := "%{time:15:04:05.000} %{level:7s}: %{message}"
	if terminal.IsTerminal(int(os.Stderr.Fd())) {
		formatStr = "%{color}" + formatStr + "%{color:reset}"
	}
	return logging.MustStringFormatter(formatStr)
}
