package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/b4nd/git-profile/pkg/domain"
	"github.com/jaswdr/faker"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	MsgProfileCreatedSuccessfully = "Profile \"%s\" created successfully"
	MsgProfileUpdatedSuccessfully = "Profile \"%s\" updated successfully"
	MsgProfileDeleted             = "Profile \"%s\" deleted"
	MsgProfileUnset               = "Unset profile \"%s\""
	ErrProfileInUse               = "Profile \"%s\" is now in use"
	ErrProfileAlreadyExists       = "Profile \"%s\" already exists"
	ErrProfileNotExist            = "Profile \"%s\" does not exist"
	ErrProfilesNotFound           = "No profiles found"
	ErrProfileNotFound            = "Profile not found"
)

func initializateGitRepository(t *testing.T) string {
	path := t.TempDir()
	assert.NotEmpty(t, path)

	cmd := exec.Command("git", "init")
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	return path
}

func configureGit(t *testing.T, path string, name string, email string, env string) {
	cmd := exec.Command("git", "config", "--"+env, "user.name", name)
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	cmd = exec.Command("git", "config", "--"+env, "user.email", email)
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)
}

func emptyCommit(t *testing.T, path string, message string, author string, email string) {
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", message, "--author", author+" <"+email+">")
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	assert.NoError(t, err)
}

func lastCommit(t *testing.T, path string) string {
	cmd := exec.Command("git", "log", "-1", "--format=%an,%ae")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	return string(output)
}

func initializateRootContainer(t *testing.T, option *RootComponentOption) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "git profile [command]",
		Short: "Manage your git profiles",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "git profile",
		},
	}

	rootComponent, err := NewRootComponent(option)
	rootComponent.VersionCommand.Register(rootCmd)
	rootComponent.UpsertProfileCommand.Register(rootCmd)
	rootComponent.GetProfileCommand.Register(rootCmd)
	rootComponent.ListProfileCommand.Register(rootCmd)
	rootComponent.DeleteProfileCommand.Register(rootCmd)
	rootComponent.SetProfileCommand.Register(rootCmd)
	rootComponent.UnsetProfileCommand.Register(rootCmd)
	rootComponent.CurrentProfileCommand.Register(rootCmd)
	rootComponent.AmendProfileCommand.Register(rootCmd)

	assert.Nil(t, err)

	return rootCmd
}

