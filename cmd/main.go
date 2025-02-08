package main

import (
	"os"

	"github.com/spf13/cobra"
)

// These variables are set in build step using
// -ldflags="-X main.version=... -X main.gitCommit=... -X main.buildDate=..."
var (
	version   = "v0.0.0-develop"
	gitCommit = "0000000"
	buildDate = "0000-00-00T00:00:00Z"
)

// Environment Variables
const (
	profileEnvName = "GIT_PROFILE_PATH"
)

// Global Flags
var (
	profileFlag string = ""
	localFlag   bool   = false
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "git-profile [command]",
		Short: "Manage your git profiles",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "git profile",
		},
	}

	// Global Flags
	rootCmd.PersistentFlags().BoolVarP(&localFlag, "local", "l", false, "Set the local profile (default is .gitprofile in the current directory)")
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "file", "f", os.Getenv(profileEnvName), "Set the profile file path (default is $HOME/.gitprofile)")

	// nolint
	rootCmd.ParseFlags(os.Args[1:]) // #nosec G104

	rootComponent, err := NewRootComponent(&RootComponentOption{
		profile: profileFlag,
		local:   localFlag,
	})

	if err != nil {
		panic(err)
	}

	// Register all commands to the root command
	rootComponent.VersionCommand.Register(rootCmd)
	rootComponent.UpsertProfileCommand.Register(rootCmd)
	rootComponent.GetProfileCommand.Register(rootCmd)
	rootComponent.ListProfileCommand.Register(rootCmd)
	rootComponent.DeleteProfileCommand.Register(rootCmd)
	rootComponent.UseProfileCommand.Register(rootCmd)
	rootComponent.CurrentProfileCommand.Register(rootCmd)
	rootComponent.AmendProfileCommand.Register(rootCmd)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
