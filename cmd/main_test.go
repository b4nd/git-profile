//go:build !build

package main

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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
		Use:   "git-profile [command]",
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
	rootComponent.UseProfileCommand.Register(rootCmd)
	rootComponent.CurrentProfileCommand.Register(rootCmd)
	rootComponent.AmendProfileCommand.Register(rootCmd)

	assert.Nil(t, err)

	return rootCmd
}

func TestMain_Command(t *testing.T) {
	faker := faker.New()
	stdout := new(bytes.Buffer)

	t.Cleanup(func() {
		stdout.Reset()
	})

	t.Run("should show the version of the application", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})
		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"version"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Git Version")
		assert.Contains(t, stdout.String(), "Git Commit")
		assert.Contains(t, stdout.String(), "Build Date")
		assert.Contains(t, stdout.String(), "Go Version")
		assert.Contains(t, stdout.String(), "Compiler")
		assert.Contains(t, stdout.String(), "Platform")
		stdout.Reset()
	})

	t.Run("should show the help message", func(t *testing.T) {
		workingDir := initializateGitRepository(t)
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
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
		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)
		rootCmd.SetArgs([]string{"list"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "No profiles found")
		stdout.Reset()
	})

	t.Run("should show the profile created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})
		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Workspace: "+workspace)
		assert.Contains(t, stdout.String(), "Email: "+email)
		assert.Contains(t, stdout.String(), "Name: "+name)
		stdout.Reset()

		// List all profiles verbose
		rootCmd.SetArgs([]string{"list", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Workspace: "+workspace)
		assert.Contains(t, stdout.String(), "Email: "+email)
		assert.Contains(t, stdout.String(), "Name: "+name)
		stdout.Reset()
	})

	t.Run("should show correct deletion of the created profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})
		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Delete the current profile
		rootCmd.SetArgs([]string{"delete", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Profile \""+workspace+"\" deleted")
		stdout.Reset()

		// List all profiles
		rootCmd.SetArgs([]string{"list"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "No profiles found")
		stdout.Reset()
	})

	t.Run("should show the error when the profile does not exist", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)
		workspace := faker.Internet().User()

		// Delete the current profile
		rootCmd.SetArgs([]string{"delete", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Profile \""+workspace+"\" does not exist")
		stdout.Reset()
	})

	t.Run("should show current profile not found", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)
		workspace := faker.Internet().User()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Profile \""+workspace+"\" does not exist")
		stdout.Reset()
	})

	t.Run("should show the profile created and the current profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"use", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "is now in use")
		assert.Contains(t, stdout.String(), "Email: "+email)
		assert.Contains(t, stdout.String(), "Name: "+name)
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"current"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), workspace)
		stdout.Reset()
	})

	t.Run("should show list of profiles created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
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
			rootCmd.SetArgs([]string{"set", "-w", profile.Workspace, "-n", profile.Name, "-e", profile.Email})
			err := rootCmd.Execute()

			assert.Nil(t, err)
			assert.Contains(t, stdout.String(), "created successfully")
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
		rootCmd.SetArgs([]string{"use", "-w", profiels[0].Workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "is now in use")
		stdout.Reset()

		// List all profiles verbose
		rootCmd.SetArgs([]string{"list", "-v"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		for _, profile := range profiels {
			assert.Contains(t, stdout.String(), "Workspace: "+profile.Workspace)
			assert.Contains(t, stdout.String(), "Email: "+profile.Email)
			assert.Contains(t, stdout.String(), "Name: "+profile.Name)
			if profile.Workspace == profiels[0].Workspace {
				assert.Contains(t, stdout.String(), "Current: true")
			}
		}
		stdout.Reset()
	})

	t.Run("should update the profile created", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()
		email := faker.Internet().Email()
		name := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", workspace, "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Workspace: "+workspace)
		assert.Contains(t, stdout.String(), "Email: "+email)
		assert.Contains(t, stdout.String(), "Name: "+name)
		stdout.Reset()

		// Update the current profile
		newEmail := faker.Internet().Email()
		newName := faker.Person().Name()
		rootCmd.SetArgs([]string{"set", "-w", workspace, "-n", newName, "-e", newEmail, "--force"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "updated successfully")
		stdout.Reset()

		// Get the current profile
		rootCmd.SetArgs([]string{"get", "-w", workspace})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Workspace: "+workspace)
		assert.Contains(t, stdout.String(), "Email: "+newEmail)
		assert.Contains(t, stdout.String(), "Name: "+newName)
		stdout.Reset()
	})

	t.Run("should show the error when set current profile does not exist", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		workspace := faker.Internet().User()

		// Set the current profile
		rootCmd.SetArgs([]string{"use", "-w", workspace})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "Profile \""+workspace+"\" does not exist")
		stdout.Reset()
	})

	t.Run("should amend last commit with the profile", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		email := faker.Internet().Email()
		name := faker.Person().Name()

		emptyCommit(t, workingDir, "Initial commit", name, email)
		commit := lastCommit(t, workingDir)

		assert.Equal(t, name+","+email+"\n", commit)

		newEmail := faker.Internet().Email()
		newName := faker.Person().Name()

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", "test", "-n", newName, "-e", newEmail})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend", "-w", "test"})

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

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		email := faker.Internet().Email()
		name := faker.Person().Name()

		emptyCommit(t, workingDir, "Initial commit", name, email)
		commit := lastCommit(t, workingDir)

		assert.Equal(t, name+","+email+"\n", commit)

		// Create a new profile
		rootCmd.SetArgs([]string{"set", "-w", "test", "-n", name, "-e", email})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "created successfully")
		stdout.Reset()

		// Set the current profile
		rootCmd.SetArgs([]string{"use", "-w", "test"})
		err = rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "is now in use")
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

	t.Run("should show the error when the profile does not exist", func(t *testing.T) {
		workingDir := initializateGitRepository(t)

		rootCmd := initializateRootContainer(t, &RootComponentOption{
			profile: t.TempDir(),
			local:   false,
			pwd:     workingDir,
		})

		rootCmd.SetOutput(stdout)

		// Amend the last commit
		rootCmd.SetArgs([]string{"amend", "-w", "test"})
		err := rootCmd.Execute()

		assert.Nil(t, err)
		assert.Contains(t, stdout.String(), "profile already exists")
		stdout.Reset()
	})
}