func TestMainCommand(t *testing.T) {
	faker := faker.New()
	stdout := new(bytes.Buffer)

	t.Cleanup(func() {
		stdout.Reset()
	})

	t.Run("should show the version of the application", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"version"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Version")
		assert.Contains(t, stdout.String(), "Git Commit")
		assert.Contains(t, stdout.String(), "Build Date")
		assert.Contains(t, stdout.String(), "Go Version")
		assert.Contains(t, stdout.String(), "Compiler")
		assert.Contains(t, stdout.String(), "Platform")
		stdout.Reset()
	})

	t.Run("should show the help message", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"--help"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Manage your git profiles")
		stdout.Reset()
	})

	t.Run("should show an error when there are no profiles", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"list"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfilesNotFound)
		stdout.Reset()
	})

	t.Run("should show the profile created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()

		// List all profiles verbose
		rootCmd.SetArgs([]string{"list", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()
	})

	t.Run("should show correct deletion of the created profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Delete the current profile
		rootCmd.SetArgs([]string{"delete", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileDeleted, workspace))
		stdout.Reset()

		// List all profiles
		rootCmd.SetArgs([]string{"list"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfilesNotFound)
		stdout.Reset()
	})

	t.Run("should show the error when the profile does not exist", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)
		workspace := faker.Internet().User()

		// Delete the current profile
		rootCmd.SetArgs([]string{"delete", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileNotExist, workspace))
		stdout.Reset()
	})

	t.Run("should show current profile not found", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)
		workspace := faker.Internet().User()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileNotExist, workspace))
		stdout.Reset()
	})

	t.Run("should show the profile created and the current profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"set", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, workspace))
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()
	})

	t.Run("should show the error when non selected current profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		// Get the current profile
		rootCmd.SetArgs([]string{"current"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should show name and email profile not configured", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		name, email := faker.Person().Name(), faker.Internet().Email()
		configureGit(t, workingDir, name, email, "local")

		// Get the current profile
		rootCmd.SetArgs([]string{"current"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), name)
		assert.Contains(t, stdout.String(), email)
		assert.NotContains(t, stdout.String(), "Workspace")
		stdout.Reset()

		// Get the current profile verbose
		rootCmd.SetArgs([]string{"current", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), domain.NotConfiguredWorkspace)
		assert.Contains(t, stdout.String(), name)
		assert.Contains(t, stdout.String(), email)
		stdout.Reset()
	})

	t.Run("should show current global profile not found", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       true,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "--global"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "-v", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should show the profile created and delete the current global profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       true,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"set", workspace, "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, workspace))
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "-v", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()

		// unset the current profile
		rootCmd.SetArgs([]string{"unset", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileUnset, workspace))
		stdout.Reset()

		// List all profiles
		rootCmd.SetArgs([]string{"current", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should delete the current global profile not found", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       true,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		email := faker.Internet().Email()
		name := faker.Person().Name()

		configureGit(t, workingDir, name, email, "global")

		rootCmd.SetOutput(stdout)

		rootCmd.SetArgs([]string{"unset", "--global"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Unset profile")
		stdout.Reset()

		rootCmd.SetArgs([]string{"current", "--global"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should show list of profiles created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		profiels := []struct {
			Workspace string
			Email     string
			Name      string
		}{
			{
				Workspace: faker.Internet().User(),
				Email:     faker.Internet().Email(),
				Name:      faker.Person().Name(),
			},
			{
				Workspace: faker.Internet().User(),
				Email:     faker.Internet().Email(),
				Name:      faker.Person().Name(),
			},
			{
				Workspace: faker.Internet().User(),
				Email:     faker.Internet().Email(),
				Name:      faker.Person().Name(),
			},
		}

		for _, profile := range profiels {
			rootCmd.SetArgs([]string{"add", "-w", profile.Workspace, "-n", profile.Name, "-e", profile.Email})
			err := rootCmd.Execute()

			assert.Nil(t, err)
			assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, profile.Workspace))
			stdout.Reset()
		}

		// List all profiles
		rootCmd.SetArgs([]string{"list"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		for _, profile := range profiels {
			assert.Contains(t, stdout.String(), profile.Workspace)
		}
		stdout.Reset()

		// Set the current profile
		rootCmd.SetArgs([]string{"set", "-w", profiels[0].Workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, profiels[0].Workspace))
		stdout.Reset()

		// List all profiles verbose
		rootCmd.SetArgs([]string{"list"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		for _, profile := range profiels {
			assert.Contains(t, stdout.String(), profile.Workspace)
		}

		// List all profiles verbose
		rootCmd.SetArgs([]string{"list", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		for _, profile := range profiels {
			assert.Contains(t, stdout.String(), profile.Workspace)
			assert.Contains(t, stdout.String(), profile.Email)
			assert.Contains(t, stdout.String(), profile.Name)
			if profile.Workspace == profiels[0].Workspace {
				assert.Contains(t, stdout.String(), "Current: true")
			}
		}
		stdout.Reset()
	})

	t.Run("should update the profile created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), email)
		assert.Contains(t, stdout.String(), name)
		stdout.Reset()

		// Update the current profile
		newEmail := faker.Internet().Email()
		newName := faker.Person().Name()
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", newName, "-e", newEmail, "--force"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileUpdatedSuccessfully, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		assert.Contains(t, stdout.String(), newEmail)
		assert.Contains(t, stdout.String(), newName)
		stdout.Reset()
	})

	t.Run("should show the error when set current profile does not exist", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()

		// Set the current profile
		rootCmd.SetArgs([]string{"set", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileNotExist, workspace))
		stdout.Reset()
	})

	t.Run("should amend last commit with the profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		emptyCommit(t, workingDir, "Initial commit", name, email)
		commit := lastCommit(t, workingDir)

		assert.Equal(t, name+","+email+"\n", commit)

		newEmail := faker.Internet().Email()
		newName := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", newName, "-e", newEmail})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend", workspace})

		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Amended commit author to "+newName+" <"+newEmail+">")
		stdout.Reset()

		commit = lastCommit(t, workingDir)
		assert.Equal(t, newName+","+newEmail+"\n", commit)
		stdout.Reset()
	})

	t.Run("should amend last commit with the current profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		emptyCommit(t, workingDir, "Initial commit", name, email)
		commit := lastCommit(t, workingDir)

		assert.Equal(t, name+","+email+"\n", commit)

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Set the current profile
		rootCmd.SetArgs([]string{"set", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, workspace))
		stdout.Reset()

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend"})

		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Amended commit author to "+name+" <"+email+">")
		stdout.Reset()

		commit = lastCommit(t, workingDir)
		assert.Equal(t, name+","+email+"\n", commit)
		stdout.Reset()
	})

	t.Run("should show the error when the profile does not exist in amend command", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend", "-w", faker.Internet().User()})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "does not exist")
		stdout.Reset()
	})

	t.Run("should show the error when the profile empty in amend command", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should show the error when the profile invalid in amend command", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend", "-w", "test test"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	t.Run("should remove the current profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"add", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		// Set the current profile
		rootCmd.SetArgs([]string{"set", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, workspace))
		stdout.Reset()

		// Unset the current profile
		rootCmd.SetArgs([]string{"unset"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileUnset, workspace))
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)

		assert.Contains(t, stdout.String(), ErrProfileNotFound)
		stdout.Reset()
	})

	// Test Interactive Mode

	t.Run("should create a new profile in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()
	})

	t.Run("should update a profile in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\ny\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileUpdatedSuccessfully, workspace))
		stdout.Reset()
	})

	t.Run("should show the error when the profile already exists in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\nN\n"))
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileAlreadyExists, workspace))
		stdout.Reset()
	})

	t.Run("should use a profile in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		rootCmd.SetArgs([]string{"set"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n"))
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileInUse, workspace))
		stdout.Reset()
	})

	t.Run("should show the error when the profile does not exist in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"set"})
		rootCmd.SetIn(bytes.NewBufferString("test\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), ErrProfilesNotFound)
		stdout.Reset()
	})

	t.Run("should delete a profile in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		rootCmd.SetArgs([]string{"delete"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n"))
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileDeleted, workspace))
		stdout.Reset()
	})

	t.Run("should show the error when the profile does not exist in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})

		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"delete"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(ErrProfileNotExist, workspace))
		stdout.Reset()
	})

	t.Run("should show name profile in interactive mode", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		userHomeDir := t.TempDir()
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile:     t.TempDir(),
			local:       false,
			workingDir:  workingDir,
			userHomeDir: userHomeDir,
		})
		workspace := faker.Internet().User()

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"add"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n" + faker.Internet().Email() + "\n" + faker.Person().Name() + "\n"))
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(MsgProfileCreatedSuccessfully, workspace))
		stdout.Reset()

		rootCmd.SetArgs([]string{"get"})
		rootCmd.SetIn(bytes.NewBufferString(workspace + "\n"))
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Workspace: "+workspace)
		stdout.Reset()
	})
}
