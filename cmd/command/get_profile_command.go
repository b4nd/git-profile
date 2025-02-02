package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

type GetProfileCommand struct {
	getProfileService *application.GetProfileService
}

func NewGetProfileCommand(
	getProfileService *application.GetProfileService,
) *GetProfileCommand {
	return &GetProfileCommand{
		getProfileService,
	}
}

func (c *GetProfileCommand) Register(rootCmd *cobra.Command) {
	var workspace string

	cmd := &cobra.Command{
		Use:   "get [-w workspace]",
		Short: "Retrieves details of a specific profile.",
		Long:  `get profile with the given workspace, email and name.`,
		Example: `git-profile get
git-profile get work
git-profile get --workspace work 
git-profile get -w work`,
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

func (c *GetProfileCommand) Execute(cmd *cobra.Command, workspace string) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	params := application.GetProfileServiceParams{
		Workspace: workspace,
	}

	if params.Workspace == "" {
		cmd.Print("Enter workspace: ")
		input, _ := reader.ReadString('\n')
		params.Workspace = strings.TrimSpace(input)
	}

	profile, err := c.getProfileService.Execute(params)

	if err != nil {
		cmd.Printf(errorMessages[err], params.Workspace)
		cmd.Printf("\nSuggest to create a new profile with the following command:\n")
		cmd.Printf("  git-profile set %s\n", params.Workspace)
		return nil
	}

	cmd.Printf("Workspace: %s\n", profile.Workspace().String())
	cmd.Printf("Email: %s\n", profile.Email().String())
	cmd.Printf("Name: %s\n", profile.Name().String())

	return nil
}
