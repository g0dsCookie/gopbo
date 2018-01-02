package cmd

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/google/subcommands"
)

var (
	buildTime string
	gitHash   string
	gitBranch string
)

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
	const format = "%10s %s\n"
	fmt.Printf(format, "Built", buildTime)
	fmt.Printf(format, "Git Hash", gitHash)
	fmt.Printf(format, "Git Branch", gitBranch)
	fmt.Printf(format, "Runtime", runtime.Version())
	fmt.Printf(format, "OS", runtime.GOOS)
	return subcommands.ExitSuccess
}
