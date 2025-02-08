package command

import (
	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/spf13/cobra"
)

type AmendProfileCommitCommand struct {
	currentProfileService *application.CurrentProfileService
	amendProfileService   *application.AmendProfileService
}

func NewAmendProfileCommitCommnad(
	currentProfileService *application.CurrentProfileService,
	amendProfileService *application.AmendProfileService,
) *AmendProfileCommitCommand {
	return &AmendProfileCommitCommand{
		currentProfileService: currentProfileService,
		amendProfileService:   amendProfileService,
	}
}

func (c *AmendProfileCommitCommand) Register(rootCmd *cobra.Command) {
	var workspace string

	cmd := &cobra.Command{
		Use:     "amend [-w workspace]",
		Short:   "Amend author of last commit",
		Long:    `Amend author of last commit`,
		Example: `git-profile amend`,
		Args:    cobra.NoArgs,
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

func (c *AmendProfileCommitCommand) Execute(cmd *cobra.Command, workspace string) error {
	if workspace == "" {
		// If no workspace is provided, use the current profile workspace
		profileWorkspace, err := c.currentProfileService.Execute()
		if err != nil {
			cmd.Println("Profile not found")
			return nil
		}

		workspace = profileWorkspace.Workspace().String()
	} else {
		profileWorkspace, err := domain.NewProfileWorkspace(workspace)
		if err != nil {
			cmd.Println("Profile not found")
			return nil
		}

		workspace = profileWorkspace.String()
	}

	commit, err := c.amendProfileService.Execute(application.AmendProfileServiceParams{
		Workspace: workspace,
	})

	if err != nil {
		cmd.Printf(errorMessages[err], workspace)
		return nil
	}

	cmd.Printf("Amended commit author to %s <%s>\n", commit.Author.Name(), commit.Author.Email())
	cmd.Printf("\nSuggest to check the commit with the following command:\n")
	cmd.Printf("  git log -1\n")

	return nil
}
