package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/g0dsCookie/gopbo/pbo"
	"github.com/google/subcommands"
)

type ValidateCmd struct {
}

func (*ValidateCmd) Name() string { return "validate" }

func (*ValidateCmd) Synopsis() string { return "Validates a PBO" }

func (*ValidateCmd) Usage() string { return "validate <pbo>\n" }

func (*ValidateCmd) SetFlags(_ *flag.FlagSet) {}

func (*ValidateCmd) Execute(_ context.Context, flags *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	path := flags.Arg(0)
	if path == "" {
		fmt.Println("Please provide the pbo to validate.\n")
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
