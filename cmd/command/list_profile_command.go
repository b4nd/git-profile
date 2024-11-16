package command

import (
	"backend/git-profile/pkg/application"
	"fmt"

	"github.com/spf13/cobra"
)

type ListProfileCommand struct {
	listProfileService    *application.ListProfileService
	currentProfileService *application.CurrentProfileService
}

func NewListProfileCommand(
	listProfileService *application.ListProfileService,
	currentProfileService *application.CurrentProfileService,
) *ListProfileCommand {
	return &ListProfileCommand{
		listProfileService,
		currentProfileService,
	}
}

func (c *ListProfileCommand) Register(rootCmd *cobra.Command) {
	var verboseFlag bool

	cmd := &cobra.Command{
		Use:   "list [-v verbose]",
		Short: "List profiles",
		Example: `git-profile list
git-profile list --verbose
git-profile list -v`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := c.listProfileService.Execute()

			if err != nil {
				fmt.Println(errorMessages[err])
				return
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found")
				return
			}

			currentProfile, err := c.currentProfileService.Execute()
			if err != nil {
				fmt.Println(errorMessages[err])
				return
			}

			for _, profile := range profiles {
				isCurrentProfile := currentProfile != nil && currentProfile.Workspace().String() == profile.Workspace().String()

				if !verboseFlag {
					if isCurrentProfile {
						fmt.Printf("*%s\n", profile.Workspace().String())
					} else {
						fmt.Printf("%s\n", profile.Workspace().String())
					}
				} else {
					fmt.Printf("Workspace: %s\n", profile.Workspace().String())
					fmt.Printf("Email: %s\n", profile.Email().String())
					fmt.Printf("Name: %s\n", profile.Name().String())
					if isCurrentProfile {
						fmt.Printf("Current: true\n")
					}
					fmt.Println()
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Show full profile details")

	rootCmd.AddCommand(cmd)
}
