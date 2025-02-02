package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

type UseProfileCommand struct {
	useProfileService  *application.UseProfileService
	getProfileService  *application.GetProfileService
	listProfileService *application.ListProfileService
}

func NewUseProfileCommand(
	useProfileService *application.UseProfileService,
	getProfileService *application.GetProfileService,
	listProfileService *application.ListProfileService,
) *UseProfileCommand {
	return &UseProfileCommand{
		useProfileService,
		getProfileService,
		listProfileService,
	}
}

type UseProfileCommandParams struct {
	Workspace string
}

func (c *UseProfileCommand) Register(rootCmd *cobra.Command) {
	var workspace string

	cmd := &cobra.Command{
		Use: "use [-w workspace]",
		Aliases: []string{
			"switch",
		},
		Short: "Switches to a specific profile for operations.",
		Long: `Switch to a profile with the given workspace.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `git-profile use
git-profile use work
git-profile use --workspace work
git-profile use -w work`,
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

func (c *UseProfileCommand) Execute(cmd *cobra.Command, workspace string) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	params := UseProfileCommandParams{
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

	profile, err := c.useProfileService.Execute(application.UseProfileServiceParams{
		Workspace: params.Workspace,
	})

	if err != nil {
		cmd.Printf(errorMessages[err], workspace)
		return nil
	}

	cmd.Printf("Profile \"%s\" is now in use\n", profile.Workspace().String())
	cmd.Printf("Email: %s\n", profile.Email().String())
	cmd.Printf("Name: %s\n", profile.Name().String())
	cmd.Printf("\nSuggest to use the profile with the following command:\n")
	cmd.Printf("  git-profile list\n")

	return nil
}
