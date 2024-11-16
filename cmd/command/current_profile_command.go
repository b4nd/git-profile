package command

import (
	"backend/git-profile/pkg/application"
	"fmt"

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
		Short:   "Current a profile",
		Example: `git-profile current`,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			profile, err := c.currentProfileService.Execute()
			if err != nil {
				fmt.Println(errorMessages[err])
				return
			}

			fmt.Println(profile.Workspace().String())
		},
	}

	rootCmd.AddCommand(cmd)
}
