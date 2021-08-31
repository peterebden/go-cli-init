package flags

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/thought-machine/go-flags"
)

// A CompletionHandler is the type of function that our flags library uses to handle completions.
type CompletionHandler func(parser *flags.Parser, items []flags.Completion)

// AdditionalUsageInfo is a function that can seek out auxiliary flags and add
// them to the core list of options.
type AdditionalUsageInfo func(parser *flags.Parser)

// ParseFlags parses the app's flags and returns the parser, any extra arguments, and any error encountered.
// It may exit if certain options are encountered (eg. --help).
func ParseFlags(appname string, data interface{}, args []string, opts flags.Options, completionHandler CompletionHandler, additionalUsageInfo AdditionalUsageInfo) (*flags.Parser, []string, error) {
	parser := flags.NewNamedParser(path.Base(args[0]), opts)
	parser.NamespaceDelimiter = "_"
	if completionHandler != nil {
		parser.CompletionHandler = func(items []flags.Completion) { completionHandler(parser, items) }
	}
	if additionalUsageInfo != nil {
		parser.PrintAdditionalUsageInfo = func() { additionalUsageInfo(parser) }
	}
	if _, err := parser.AddGroup(appname+" options", "", data); err != nil {
		return nil, nil, err
	}
	extraArgs, err := parser.ParseArgs(args[1:])
	if err != nil {
		if t, ok := err.(*flags.Error); ok && t.Type == flags.ErrHelp {
			writeUsage(data)
			fmt.Printf("%s\n", err)
			os.Exit(0)
		}
	}
	return parser, extraArgs, err
}

// ParseFlagsOrDie parses the app's flags and dies if unsuccessful.
// Also dies if any unexpected arguments are passed.
// It returns the active command if there is one.
func ParseFlagsOrDie(appname string, data interface{}) string {
	return ParseFlagsFromArgsOrDie(appname, data, os.Args)
}

// ParseFlagsFromArgsOrDie is similar to ParseFlagsOrDie but allows control over the
// flags passed.
// It returns the active command if there is one.
func ParseFlagsFromArgsOrDie(appname string, data interface{}, args []string) string {
	parser, extraArgs, err := ParseFlags(appname, data, args, flags.HelpFlag|flags.PassDoubleDash, nil, nil)
	if err != nil && parser == nil {
		// Most likely this is something structurally wrong with the flags setup.
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else if err != nil {
		writeUsage(data)
		parser.WriteHelp(os.Stderr)
		fmt.Fprintf(os.Stderr, "\n%s\n", err)
		os.Exit(1)
	} else if len(extraArgs) > 0 {
		writeUsage(data)
		parser.WriteHelp(os.Stderr)
		fmt.Fprintf(os.Stderr, "Unknown option %s\n", extraArgs)
		os.Exit(1)
	}
	if parser.Command != nil {
		return ActiveCommand(parser.Command)
	}
	return ""
}

// ActiveCommand returns the name of the currently active command.
func ActiveCommand(command *flags.Command) string {
	if command.Active != nil {
		return ActiveCommand(command.Active)
	}
	return command.Name
}

// writeUsage prints any usage specified on the flag struct.
func writeUsage(opts interface{}) {
	if s := getUsage(opts); s != "" {
		fmt.Println(strings.TrimSpace(s) + "\n")
	}
}

// getUsage extracts any usage specified on a flag struct.
// It is set on a field named Usage, either by value or in a struct tag named usage.
func getUsage(opts interface{}) string {
	if field := reflect.ValueOf(opts).Elem().FieldByName("Usage"); field.IsValid() && field.String() != "" {
		return strings.TrimSpace(field.String())
	}
	if field, present := reflect.TypeOf(opts).Elem().FieldByName("Usage"); present {
		return field.Tag.Get("usage")
	}
	return ""
}

// A Duration is used for flags that represent a time duration; it's just a wrapper
// around time.Duration that implements the flags.Unmarshaler and
// encoding.TextUnmarshaler interfaces.
type Duration time.Duration

// UnmarshalFlag implements the flags.Unmarshaler interface.
func (d *Duration) UnmarshalFlag(in string) error {
	d2, err := time.ParseDuration(in)
	// For backwards compatibility, treat missing units as seconds.
	if err != nil {
		if d3, err := strconv.Atoi(in); err == nil {
			*d = Duration(time.Duration(d3) * time.Second)
			return nil
		}
	}
	if err != nil {
		return &flags.Error{Type: flags.ErrMarshal, Message: err.Error()}
	}
	*d = Duration(d2)
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (d *Duration) UnmarshalText(text []byte) error {
	return d.UnmarshalFlag(string(text))
}

// A ByteSize is used for flags that represent some quantity of bytes that can be
// passed as human-readable quantities (eg. "10G").
type ByteSize uint64

// UnmarshalFlag implements the flags.Unmarshaler interface.
func (b *ByteSize) UnmarshalFlag(in string) error {
	b2, err := humanize.ParseBytes(in)
	*b = ByteSize(b2)
	if err != nil {
		return &flags.Error{Type: flags.ErrMarshal, Message: err.Error()}
	}
	return nil
}
