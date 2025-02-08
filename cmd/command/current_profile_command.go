package command

import (
	"github.com/b4nd/git-profile/pkg/application"

	"github.com/spf13/cobra"
)

type CurrentProfileCommand struct {
	currentProfileService *application.CurrentProfileService
}

func NewCurrentProfileCommand(
	currentProfileService *application.CurrentProfileService,
) *CurrentProfileCommand {
	return &CurrentProfileCommand{
		currentProfileService,
	}
}

func (c *CurrentProfileCommand) Register(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "current [-w workspace]",
		Short:   "Displays the currently active profile",
		Example: `git-profile current`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Execute(cmd)
		},
	}

	rootCmd.AddCommand(cmd)
}

func (c *CurrentProfileCommand) Execute(cmd *cobra.Command) error {
	profile, err := c.currentProfileService.Execute()
	if err != nil {
		cmd.Println("Profile not found")
		return nil
	}

	cmd.Println(profile.Workspace().String())

	return nil
}
