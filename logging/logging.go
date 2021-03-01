// Package cli contains initialisation functions for go-flags and go-logging.
// It facilitates sharing them between several projects.
package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/op/go-logging.v1"
)

var log = MustGetLogger()

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
func InitLogging(verbosity Verbosity) LogLevelInfo {
	return initialiseLogging(verbosity, false)
}

func initialiseLogging(verbosity Verbosity, structured bool) LogLevelInfo {
	backend := initLogging(verbosity, os.Stderr, structured)
	logging.SetBackend(backend)
	logInfo.backend = backend
	return &logInfo
}

// InitFileLogging initialises logging backends, both to stderr and to a file.
// If the file path is empty then it will be ignored.
func InitFileLogging(stderrVerbosity, fileVerbosity Verbosity, filename string) (LogLevelInfo, error) {
	return InitStructuredLogging(stderrVerbosity, fileVerbosity, filename, false)
}

// MustInitFileLogging is like InitFileLogging but dies on any errors.
func MustInitFileLogging(stderrVerbosity, fileVerbosity Verbosity, filename string) LogLevelInfo {
	return MustInitStructuredLogging(stderrVerbosity, fileVerbosity, filename, false)
}

// InitStructuredLogging is like InitFileLogging but allows specifying whether the output should be
// structured as JSON.
func InitStructuredLogging(stderrVerbosity, fileVerbosity Verbosity, filename string, structured bool) (LogLevelInfo, error) {
	info := initialiseLogging(stderrVerbosity, structured)
	if filename == "" {
		return info, nil
	}
	if err := os.MkdirAll(path.Dir(filename), os.ModeDir|0755); err != nil {
		return info, err
	}
	f, err := os.Create(filename)
	if err != nil {
		return info, err
	}
	logging.SetBackend(logInfo.backend, initLogging(fileVerbosity, f, structured))
	return info, nil
}

// MustInitStructuredLogging is like InitStructuredLogging but dies on any errors.
func MustInitStructuredLogging(stderrVerbosity, fileVerbosity Verbosity, filename string, structured bool) LogLevelInfo {
	info, err := InitStructuredLogging(stderrVerbosity, fileVerbosity, filename, structured)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	return info
}

func initLogging(verbosity Verbosity, out *os.File, structured bool) logging.LeveledBackend {
	level := logging.Level(verbosity)
	backend := logging.NewLogBackend(out, "", 0)
	backendFormatted := logging.NewBackendFormatter(backend, logFormatter(out, structured))
	backendLeveled := logging.AddModuleLevel(backendFormatted)
	backendLeveled.SetLevel(level, "")
	return backendLeveled
}

func logFormatter(f *os.File, structured bool) logging.Formatter {
	if structured {
		return jsonFormatter{}
	}
	formatStr := "%{time:15:04:05.000} %{level:7s}: %{message}"
	if terminal.IsTerminal(int(f.Fd())) {
		formatStr = "%{color}" + formatStr + "%{color:reset}"
	}
	return logging.MustStringFormatter(formatStr)
}

// getLoggerName returns the name of the calling package as a logger name (e.g. "github.com.peterebden.cli")
func getLoggerName(skip int) string {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return "<unknown>" // Shouldn't really happen but best to handle it.
	}
	return strings.Replace(strings.TrimPrefix(path.Dir(file), ".go"), "/", ".", -1)
}

// MustGetLogger is a wrapper around go-logging's function of the same name. It automatically determines a logger name.
// The logger is registered and will be returned by ModuleLevels().
func MustGetLogger() *logging.Logger {
	return MustGetLoggerNamed(getLoggerName(2)) // 2 to skip back to the calling function.
}

// MustGetLoggerNamed is like MustGetLogger but lets the caller choose the name.
func MustGetLoggerNamed(name string) *logging.Logger {
	logInfo.Register(name)
	return logging.MustGetLogger(name)
}

// A LogLevelInfo describes and can modify levels of the set of registered loggers.
type LogLevelInfo interface {
	// ModuleLevels returns the level of all loggers retrieved by MustGetLogger().
	ModuleLevels() map[string]logging.Level
	// SetLevel modifies the level of a specific logger.
	SetLevel(level logging.Level, module string)
}

type logLevelInfo struct {
	backend logging.LeveledBackend
	modules map[string]struct{}
	mutex   sync.Mutex
}

func (info *logLevelInfo) Register(name string) {
	info.mutex.Lock()
	defer info.mutex.Unlock()
	info.modules[name] = struct{}{}
}

func (info *logLevelInfo) ModuleLevels() map[string]logging.Level {
	info.mutex.Lock()
	defer info.mutex.Unlock()
	levels := map[string]logging.Level{}
	levels[""] = info.backend.GetLevel("")
	for module := range info.modules {
		levels[module] = info.backend.GetLevel(module)
	}
	return levels
}

func (info *logLevelInfo) SetLevel(level logging.Level, module string) {
	info.backend.SetLevel(level, module)
}

var logInfo = logLevelInfo{modules: map[string]struct{}{}}

type jsonFormatter struct{}

func (f jsonFormatter) Format(calldepth int, r *logging.Record, w io.Writer) error {
	fn := ""
	pc, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	if f := runtime.FuncForPC(pc); f != nil {
		fn = f.Name()
	}
	return json.NewEncoder(w).Encode(&jsonEntry{
		File:   fmt.Sprintf("%s:%d", file, line),
		Func:   fn,
		Level:  jsonLevelNames[r.Level],
		Module: r.Module,
		Time:   r.Time.Format(jsonTimestampFormat),
		Msg:    r.Message(),
	})
}

var jsonLevelNames = []string{
	"critical",
	"error",
	"warning",
	"notice",
	"info",
	"debug",
}

type jsonEntry struct {
	File   string `json:"file"`
	Func   string `json:"func"`
	Level  string `json:"level"`
	Module string `json:"module"`
	Msg    string `json:"msg"`
	Time   string `json:"time"`
}

const jsonTimestampFormat = "2006-01-02T15:04:05.000Z07:00"
