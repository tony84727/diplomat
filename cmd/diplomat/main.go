package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(&createCommand{}, "")
	subcommands.Register(&generateCommand{}, "")
	flag.Parse()
	os.Exit(int(subcommands.Execute(context.Background(), os.Args)))
}
