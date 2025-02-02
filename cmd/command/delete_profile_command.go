package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

type DeleteProfileCommand struct {
	getProfileService    *application.GetProfileService
	deleteProfileService *application.DeleteProfileService
}

func NewDeleteProfileCommand(
	getProfileService *application.GetProfileService,
	createProfileService *application.DeleteProfileService,
) *DeleteProfileCommand {
	return &DeleteProfileCommand{
		getProfileService,
		createProfileService,
	}
}

type DeleteProfileCommandParams struct {
	Workspace string
}

func (c *DeleteProfileCommand) Register(rootCmd *cobra.Command) {
	var workspace string

	cmd := &cobra.Command{
		Use: "delete [-w workspace]",
		Aliases: []string{
			"del",
		},
		Short: "Deletes a specified profile from the system.",
		Long: `Delete a profile with the given workspace.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `git-profile delete
git-profile delete work
git-profile delete --workspace work
git-profile delete -w work`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if workspace == "" && len(args) > 0 {
				workspace = args[0]
			}

			return c.Execute(cmd, workspace)
		},
	}

	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "The workspace of the profile")

	rootCmd.AddCommand(cmd)
}

func (c *DeleteProfileCommand) Execute(cmd *cobra.Command, workspace string) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	params := DeleteProfileCommandParams{
		Workspace: workspace,
	}

	if workspace == "" {
		cmd.Print("Enter workspace: ")
		input, _ := reader.ReadString('\n')
		params.Workspace = strings.TrimSpace(input)
	}

	err := c.deleteProfileService.Execute(application.DeleteProfileServiceParams{
		Workspace: params.Workspace,
	})

	if err != nil {
		cmd.Printf(errorMessages[err], params.Workspace)
		return nil
	}

	cmd.Printf("Profile \"%s\" deleted\n", params.Workspace)
	cmd.Printf("\nSuggest to list all profiles with the following command:\n")
	cmd.Printf("  git-profile list\n")

	return nil
}
