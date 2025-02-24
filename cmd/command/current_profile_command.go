package command

import (
	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/spf13/cobra"
)

type CurrentProfileCommand struct {
	currentProfileService       *application.CurrentProfileService
	currentProfileGlobalService *application.CurrentProfileService
}

func NewCurrentProfileCommand(
	currentProfileService *application.CurrentProfileService,
	currentProfileGlobalService *application.CurrentProfileService,
) *CurrentProfileCommand {
	return &CurrentProfileCommand{
		currentProfileService,
		currentProfileGlobalService,
	}
}

func (c *CurrentProfileCommand) Register(rootCmd *cobra.Command) {
	var verbose bool
	var global bool

	cmd := &cobra.Command{
		Use:     "current [--verbose] [--global]",
		Short:   "Displays the currently active profile",
		Example: `  git profile current`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Execute(cmd, verbose, global)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display the profile in verbose mode")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Display the global profile (default: false)")

	rootCmd.AddCommand(cmd)
}

func (c *CurrentProfileCommand) Execute(cmd *cobra.Command, verbose bool, global bool) error {
	service := c.currentProfileService
	if global {
		service = c.currentProfileGlobalService
	}

	profile, err := service.Execute()
	if err != nil {
		cmd.Println("Profile not found")
		return nil
	}

	if verbose {
		cmd.Printf("Workspace: %s\n", profile.Workspace().String())
		cmd.Printf("Email: %s\n", profile.Email().String())
		cmd.Printf("Name: %s\n", profile.Name().String())
		return nil
	}

	// If the profile is not configured and show the email and name
	if profile.Workspace().String() == domain.NotConfiguredWorkspace {
		cmd.Printf("Email: %s\n", profile.Email().String())
		cmd.Printf("Name: %s\n", profile.Name().String())
		cmd.Printf("\nProfile not configured, suggest to use the new profile with the following command:\n")
		cmd.Printf("  git profile set\n")
		return nil
	}

	cmd.Println(profile.Workspace().String())
	return nil
}
