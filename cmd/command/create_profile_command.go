package command

import (
	"bufio"
	"strings"

	"github.com/b4nd/git-profile/pkg/application"

	"github.com/spf13/cobra"
)

type CreateProfileCommand struct {
	createProfileService *application.CreateProfileService
	updateProfileService *application.UpdateProfileService
	getProfileService    *application.GetProfileService
}

func NewCreateProfileCommand(
	createProfileService *application.CreateProfileService,
	updateProfileService *application.UpdateProfileService,
	getProfileService *application.GetProfileService,
) *CreateProfileCommand {
	return &CreateProfileCommand{
		createProfileService,
		updateProfileService,
		getProfileService,
	}
}

type CreateProfileCommandParams struct {
	Workspace string
	Email     string
	Name      string
}

func (c *CreateProfileCommand) Register(rootCmd *cobra.Command) {
	var workspace string
	var email string
	var name string
	var force bool

	cmd := &cobra.Command{
		Use: "add [-w workspace] [-e email] [-n name] [--force]",
		Aliases: []string{
			"create",
		},
		Short: "Add or updates a profile configuration.",
		Long: `Add or update a profile with the given workspace, email and name.
If no arguments are provided, the command will prompt for the missing values.
`,
		Example: `  git profile add
  git profile add work
  git profile add --workspace work --email email@example.com --name "Firstname Lastname"
  git profile add -w work -e email@example.com -n "Firstname Lastname"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if workspace == "" && len(args) > 0 {
				workspace = args[0]
			}

			return c.Execute(cmd, workspace, email, name, force)
		},
	}

	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "The workspace of the profile")
	cmd.Flags().StringVarP(&email, "email", "e", "", "The email of the profile")
	cmd.Flags().StringVarP(&name, "name", "n", "", "The name of the profile")
	cmd.Flags().BoolVar(&force, "force", false, "Force the update of an existing profile")

	rootCmd.AddCommand(cmd)
}

func (c *CreateProfileCommand) Execute(cmd *cobra.Command, workspace string, email string, name string, force bool) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	params := CreateProfileCommandParams{
		Workspace: workspace,
		Email:     email,
		Name:      name,
	}

	if workspace == "" {
		cmd.Print("Enter workspace: ")
		input, _ := reader.ReadString('\n')
		params.Workspace = strings.TrimSpace(input)
	}

	updateProfile, params := c.checkAndUpdateProfile(cmd, reader, params, force)

	if email == "" {
		cmd.Print("Enter email [" + params.Email + "]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			params.Email = input
		}
	}

	if name == "" {
		cmd.Print("Enter name [" + params.Name + "]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			params.Name = input
		}
	}

	if updateProfile {
		profile, err := c.updateProfileService.Execute(application.UpdateProfileServiceParams{
			Workspace: params.Workspace,
			Email:     params.Email,
			Name:      params.Name,
		})

		if err != nil {
			cmd.Printf(errorMessages[err], params.Workspace)
			return nil
		}

		cmd.Printf("Profile \"%s\" updated successfully", profile.Workspace().String())
		cmd.Printf("\nSuggest to set the updated profile with the following command:\n")
		cmd.Printf("  git profile set %s\n", params.Workspace)

		return nil
	}

	profile, err := c.createProfileService.Execute(application.CreateProfileServiceParams{
		Workspace: params.Workspace,
		Email:     params.Email,
		Name:      params.Name,
	})

	if err != nil {
		cmd.Printf(errorMessages[err], params.Workspace)
		return nil
	}

	cmd.Printf("Profile \"%s\" created successfully", profile.Workspace().String())
	cmd.Printf("\nSuggest to set the new profile with the following command:\n")
	cmd.Printf("  git profile set %s\n", params.Workspace)

	return nil
}

func (c *CreateProfileCommand) checkAndUpdateProfile(cmd *cobra.Command, reader *bufio.Reader, params CreateProfileCommandParams, force bool) (bool, CreateProfileCommandParams) {
	profile, err := c.getProfileService.Execute(application.GetProfileServiceParams{Workspace: params.Workspace})
	if err != nil {
		return false, params
	}

	if !force {
		cmd.Printf("Profile \"%s\" already exists, do you want to update it? (y/N): ", profile.Workspace())

		var answer string
		answer, _ = reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if answer != "y" {
			return false, params
		}
	}

	if params.Email == "" {
		params.Email = profile.Email().String()
	}

	if params.Name == "" {
		params.Name = profile.Name().String()
	}

	return true, params
}
