package command

import (
	"backend/git-profile/pkg/application"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type SetProfileCommand struct {
	createProfileService *application.CreateProfileService
	updateProfileService *application.UpdateProfileService
	getProfileService    *application.GetProfileService
}

func NewSetProfileCommand(
	createProfileService *application.CreateProfileService,
	updateProfileService *application.UpdateProfileService,
	getProfileService *application.GetProfileService,
) *SetProfileCommand {
	return &SetProfileCommand{
		createProfileService,
		updateProfileService,
		getProfileService,
	}
}

type SetProfileCommandParams struct {
	Workspace string
	Email     string
	Name      string
}

func (c *SetProfileCommand) Register(rootCmd *cobra.Command) {
	reader := bufio.NewReader(os.Stdin)
	var workspaceFlag string
	var emailFlag string
	var nameFlag string

	cmd := &cobra.Command{
		Use:   "set [-w workspace] [-e email] [-n name]",
		Short: "Create or update a profile",
		Long: `Create or update a profile with the given workspace, email and name.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `git-profile set
git-profile set work
git-profile set --workspace work --email email@example.com --name "Firstname Lastname"
git-profile set -w work -e email@example.com -n "Firstname Lastname""`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if workspaceFlag == "" && len(args) > 0 {
				workspaceFlag = args[0]
			}

			params := SetProfileCommandParams{
				Workspace: workspaceFlag,
				Email:     emailFlag,
				Name:      nameFlag,
			}

			if workspaceFlag == "" {
				fmt.Print("Enter workspace: ")
				input, _ := reader.ReadString('\n')
				params.Workspace = strings.TrimSpace(input)
			}

			updateProfile := false

			if profile, err := c.getProfileService.Execute(application.GetProfileServiceParams{Workspace: params.Workspace}); err == nil {
				fmt.Printf("Profile \"%s\" already exists, do you want to update it? (y/N): ", profile.Name())
				var answer string
				fmt.Scanln(&answer)

				if answer != "y" {
					return
				}

				updateProfile = true
				params.Email = profile.Email().String()
				params.Name = profile.Name().String()
			}

			if emailFlag == "" {
				fmt.Print("Enter email [" + params.Email + "]: ")
				input, _ := reader.ReadString('\n')
				params.Email = strings.TrimSpace(input)
			}

			if nameFlag == "" {
				fmt.Print("Enter name [" + params.Name + "]: ")
				input, _ := reader.ReadString('\n')
				params.Name = strings.TrimSpace(input)
			}

			if updateProfile {
				profile, err := c.updateProfileService.Execute(application.UpdateProfileServiceParams{
					Workspace: params.Workspace,
					Email:     params.Email,
					Name:      params.Name,
				})

				if err != nil {
					fmt.Printf(errorMessages[err], params.Workspace)
					return
				}

				fmt.Printf("Profile \"%s\" updated successfully", profile.Workspace().String())
				fmt.Printf("\nSuggest to use the updated profile with the following command:\n")
				fmt.Printf("  git-profile use %s\n", params.Workspace)
			} else {
				profile, err := c.createProfileService.Execute(application.CreateProfileServiceParams{
					Workspace: params.Workspace,
					Email:     params.Email,
					Name:      params.Name,
				})

				if err != nil {
					fmt.Printf(errorMessages[err], params.Workspace)
					return
				}

				fmt.Printf("Profile \"%s\" created successfully", profile.Workspace().String())
				fmt.Printf("\nSuggest to use the new profile with the following command:\n")
				fmt.Printf("  git-profile use %s\n", params.Workspace)
			}
		},
	}

	cmd.Flags().StringVarP(&workspaceFlag, "workspace", "w", "", "The workspace of the profile")
	cmd.Flags().StringVarP(&emailFlag, "email", "e", "", "The email of the profile")
	cmd.Flags().StringVarP(&nameFlag, "name", "n", "", "The name of the profile")

	rootCmd.AddCommand(cmd)
}
