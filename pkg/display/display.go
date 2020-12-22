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
}

// Displayer executes a display action based on
// the provided displayable object.
type Displayer interface {
	Display(Displayable) error
	DisplayMany([]Displayable) error
}

type stdDisplayer struct {
	output io.Writer
}

// Display returns an error if the action of
// printing output to the CLI fails.
func (sd *stdDisplayer) Display(d Displayable) error {
	w := newTabWritter(sd.output)

	displayablePrinter(d, w)

	return w.Flush()
}

// DisplayMany executes the displaying process on multiple
// displayable structs.
//
// A popular use case is an object with nested objects inside
// each of which requires a specific dispaying structure.
func (sd *stdDisplayer) DisplayMany(ds []Displayable) error {
	for _, d := range ds {
		err := sd.Display(d)
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

func displayablePrinter(d Displayable, w io.Writer) {
	if !d.NoHeaders() {
		fmt.Fprintln(w, strings.Join(d.Cols(), "\t"))
	}

	for _, r := range d.KV() {
		values := []interface{}{}
		formats := []string{}

		for _, col := range d.Cols() {
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
// The output appeareance is similar to the one provided
// by docker's cli.
func DefaultDisplayer(output io.Writer) Displayer {
	if output == nil {
		output = os.Stdout
	}
	return &stdDisplayer{
		output: output,
	}
}
