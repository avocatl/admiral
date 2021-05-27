package commander

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Cols contains the header for the column
// based CLI displayer.
type Cols []string

// NoCols is a 0 column indicator.
func NoCols() Cols {
	return []string{}
}

// NewCols returns a cols type.
func NewCols(vals ...string) Cols {
	return vals
}

// Config contains a command configuration.
type Config struct {
	DisableAutoGenTag     bool
	DisableSuggentions    bool
	Hidden                bool
	SuggestMinDistance    int
	LongDesc              string
	Deprecated            string
	Example               string
	Namespace             string
	ShortDesc             string
	Version               string
	Aliases               []string
	SuggestFor            []string
	ValidArgs             []string
	Execute               func(cmd *cobra.Command, args []string)
	PersistentPostHook    func(cmd *cobra.Command, args []string)
	PersistentPreHook     func(cmd *cobra.Command, args []string)
	PostHook              func(cmd *cobra.Command, args []string)
	PreHook               func(cmd *cobra.Command, args []string)
	ExecuteErr            func(cmd *cobra.Command, args []string) error
	PersistentPostHookErr func(cmd *cobra.Command, args []string) error
	PersistentPreHookErr  func(cmd *cobra.Command, args []string) error
	PostHookErr           func(cmd *cobra.Command, args []string) error
	PreHookErr            func(cmd *cobra.Command, args []string) error
	ValidArgsFunc         func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
}

// Supported flags.
const (
	StringFlag = iota
	IntFlag
	Int64Flag
	Float64Flag
	BoolFlag
)

// FlagBindOptions exposes the parameters
// used to bind a flag to a passed pointer.
type FlagBindOptions struct {
	Bound       bool
	BindInt     *int
	BindInt64   *int64
	BindString  *string
	BindBool    *bool
	BindFloat64 *float64
}

// FlagConfig defines the configuration of a flag.
type FlagConfig struct {
	FlagType   int
	Name       string
	Shorthand  string
	Usage      string
	Default    interface{}
	Required   bool
	Persistent bool
	Binding    FlagBindOptions
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
func Builder(parent *Command, config Config, cols Cols) *Command {
	cc := &cobra.Command{
		Use:                config.Namespace,
		Short:              config.ShortDesc,
		Long:               strings.TrimSpace(config.LongDesc),
		Run:                config.Execute,
		RunE:               config.ExecuteErr,
		PreRun:             config.PreHook,
		PostRun:            config.PostHook,
		Hidden:             config.Hidden,
		ValidArgs:          config.ValidArgs,
		ValidArgsFunction:  config.ValidArgsFunc,
		Example:            config.Example,
		Aliases:            config.Aliases,
		DisableAutoGenTag:  config.DisableAutoGenTag,
		DisableSuggestions: config.DisableSuggentions,
		PersistentPreRun:   config.PersistentPreHook,
		PersistentPostRun:  config.PersistentPostHook,
		PersistentPreRunE:  config.PersistentPreHookErr,
		PersistentPostRunE: config.PersistentPostHookErr,
		Version:            config.Version,
		SuggestFor:         config.SuggestFor,
	}

	if config.SuggestMinDistance > 1 {
		cc.SuggestionsMinimumDistance = config.SuggestMinDistance
	}

	c := &Command{Command: cc, cols: cols}

	if parent != nil {
		parent.AddCommand(c)
	}

	if cols := c.cols; len(cols) > 0 {
		addDisplayerFlags(c)
	}

	return c
}

func addDisplayerFlags(c *Command) {
	formatHelpText := fmt.Sprintf(
		"select displayable fields to filter the console output, possible values are %s",
		strings.Join(c.cols, ","),
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
		addIntFlag(flagger, &config)
	case Int64Flag:
		addInt64Flag(flagger, &config)
	case Float64Flag:
		addFloat64Flag(flagger, &config)
	case BoolFlag:
		addBoolFlag(flagger, &config)
	default:
		addStringFlag(flagger, &config)
	}

	if config.Required {
		err := cmd.MarkFlagRequired(config.Name)
		if err != nil {
			log.Fatal(err)
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

func addIntFlag(flagger *pflag.FlagSet, config *FlagConfig) {
	val := config.Default.(int)

	if config.Binding.Bound {
		flagger.IntVarP(config.Binding.BindInt, config.Name, config.Shorthand, val, config.Usage)
	} else {
		flagger.IntP(config.Name, config.Shorthand, val, config.Usage)
	}
}

func addInt64Flag(flagger *pflag.FlagSet, config *FlagConfig) {
	val := config.Default.(int64)

	if config.Binding.Bound {
		flagger.Int64VarP(config.Binding.BindInt64, config.Name, config.Shorthand, val, config.Usage)
	} else {
		flagger.Int64P(config.Name, config.Shorthand, val, config.Usage)
	}
}

func addFloat64Flag(flagger *pflag.FlagSet, config *FlagConfig) {
	val := config.Default.(float64)

	if config.Binding.Bound {
		flagger.Float64VarP(config.Binding.BindFloat64, config.Name, config.Shorthand, val, config.Usage)
	} else {
		flagger.Float64P(config.Name, config.Shorthand, val, config.Usage)
	}
}

func addBoolFlag(flagger *pflag.FlagSet, config *FlagConfig) {
	val := config.Default.(bool)

	if config.Binding.Bound {
		flagger.BoolVarP(config.Binding.BindBool, config.Name, config.Shorthand, val, config.Usage)
	} else {
		flagger.BoolP(config.Name, config.Shorthand, val, config.Usage)
	}
}

func addStringFlag(flagger *pflag.FlagSet, config *FlagConfig) {
	if config.Default == nil {
		config.Default = ""
	}

	val := config.Default.(string)

	if config.Binding.Bound {
		flagger.StringVarP(config.Binding.BindString, config.Name, config.Shorthand, val, config.Usage)
	} else {
		flagger.StringP(config.Name, config.Shorthand, val, config.Usage)
	}
}
