package main

import (
	"context"
	"flag"
	"os"

	"github.com/g0dsCookie/gopbo/cmd"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.UnpackCmd{}, "")
	subcommands.Register(&cmd.ValidateCmd{}, "")
	subcommands.Register(&cmd.PackCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
