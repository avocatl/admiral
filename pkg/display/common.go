package display

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/avocatl/admiral/pkg/commander"
)

type jsonDisplayer struct {
	Data   interface{}
	Pretty bool
}

// KV is a displayable group of key value.
func (jd *jsonDisplayer) KV() []map[string]interface{} {
	var out []map[string]interface{}

	var v []byte

	var err error
	if jd.Pretty {
		v, err = json.MarshalIndent(jd.Data, "", "    ")
	} else {
		v, err = json.Marshal(jd.Data)
	}

	if err != nil {
		v = []byte(err.Error())
	}

	out = append(out, map[string]interface{}{
		"": string(v),
	})

	return out
}

// Cols returns an array of columns available for displaying.
func (jd *jsonDisplayer) Cols() []string {
	return commander.NewCols("")
}

// ColMap returns a list of columns and its description.
func (jd *jsonDisplayer) ColMap() map[string]string {
	return map[string]string{}
}

// NoHeaders returns a boolean indicating if headers should be displayed
// or not to the provided output.
func (jd *jsonDisplayer) NoHeaders() bool {
	return true
}

// Filterable defines if the displayable can be filtered (columns -> fields).
func (jd *jsonDisplayer) Filterable() bool {
	return false
}

// JSON retrieves a json displayable struct.
//
// To display prettified JSON pass true as second parameter.
// Under the hood is uses the default json marshaler.
func JSON(d interface{}, pretty bool) Displayable {
	return &jsonDisplayer{
		Data:   d,
		Pretty: pretty,
	}
}

type textDisplayable struct {
	Divider string
	Text    string
}

// KV is a displayable group of key value.
func (td *textDisplayable) KV() []map[string]interface{} {
	var out []map[string]interface{}

	br := strings.Repeat(td.Divider, 50)

	var text string
	{
		text = td.Text

		if br != "" {
			text = fmt.Sprintf("\n%s\n%s\n%s\n",
				br,
				strings.TrimSuffix(td.Text, "\n"),
				br,
			)
		}
	}

	out = append(out, map[string]interface{}{
		"": text,
	})

	return out
}

// Cols returns an array of columns available for displaying.
func (td *textDisplayable) Cols() []string {
	return commander.NewCols("")
}

// ColMap returns a list of columns and its description.
func (td *textDisplayable) ColMap() map[string]string {
	return map[string]string{}
}

// NoHeaders returns a boolean indicating if headers should be displayed
// or not to the provided output.
func (td *textDisplayable) NoHeaders() bool {
	return true
}

// Filterable defines if the displayable can be filtered (columns -> fields).
func (td *textDisplayable) Filterable() bool {
	return false
}

// Text programatically builds a simple text displayer.
// If you pass a non empty displayer  then your text will be
// wrapped by the value of divider repeated 50 times.
func Text(divider, text string) Displayable {
	return &textDisplayable{
		Divider: divider,
		Text:    text,
	}
}
