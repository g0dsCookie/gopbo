package cmd

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/google/subcommands"
)

const version = "1.0.0"

type VersionCmd struct {
}

func (*VersionCmd) Name() string { return "version" }

func (*VersionCmd) Synopsis() string { return "Prints the tools version" }

func (*VersionCmd) Usage() string {
	return `version
`
}

func (*VersionCmd) SetFlags(_ *flag.FlagSet) {}

func (*VersionCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println("gopbo-"+version, "|", runtime.Version(), "|", runtime.GOOS)
	return subcommands.ExitSuccess
}
