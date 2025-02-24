package command

import (
	"github.com/b4nd/git-profile/pkg/application"

	"github.com/spf13/cobra"
)

type UnsetProfileCommand struct {
	unsetProfileService         *application.UnsetProfileService
	unsetGlobalProfileService   *application.UnsetProfileService
	currentProfileService       *application.CurrentProfileService
	currentProfileGlobalService *application.CurrentProfileService
}

func NewUnsetProfileCommand(
	unsetProfileService *application.UnsetProfileService,
	unsetGlobalProfileService *application.UnsetProfileService,
	currentProfileService *application.CurrentProfileService,
	currentProfileGlobalService *application.CurrentProfileService,
) *UnsetProfileCommand {
	return &UnsetProfileCommand{
		unsetProfileService,
		unsetGlobalProfileService,
		currentProfileService,
		currentProfileGlobalService,
	}
}

func (c *UnsetProfileCommand) Register(rootCmd *cobra.Command) {
	var global bool

	cmd := &cobra.Command{
		Use:   "unset [--global]",
		Short: "Unset the current profile.",
		Long:  `Unset the current profile.`,
		Aliases: []string{
			"unuse",
		},
		Example: `  git-profile unset`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if global {
				return c.Execute(cmd, c.unsetGlobalProfileService, c.currentProfileGlobalService)
			}
			return c.Execute(cmd, c.unsetProfileService, c.currentProfileService)
		},
	}

	cmd.Flags().BoolVarP(&global, "global", "g", false, "Use the profile globally for all repositories (default: false)")

	rootCmd.AddCommand(cmd)
}

func (c *UnsetProfileCommand) Execute(
	cmd *cobra.Command,
	unsetProfileService *application.UnsetProfileService,
	currentProfileService *application.CurrentProfileService,
) error {
	profile, _ := currentProfileService.Execute()

	err := unsetProfileService.Execute()
	if err != nil {
		return err
	}

	if profile != nil {
		cmd.Printf("Unset profile \"%s\"\n", profile.Workspace())
	} else {
		cmd.Println("Unset profile")
	}

	return nil
}
