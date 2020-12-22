package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Config contains a command configuration
type Config struct {
	Aliases   []string
	Example   string
	Execute   func(cmd *cobra.Command, args []string)
	Hidden    bool
	LongDesc  string
	Namespace string
	PostHook  func(cmd *cobra.Command, args []string)
	PreHook   func(cmd *cobra.Command, args []string)
	ShortDesc string
	ValidArgs []string
}

// Supported flags
const (
	StringFlag = iota
	IntFlag
	Int64Flag
	Float64Flag
	BoolFlag
)

// FlagConfig defines the configuration of a flag.
type FlagConfig struct {
	FlagType   int
	Name       string
	Shorthand  string
	Usage      string
	Default    interface{}
	Required   bool
	Persistent bool
}

// Command wraps a base cobra command to add some
// custom functionality.
type Command struct {
	*cobra.Command
	cols     []string
	children []*Command
}

// AddCommand adds child commands and also to cobra.
func (c *Command) AddCommand(commands ...*Command) {
	c.children = append(c.children, commands...)
	for _, cmd := range commands {
		c.Command.AddCommand(cmd.Command)
	}
}

// GetSubCommands returns the command sub commands.
func (c *Command) GetSubCommands() []*Command {
	return c.children
}

// Builder constructs a new command.
func Builder(parent *Command, config Config, cols []string) *Command {
	cc := &cobra.Command{
		Use:       config.Namespace,
		Short:     config.ShortDesc,
		Long:      strings.TrimSpace(config.LongDesc),
		Run:       config.Execute,
		PreRun:    config.PreHook,
		PostRun:   config.PostHook,
		Hidden:    config.Hidden,
		ValidArgs: config.ValidArgs,
		Example:   config.Example,
		Aliases:   config.Aliases,
	}

	c := &Command{Command: cc, cols: cols}

	if parent != nil {
		parent.AddCommand(c)
	}

	if cols := c.cols; len(cols) > 0 {
		formatHelpText := fmt.Sprintf(
			"select displayable fields to filter the console output, possible values are %s",
			strings.Join(cols, ","),
		)

		AddFlag(
			c,
			FlagConfig{
				Name:       "fields",
				Persistent: true,
				Shorthand:  "f",
				Usage:      formatHelpText,
			},
		)

		AddFlag(
			c,
			FlagConfig{
				Name:       "no-headers",
				FlagType:   BoolFlag,
				Persistent: true,
				Usage:      "Return raw data with no headers",
				Default:    false,
			},
		)
	}

	return c
}

// AddFlag attaches a flag of the given type with the
// specified configuration.
func AddFlag(cmd *Command, config FlagConfig) {
	var flagger *pflag.FlagSet
	{
		if config.Persistent {
			flagger = cmd.PersistentFlags()
		} else {
			flagger = cmd.Flags()
		}
	}
	switch config.FlagType {
	case IntFlag:
		val := config.Default.(int)
		flagger.IntP(config.Name, config.Shorthand, val, config.Usage)
	case Int64Flag:
		val := config.Default.(int64)
		flagger.Int64P(config.Name, config.Shorthand, val, config.Usage)
	case Float64Flag:
		val := config.Default.(float64)
		flagger.Float64P(config.Name, config.Shorthand, val, config.Usage)
	case BoolFlag:
		val := config.Default.(bool)
		flagger.BoolP(config.Name, config.Shorthand, val, config.Usage)
	default:
		if config.Default == nil {
			config.Default = ""
		}
		val := config.Default.(string)
		flagger.StringP(config.Name, config.Shorthand, val, config.Usage)
	}

	if config.Required {
		err := cmd.MarkFlagRequired(config.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}

// FilterColumns will check if the filterable flag is used
// and return only the requested set of columns.
//
// It takes the string parsed from the filterable flag as
// first argument and the default set of columns (default)
// as second parameter.
func FilterColumns(req string, def []string) []string {
	if req != "" {
		return strings.Split(req, ",")
	}

	return def
}
