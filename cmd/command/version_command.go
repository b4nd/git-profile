package command

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

type VersionCommand struct {
	GitVersion string
	GitCommit  string
	BuildDate  string
}

func NewVersionCommand(
	gitVersion string,
	gitCommit string,
	buildDate string,
) *VersionCommand {
	return &VersionCommand{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
	}
}

func (c *VersionCommand) Register(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version details",
		Example: `git-profile version`,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Git Version: %s\n", c.GitVersion)
			fmt.Printf("Git Commit: %s\n", c.GitCommit)
			fmt.Printf("Build Date: %s\n", c.BuildDate)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("Compiler: %s\n", runtime.Compiler)
			fmt.Printf("Platform: %s\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		},
	}

	rootCmd.AddCommand(cmd)
}
