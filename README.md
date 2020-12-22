# Admiral

Create [cobra](https://github.com/spf13/cobra) based command line tools without guessing so much about how to keep things consistent and testable.

Admiral provides some common base functionality such as column display filtering and sorting.

## Usage

### Installation

```bash
go get -u github.com/avocatl/admiral
```

### Usage

```go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/avocatl/admiral/pkg/cmd"
	"github.com/spf13/cobra"
)

var hello = cmd.Builder(
	nil,
	cmd.Config{
		Namespace: "hello",
		Execute:   runHelloAction,
	},
	cmd.NoCols,
)

func init() {
	cmd.AddFlag(hello, cmd.FlagConfig{
		Name:      "name",
		Shorthand: "n",
		Usage:     "who are you greeting?",
		Default:   "world",
	})
}

func main() {
	if err := hello.Execute(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func runHelloAction(command *cobra.Command, args []string) {
	name, err := command.Flags().GetString("name")
	{
		name = strings.Title(name)
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello, %s!\n", name)
}
```

**an example with subcommands is available at [hello_world](examples/hello_world/main.go)**