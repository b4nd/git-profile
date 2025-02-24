package command

import (
	"bufio"
	"strings"

	"github.com/b4nd/git-profile/pkg/application"

	"github.com/spf13/cobra"
)

type SetProfileCommand struct {
	setProfileService       *application.SetProfileService
	setGlobalProfileService *application.SetProfileService
	getProfileService       *application.GetProfileService
	listProfileService      *application.ListProfileService
}

func NewSetProfileCommand(
	setProfileService *application.SetProfileService,
	setGlobalProfileService *application.SetProfileService,
	getProfileService *application.GetProfileService,
	listProfileService *application.ListProfileService,
) *SetProfileCommand {
	return &SetProfileCommand{
		setProfileService,
		setGlobalProfileService,
		getProfileService,
		listProfileService,
	}
}

type SetProfileCommandParams struct {
	Workspace string
}

func (c *SetProfileCommand) Register(rootCmd *cobra.Command) {
	var workspace string
	var global bool

	cmd := &cobra.Command{
		Use: "set [-w workspace] [--global]",
		Aliases: []string{
			"use",
			"switch",
		},
		Short: "Switches to a specific profile for operations.",
		Long: `Switch to a profile with the given workspace.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `  git profile set
  git profile set work
  git profile set --workspace work
  git profile set -w work`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if workspace == "" && len(args) > 0 {
				workspace = args[0]
			}

			return c.Execute(cmd, workspace, global)
		},
	}

	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "The workspace of the profile")
	cmd.Flags().BoolVarP(&global, "global", "g", false, "Use the profile globally for all repositories (default: false)")

	rootCmd.AddCommand(cmd)
}

func (c *SetProfileCommand) Execute(cmd *cobra.Command, workspace string, global bool) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	params := SetProfileCommandParams{
		Workspace: workspace,
	}

	if workspace == "" {
		profiles, err := c.listProfileService.Execute()

		if err != nil {
			cmd.Print(errorMessages[err])
			return nil
		}

		if len(profiles) == 0 {
			cmd.Println("No profiles found")
			return nil
		}

		for _, profile := range profiles {
			cmd.Printf("%s\n", profile.Workspace().String())
		}

		cmd.Print("Enter workspace: ")
		input, _ := reader.ReadString('\n')
		params.Workspace = strings.TrimSpace(input)
	}

	service := c.setProfileService
	if global {
		service = c.setGlobalProfileService
	}

	profile, err := service.Execute(application.SetProfileServiceParams{
		Workspace: params.Workspace,
	})

	if err != nil {
		cmd.Printf(errorMessages[err], workspace)
		return nil
	}

	cmd.Printf("Profile \"%s\" is now in use\n", profile.Workspace().String())
	cmd.Printf("Workspace: %s\n", profile.Workspace().String())
	cmd.Printf("Email: %s\n", profile.Email().String())
	cmd.Printf("Name: %s\n", profile.Name().String())

	return nil
}
