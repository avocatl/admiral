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
		Name:       "name",
		Shorthand:  "n",
		Usage:      "who are you greeting?",
		Default:    "world",
		Persistent: true,
	})

	es := cmd.Builder(
		hello,
		cmd.Config{
			Namespace: "es",
			Execute:   runHolaAction,
		},
		cmd.NoCols,
	)
	hello.AddCommand(es)
}

func main() {
	if err := hello.Execute(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func runHelloAction(command *cobra.Command, args []string) {
	fmt.Printf("Hello, %s!\n", getNameCapitalized(command.Flags().GetString("name")))
}

func runHolaAction(command *cobra.Command, args []string) {
	fmt.Printf("Hola, %s!\n", getNameCapitalized(command.Flags().GetString("name")))
}

func getNameCapitalized(n string, err error) string {
	if err != nil {
		log.Fatal(err)
	}

	return strings.Title(n)
}
