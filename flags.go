package cli

import (
	"fmt"
	"strings"

	"github.com/jessevdk/go-flags"
)

// A CompletionHandler is the type of function that our flags library uses to handle completions.
type CompletionHandler func(parser *flags.Parser, items []flags.Completion)

// ParseFlags parses the app's flags and returns the parser, any extra arguments, and any error encountered.
// It may exit if certain options are encountered (eg. --help).
func ParseFlags(appname string, data interface{}, args []string, opts flags.Options, completionHandler CompletionHandler) (*flags.Parser, []string, error) {
	parser := flags.NewNamedParser(path.Base(args[0]), opts)
	if completionHandler != nil {
		parser.CompletionHandler = func(items []flags.Completion) { completionHandler(parser, items) }
	}
	parser.AddGroup(appname+" options", "", data)
	extraArgs, err := parser.ParseArgs(args[1:])
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
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
	return ParseFlagsFromArgsOrDie(appname, version, data, os.Args)
}

// ParseFlagsFromArgsOrDie is similar to ParseFlagsOrDie but allows control over the
// flags passed.
// It returns the active command if there is one.
func ParseFlagsFromArgsOrDie(appname string, data interface{}, args []string) string {
	parser, extraArgs, err := ParseFlags(appname, data, args, flags.HelpFlag|flags.PassDoubleDash, nil)
	if err != nil {
		writeUsage(data)
		parser.WriteHelp(os.Stderr)
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	} else if len(extraArgs) > 0 {
		writeUsage(data)
		fmt.Printf("Unknown option %s\n", extraArgs)
		parser.WriteHelp(os.Stderr)
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