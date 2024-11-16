package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"fmt"
	"os"
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
	reader := bufio.NewReader(os.Stdin)
	var workspaceFlag string

	cmd := &cobra.Command{
		Use:   "delete [-w workspace]",
		Short: "Delete a profile",
		Long: `Delete a profile with the given workspace.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `git-profile delete
git-profile delete work
git-profile delete --workspace work
git-profile delete -w work"`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if workspaceFlag == "" && len(args) > 0 {
				workspaceFlag = args[0]
			}

			params := DeleteProfileCommandParams{
				Workspace: workspaceFlag,
			}

			if workspaceFlag == "" {
				fmt.Print("Enter workspace: ")
				input, _ := reader.ReadString('\n')
				params.Workspace = strings.TrimSpace(input)
			}

			err := c.deleteProfileService.Execute(application.DeleteProfileServiceParams{
				Workspace: params.Workspace,
			})

			if err != nil {
				fmt.Println(errorMessages[err], params.Workspace)
				return
			}

			fmt.Printf("Profile \"%s\" deleted\n", params.Workspace)
			fmt.Printf("\nSuggest to list all profiles with the following command:\n")
			fmt.Printf("  git-profile list")
		},
	}

	cmd.Flags().StringVarP(&workspaceFlag, "workspace", "w", "", "The workspace of the profile")

	rootCmd.AddCommand(cmd)
}
