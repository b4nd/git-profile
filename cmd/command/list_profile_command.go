package command

import (
	"backend/git-profile/pkg/application"

	"github.com/spf13/cobra"
)

type ListProfileCommand struct {
	listProfileService    *application.ListProfileService
	currentProfileService *application.CurrentProfileService
}

func NewListProfileCommand(
	listProfileService *application.ListProfileService,
	currentProfileService *application.CurrentProfileService,
) *ListProfileCommand {
	return &ListProfileCommand{
		listProfileService,
		currentProfileService,
	}
}

func (c *ListProfileCommand) Register(rootCmd *cobra.Command) {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "list [-v verbose]",
		Short: "Lists all available profiles.",
		Aliases: []string{
			"ls",
		},
		Example: `git-profile list
git-profile list --verbose
git-profile list -v`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Execute(cmd, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show full profile details")

	rootCmd.AddCommand(cmd)
}

func (c *ListProfileCommand) Execute(cmd *cobra.Command, verbose bool) error {
	profiles, err := c.listProfileService.Execute()

	if err != nil {
		cmd.Print(errorMessages[err])
		return nil
	}

	if len(profiles) == 0 {
		cmd.Println("No profiles found")
		return nil
	}

	currentProfile, _ := c.currentProfileService.Execute()
	for _, profile := range profiles {
		isCurrentProfile := currentProfile != nil && currentProfile.Workspace().Equals(profile.Workspace())
		if !verbose {
			if isCurrentProfile {
				cmd.Printf("\033[32m%s\033[0m\n", profile.Workspace().String())
			} else {
				cmd.Printf("%s\n", profile.Workspace().String())
			}
		} else {
			cmd.Printf("Workspace: %s\n", profile.Workspace().String())
			cmd.Printf("Email: %s\n", profile.Email().String())
			cmd.Printf("Name: %s\n", profile.Name().String())
			if isCurrentProfile {
				cmd.Printf("Current: true\n")
			}
			cmd.Println()
		}
	}

	return nil
}
