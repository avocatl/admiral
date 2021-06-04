package display

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// Displayable describes an output handler for
// admiral.
type Displayable interface {
	KV() []map[string]interface{}
	Cols() []string
	ColMap() map[string]string
	NoHeaders() bool
	Filterable() bool
}

// Displayer executes a display action based on
// the provided displayable object.
type Displayer interface {
	Display(Displayable, []string) error
	DisplayMany([]Displayable, []string) error
}

type stdDisplayer struct {
	output io.Writer
}

// Display returns an error if the action of
// printing output to the CLI fails.
func (sd *stdDisplayer) Display(d Displayable, f []string) error {
	w := newTabWritter(sd.output)

	displayablePrinter(d, w, f)

	return w.Flush()
}

// DisplayMany executes the displaying process on multiple
// displayable structs.
//
// A popular use case is an object with nested objects inside
// each of which requires a specific dispaying structure.
func (sd *stdDisplayer) DisplayMany(ds []Displayable, f []string) error {
	for _, d := range ds {
		err := sd.Display(d, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func newTabWritter(output io.Writer) *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(output, 0, 0, 4, ' ', 0)

	return w
}

func displayablePrinter(d Displayable, w io.Writer, f []string) {
	var cols []string
	{
		cols = d.Cols()
		if len(f) > 0 && d.Filterable() {
			cols = f
		}
	}

	if !d.NoHeaders() {
		fmt.Fprintln(w, strings.Join(cols, "\t"))
	}

	for _, r := range d.KV() {
		values := []interface{}{}
		formats := []string{}

		for _, col := range cols {
			v := r[col]
			values = append(values, v)

			switch v.(type) {
			case string:
				formats = append(formats, "%s")
			case int:
				formats = append(formats, "%d")
			case float64:
				formats = append(formats, "%f")
			case bool:
				formats = append(formats, "%v")
			default:
				formats = append(formats, "%v")
			}
		}

		format := strings.Join(formats, "\t")

		fmt.Fprintf(w, format+"\n", values...)
	}
}

// DefaultDisplayer constructs a column based output
// to the provided writer (defaults to os.Stdout).
//
// The output appearance is similar to the one provided
// by docker's cli.
func DefaultDisplayer(output io.Writer) Displayer {
	if output == nil {
		output = os.Stdout
	}

	return &stdDisplayer{
		output: output,
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
