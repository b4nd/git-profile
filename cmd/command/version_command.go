package command

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

type VersionCommand struct {
	Version     string
	GitCommit   string
	BuildDate   string
	ProfilePath string
}

func NewVersionCommand(
	version string,
	gitCommit string,
	buildDate string,
	profilePath string,
) *VersionCommand {
	return &VersionCommand{
		Version:     version,
		GitCommit:   gitCommit,
		BuildDate:   buildDate,
		ProfilePath: profilePath,
	}
}

func (c *VersionCommand) Register(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Displays the current version of the application.",
		Example: `  git-profile version`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Execute(cmd)
		},
	}
	rootCmd.AddCommand(cmd)
}

func (c *VersionCommand) Execute(cmd *cobra.Command) error {
	cmd.Printf("Version: %s\n", c.Version)
	cmd.Printf("Git Commit: %s\n", c.GitCommit)
	cmd.Printf("Build Date: %s\n", c.BuildDate)
	cmd.Printf("Go Version: %s\n", runtime.Version())
	cmd.Printf("Compiler: %s\n", runtime.Compiler)
	cmd.Printf("Platform: %s\n", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	cmd.Printf("Profile Path: %s\n", c.ProfilePath)

	return nil
}
