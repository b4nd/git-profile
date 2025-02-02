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
		Short:   "Displays the current version of the application.",
		Example: `git-profile version`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Execute(cmd)
		},
	}
	rootCmd.AddCommand(cmd)
}

func (c *VersionCommand) Execute(cmd *cobra.Command) error {
	cmd.Printf("Git Version: %s\n", c.GitVersion)
	cmd.Printf("Git Commit: %s\n", c.GitCommit)
	cmd.Printf("Build Date: %s\n", c.BuildDate)
	cmd.Printf("Go Version: %s\n", runtime.Version())
	cmd.Printf("Compiler: %s\n", runtime.Compiler)
	cmd.Printf("Platform: %s\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))

	return nil
}
