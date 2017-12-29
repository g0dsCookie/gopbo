package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/g0dsCookie/gopbo/pbo"
	"github.com/google/subcommands"
)

// PackCmd will pack any directory into a pbo file.
type PackCmd struct {
	verbose     bool
	deleteAfter bool
	destination string
}

// Name returns the name of this command.
func (*PackCmd) Name() string { return "pack" }

// Synopsis returns the synopsis of this command.
func (*PackCmd) Synopsis() string { return "Packs a directory into a PBO" }

// Usage returns the usage of this command.
func (*PackCmd) Usage() string {
	return `pack [-verbose] [-delete] [-destination <pbo>] <directory>
`
}

// SetFlags registers all flags for this command.
func (u *PackCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&u.verbose, "verbose", false, "be verbose")
	f.BoolVar(&u.deleteAfter, "delete", false, "delete directory after successful packing")
	f.StringVar(&u.destination, "destination", "", "set destination where to save the PBO. Defaults to <directory>.pbo")
}

// Execute parses and executes the input.
func (u *PackCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := f.Arg(0)
	if path == "" {
		fmt.Printf("Please provide the directory to pack.\n\n")
		return subcommands.ExitUsageError
	}

	if u.destination == "" {
		if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
			u.destination = path[:len(path)-1] + ".pbo"
		} else {
			u.destination = path + ".pbo"
		}
	}

	if err := pbo.Pack(path, u.destination, u.verbose); err != nil {
		fmt.Printf("An error occurred while packing PBO: %v\n", err)
		return subcommands.ExitFailure
	}

	if u.deleteAfter {
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Could not remove directory: %v\n", err)
			return subcommands.ExitFailure
		}
		if u.verbose {
			fmt.Printf("%s deleted\n", path)
		}
	}

	return subcommands.ExitSuccess
}
