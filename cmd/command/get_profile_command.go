package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"fmt"
	"os"
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
	reader := bufio.NewReader(os.Stdin)
	var workspaceFlag string

	cmd := &cobra.Command{
		Use:   "get [-w workspace]",
		Short: "Get profile",
		Long:  `get profile with the given workspace, email and name.`,
		Example: `git-profile get
git-profile get work
git-profile get --workspace work 
git-profile get -w work"`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if workspaceFlag == "" && len(args) > 0 {
				workspaceFlag = args[0]
			}

			params := application.GetProfileServiceParams{
				Workspace: workspaceFlag,
			}

			if params.Workspace == "" {
				fmt.Print("Enter workspace: ")
				input, _ := reader.ReadString('\n')
				params.Workspace = strings.TrimSpace(input)
			}

			profile, err := c.getProfileService.Execute(params)

			if err != nil {
				fmt.Println(errorMessages[err], params.Workspace)
				fmt.Printf("\nSuggest to create a new profile with the following command:\n")
				fmt.Printf("  git-profile set %s\n", params.Workspace)
				return
			}

			fmt.Printf("Workspace: %s\n", profile.Workspace().String())
			fmt.Printf("Email: %s\n", profile.Email().String())
			fmt.Printf("Name: %s\n", profile.Name().String())
		},
	}

	cmd.Flags().StringVarP(&workspaceFlag, "workspace", "w", "", "The workspace of the profile")

	rootCmd.AddCommand(cmd)
}
