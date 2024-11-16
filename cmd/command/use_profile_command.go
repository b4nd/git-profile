package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"fmt"
	"os"
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
	reader := bufio.NewReader(os.Stdin)
	var workspaceFlag string

	cmd := &cobra.Command{
		Use:   "use [-w workspace]",
		Short: "Use a profile",
		Long: `Use a profile with the given workspace.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `git-profile use
git-profile use work
git-profile use --workspace work
git-profile use -w work"`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if workspaceFlag == "" && len(args) > 0 {
				workspaceFlag = args[0]
			}

			params := UseProfileCommandParams{
				Workspace: workspaceFlag,
			}

			if workspaceFlag == "" {
				profiles, err := c.listProfileService.Execute()

				if err != nil {
					fmt.Println(errorMessages[err])
					return
				}

				if len(profiles) == 0 {
					fmt.Println("No profiles found")
					return
				}

				for _, profile := range profiles {
					fmt.Printf("%s\n", profile.Workspace().String())
				}

				fmt.Print("Enter workspace: ")
				input, _ := reader.ReadString('\n')
				params.Workspace = strings.TrimSpace(input)
			}

			profile, err := c.useProfileService.Execute(application.UseProfileServiceParams{
				Workspace: params.Workspace,
			})

			if err != nil {
				fmt.Println(errorMessages[err])
				return
			}

			fmt.Printf("Profile \"%s\" is now in use", profile.Workspace().String())
			fmt.Printf("\nSuggest to use the profile with the following command:\n")
			fmt.Printf("  git-profile list\n")
		},
	}

	cmd.Flags().StringVarP(&workspaceFlag, "workspace", "w", "", "The workspace of the profile")

	rootCmd.AddCommand(cmd)
}
