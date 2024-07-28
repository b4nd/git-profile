package main

import (
	"github.com/spf13/cobra"
)

// These variables are set in build step using
// -ldflags="-X main.gitVersion=... -X main.gitCommit=... -X main.buildDate=..."
var (
	gitVersion = "v0.0.0-develop"
	gitCommit  = "0000000"
	buildDate  = "0000-00-00T00:00:00Z"
)

func main() {
	rootComponent, err := NewRootComponent()
	if err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use:   "git-profile [command]",
		Short: "Manage your git profiles",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "git profile",
		},
	}

	rootComponent.VersionCommand.Register(rootCmd)
	rootComponent.UpsertProfileCommand.Register(rootCmd)
	rootComponent.GetProfileCommand.Register(rootCmd)
	rootComponent.ListProfileCommand.Register(rootCmd)
	rootComponent.DeleteProfileCommand.Register(rootCmd)

	rootCmd.Execute()
}
