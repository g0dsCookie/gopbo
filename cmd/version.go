package cmd

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/google/subcommands"
)

const version = "1.1.0"

// VersionCmd is used to show the current version.
type VersionCmd struct {
}

// Name returns the name of this command.
func (*VersionCmd) Name() string { return "version" }

// Synopsis returns the synopsis of this command.
func (*VersionCmd) Synopsis() string { return "Prints the tools version" }

// Usage returns the usage of this command.
func (*VersionCmd) Usage() string {
	return `version
`
}

// SetFlags registers all flags for this command.
func (*VersionCmd) SetFlags(_ *flag.FlagSet) {}

// Execute parses and executes the input.
func (*VersionCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println("gopbo-"+version, "|", runtime.Version(), "|", runtime.GOOS)
	return subcommands.ExitSuccess
}
