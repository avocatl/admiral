package commander

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	config := Config{
		Namespace: "test",
		Example:   "test <name>",
		Hidden:    false,
	}

	c := Builder(nil, config, []string{"test1", "col2"})

	assert.Equal(t, c.Command.Name(), config.Namespace)
	assert.False(t, c.Hidden)
	assert.Equal(t, c.Example, config.Example)
}

func TestBuilder_WithParent(t *testing.T) {
	cp := Config{
		Namespace: "parent",
		Example:   "parent <name>",
	}

	p := Builder(nil, cp, []string{})

	child := Builder(p, Config{Namespace: "child"}, []string{})

	assert.Equal(t, p.Command, child.Parent())
	assert.Nil(t, p.Parent())
	assert.ElementsMatch(t, []*Command{child}, p.GetSubCommands())
}

func TestGetCols(t *testing.T) {
	cols := []string{"test"}
	c := Builder(nil, Config{Namespace: "test"}, cols)

	assert.Equal(t, cols, c.cols)
}

func TestAddFlag_String(t *testing.T) {
	cmd := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(cmd, FlagConfig{
		FlagType:  StringFlag,
		Name:      "test-flag",
		Shorthand: "t",
		Default:   "default",
	})

	flag, err := cmd.Flags().GetString("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, "default", flag)
	assert.IsType(t, "", flag)
}

func TestAddFlag_Int(t *testing.T) {
	cmd := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(cmd, FlagConfig{
		FlagType:  IntFlag,
		Name:      "test-flag",
		Shorthand: "t",
		Default:   10,
	})

	flag, err := cmd.Flags().GetInt("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, 10, flag)
	assert.IsType(t, 0, flag)
}

func TestAddFlag_Int64(t *testing.T) {
	cmd := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(cmd, FlagConfig{
		FlagType:  Int64Flag,
		Name:      "test-flag",
		Shorthand: "t",
		Default:   int64(10),
	})

	flag, err := cmd.Flags().GetInt64("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, int64(10), flag)
	assert.IsType(t, int64(10), flag)
}

func TestAddFlag_Bool(t *testing.T) {
	cmd := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(cmd, FlagConfig{
		FlagType:  BoolFlag,
		Name:      "test-flag",
		Shorthand: "t",
		Default:   true,
	})

	flag, err := cmd.Flags().GetBool("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, true, flag)
	assert.IsType(t, false, flag)
}

func TestAddFlag_Float64(t *testing.T) {
	cmd := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(cmd, FlagConfig{
		FlagType:  Float64Flag,
		Name:      "test-flag",
		Shorthand: "t",
		Default:   float64(10.15),
	})

	flag, err := cmd.Flags().GetFloat64("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, float64(10.15), flag)
	assert.IsType(t, float64(0), flag)
}

func TestAddFlag_Persistent(t *testing.T) {
	parent := Builder(nil, Config{Namespace: "test"}, []string{})

	AddFlag(parent, FlagConfig{
		FlagType:   BoolFlag,
		Name:       "test-flag",
		Shorthand:  "t",
		Default:    true,
		Persistent: true,
	})

	child := Builder(parent, Config{Namespace: "child-test"}, []string{})

	flag, err := child.Parent().PersistentFlags().GetBool("test-flag")
	assert.Nil(t, err)
	assert.Equal(t, true, flag)
	assert.IsType(t, false, flag)
}

func TestAddFlag_BindIntFlag(t *testing.T) {
	parent := Builder(nil, Config{Namespace: "test"}, []string{})

	var fint int
	AddFlag(parent, FlagConfig{
		FlagType: IntFlag,
		Name:     "bind-int",
		Default:  10,
		Binding: FlagBindOptions{
			Bound:   true,
			BindInt: &fint,
		},
	})

	flag, err := parent.Flags().GetInt("bind-int")

	assert.Nil(t, err)
	assert.Equal(t, 10, flag)
	assert.Equal(t, fint, flag)
	assert.IsType(t, 10, flag)
}

func TestAddFlag_BindStringFlag(t *testing.T) {
	parent := Builder(nil, Config{Namespace: "test"}, []string{})

	var x string
	AddFlag(parent, FlagConfig{
		Name:    "bind-string",
		Default: "bound to var",
		Binding: FlagBindOptions{
			Bound:      true,
			BindString: &x,
		},
	})

	flag, err := parent.Flags().GetString("bind-string")

	assert.Nil(t, err)
	assert.Equal(t, "bound to var", flag)
	assert.Equal(t, x, flag)
	assert.IsType(t, "x", flag)
}

func TestAddFlag_BindFloat64Flag(t *testing.T) {
	parent := Builder(nil, Config{Namespace: "test"}, []string{})

	var x float64
	y := float64(10.1)
	AddFlag(parent, FlagConfig{
		FlagType: Float64Flag,
		Name:     "bind-float",
		Default:  y,
		Binding: FlagBindOptions{
			Bound:       true,
			BindFloat64: &x,
		},
	})

	flag, err := parent.Flags().GetFloat64("bind-float")

	assert.Nil(t, err)
	assert.Equal(t, y, flag)
	assert.Equal(t, x, flag)
	assert.IsType(t, y, flag)
}

func TestAddFlag_BindBoolFlag(t *testing.T) {
	parent := Builder(nil, Config{Namespace: "test"}, []string{})

	var x bool
	y := false
	AddFlag(parent, FlagConfig{
		FlagType: BoolFlag,
		Name:     "bind-bool",
		Default:  y,
		Binding: FlagBindOptions{
			Bound:    true,
			BindBool: &x,
		},
	})

	flag, err := parent.Flags().GetBool("bind-bool")

	assert.Nil(t, err)
	assert.Equal(t, y, flag)
	assert.Equal(t, x, flag)
	assert.IsType(t, y, flag)
}

func TestBuilder_ColsAndHeadersFlagAddition(t *testing.T) {
	cases := []struct {
		name  string
		cols  []string
		check string
	}{
		{
			"test fields flag is nil when no columns are provided",
			[]string{},
			"fields",
		},
		{
			"test no-headers flag is nil when no columns are provided",
			[]string{},
			"no-headers",
		},
		{
			"test fields flag is not nil when columns are provided",
			[]string{"name", "surname"},
			"fields",
		},
		{
			"test fields flag is not nil when columns are provided",
			[]string{"name", "surname"},
			"fields",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Builder(nil, Config{Namespace: "test"}, tt.cols)

			got := cmd.PersistentFlags().Lookup(tt.check)

			if len(tt.cols) > 0 {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
