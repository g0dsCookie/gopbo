package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/g0dsCookie/gopbo/pbo"
	"github.com/google/subcommands"
)

// UnpackCmd is used to unpack pbo files into directories.
type UnpackCmd struct {
	verbose     bool
	deleteAfter bool
	destination string
}

// Name returns the name of this command.
func (*UnpackCmd) Name() string { return "unpack" }

// Synopsis returns the synopsis of this command.
func (*UnpackCmd) Synopsis() string { return "Unpacks a PBO" }

// Usage returns the usage of this command.
func (*UnpackCmd) Usage() string {
	return `unpack [-verbose] [-delete] [-destination <dir>] <pbo>
`
}

// SetFlags registers all flags for this command.
func (u *UnpackCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&u.verbose, "verbose", false, "be verbose")
	f.BoolVar(&u.deleteAfter, "delete", false, "delete pbo after successful unpacking")
	f.StringVar(&u.destination, "destination", "", "set destination where to unpack the pbo. Defaults to <pbo> without .pbo extension")
}

// Execute parses and executes the input.
func (u *UnpackCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)
	if path == "" {
		fmt.Printf("Please provide the pbo to unpack.\n\n")
		return subcommands.ExitUsageError
	}

	if u.destination == "" {
		u.destination = path[:len(path)-4]
	}

	if err := pbo.Unpack(path, u.destination, u.verbose); err != nil {
		fmt.Printf("An error occurred while unpacking PBO: %v\n", err)
		return subcommands.ExitFailure
	}

	if u.deleteAfter {
		if err := os.Remove(path); err != nil {
			fmt.Printf("Could not remove PBO: %v\n", err)
			return subcommands.ExitFailure
		}
		if u.verbose {
			fmt.Printf("%s deleted\n", path)
		}
	}

	return subcommands.ExitSuccess
}
