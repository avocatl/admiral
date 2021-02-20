# Admiral

Create [cobra](https://github.com/spf13/cobra) based command line tools without guessing so much about how to keep things consistent and testable.

Admiral provides some common base functionality such as column display filtering and display order.

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

	"github.com/avocatl/admiral/pkg/commander"
	"github.com/spf13/cobra"
)

var hello = commander.Builder(
	nil,
	commander.Config{
		Namespace: "hello",
		Execute:   runHelloAction,
	},
	commander.NoCols,
)

func init() {
	commander.AddFlag(hello, commander.FlagConfig{
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

func runHelloAction(cmd *cobra.Command, args []string) {
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
