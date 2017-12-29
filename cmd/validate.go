package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/g0dsCookie/gopbo/pbo"
	"github.com/google/subcommands"
)

// ValidateCmd is used to validate pbo files.
type ValidateCmd struct {
}

// Name returns the name of this command.
func (*ValidateCmd) Name() string { return "validate" }

// Synopsis returns the synopsis of this command.
func (*ValidateCmd) Synopsis() string { return "Validates a PBO" }

// Usage returns the usage of this command.
func (*ValidateCmd) Usage() string {
	return `validate <pbo>
`
}

// SetFlags registers all flags for this command.
func (*ValidateCmd) SetFlags(_ *flag.FlagSet) {}

// Execute parses and executes the input.
func (*ValidateCmd) Execute(_ context.Context, flags *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := flags.Arg(0)
	if path == "" {
		fmt.Printf("Please provide the pbo to validate.\n\n")
		return subcommands.ExitUsageError
	}

	f, err := pbo.Load(path)
	if err != nil {
		fmt.Printf("Could not load PBO: %v\n", err)
		return subcommands.ExitFailure
	}
	f.ToggleCache(false) // Explicitly disable caching

	for _, v := range f.Files {
		if _, err = v.Data(); err != nil {
			fmt.Printf("Could not get data for file %s: %v\n", v.Filename, err)
			return subcommands.ExitFailure
		}
	}

	fmt.Println("PBO seems OK!")

	return subcommands.ExitSuccess
}
