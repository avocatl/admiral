package prompter

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

const TagName = "prompter"

// String creates a simple prompter that returns the value
// as a string.
//
// It accepts a label and a default value.
func String(q string, def string) (string, error) {
	p := new(q)

	if def != "" {
		p.Default = def
	}

	r, err := p.Run()
	if err != nil {
		return "", err
	}

	return r, nil
}

// Int creates a simple prompter that returns the value
// as an int.
//
// It accepts a label.
func Int(q string) (i int, err error) {
	p := new(q)

	is, err := p.Run()
	if err != nil {
		return
	}

	i, err = strconv.Atoi(is)
	if err != nil {
		return
	}

	return
}

func new(q string) *promptui.Prompt {
	return &promptui.Prompt{
		Label: q,
	}
}

// Struct fills the provided struct using reflection.
// The implementation will use the field name to prompt the
// user for a value.
//
// The implementation is basic and aimed for problems that
// require filling multiple fields on large structs where
// using tags will not be feasible nor user friendly.
//
// If you need to skip a piece of the struct, you can use the
// prompter tag with '-' as a value, similar to what happens
// with the encoding/json tag.
func Struct(s interface{}) (sc interface{}, err error) {
	mu := reflect.ValueOf(s).Elem()
	tm := reflect.TypeOf(s).Elem()

	for i := 0; i < mu.NumField(); i++ {
		v := mu.Field(i)
		f := tm.Field(i)

		// parse struct params
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}

		tag := f.Tag.Get(TagName)
		if tag == "-" {
			continue
		}

		err = checkAndSetValue(v, f)
		if err != nil {
			return
		}
	}

	return s, err
}

var errUnsetable = errors.New("the provided reflect value is not mutable")

func checkAndSetValue(v reflect.Value, f reflect.StructField) error {
	if !v.CanSet() {
		return fmt.Errorf("immutable value error: %w", errUnsetable)
	}

	p := createPrompt(f.Name, v.Kind(), v.String())

	r, err := p.Run()
	if err != nil {
		if !p.IsConfirm {
			return err
		}
	}

	if r == "y" && err == nil {
		r = "true"
	} else if r == "n" {
		r = "false"
	}

	switch v.Kind() {
	case reflect.Int:
		ri, e := strconv.ParseInt(r, 10, 64)
		if e != nil {
			return e
		}

		v.SetInt(ri)
	case reflect.Bool:
		ri, e := strconv.ParseBool(r)
		if e != nil {
			return e
		}

		v.SetBool(ri)
	default:
		v.SetString(r)
	}

	return nil
}

func createPrompt(name string, k reflect.Kind, def string) promptui.Prompt {
	p := promptui.Prompt{
		Label:       fmt.Sprintf("Define a value for %s", space(name)),
		HideEntered: false,
	}

	if k == reflect.Bool {
		p.Label = space(name)
		p.IsConfirm = true
	}

	return p
}

func space(s string) string {
	check := regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

	var a []string

	for _, sub := range check.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}

		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}

	return strings.ToLower(strings.Join(a, " "))
}
