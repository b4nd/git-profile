package main

import (
	"os"
	"path"

	"github.com/b4nd/git-profile/cmd/command"
	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"
	"github.com/b4nd/git-profile/pkg/infrastructure"
)

const (
	PROFILE_NAME string = ".gitprofile"
)

type RootComponent struct {
	ProfileRepository   domain.ProfileRepository
	ScmUserRepository   domain.ScmUserRepository
	ScmCommitRepository domain.ScmCommitRepository

	CreateProfileService  *application.CreateProfileService
	UpdateProfileService  *application.UpdateProfileService
	GetProfileService     *application.GetProfileService
	ListProfileService    *application.ListProfileService
	DeleteProfileService  *application.DeleteProfileService
	UseProfileSercice     *application.UseProfileService
	CurrentProfileService *application.CurrentProfileService
	amendProfileService   *application.AmendProfileService

	VersionCommand        *command.VersionCommand
	UpsertProfileCommand  *command.SetProfileCommand
	GetProfileCommand     *command.GetProfileCommand
	ListProfileCommand    *command.ListProfileCommand
	DeleteProfileCommand  *command.DeleteProfileCommand
	UseProfileCommand     *command.UseProfileCommand
	CurrentProfileCommand *command.CurrentProfileCommand
	AmendProfileCommand   *command.AmendProfileCommitCommand
}

type RootComponentOption struct {
	// profile flag is used to set the profile file path (default is $HOME/.gitprofile)
	profile string
	// local flag is used to set the local profile (default is .gitprofile in the current directory)
	local bool
	// pwd flag is used to set the current working directory (default is the current directory)
	pwd string
}

func NewRootComponent(option *RootComponentOption) (*RootComponent, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if option != nil && option.pwd != "" {
		workingDir = option.pwd
	}

	// Set default profile file path to $HOME/.gitprofile if not provided by user
	profiles, err := resolveProfileLocations(workingDir, option)
	if err != nil {
		return nil, err
	}

	// Repositories
	profileRepository, err := infrastructure.NewIniFileProfileRepository(profiles)
	if err != nil {
		return nil, err
	}

	scmUserRepository, err := infrastructure.NewGitUserRepository(workingDir)
	if err != nil {
		return nil, err
	}

	scmCommitRepository, err := infrastructure.NewGitCommitRepository(workingDir)
	if err != nil {
		return nil, err
	}

	// Services
	createProfileService := application.NewCreateProfileService(profileRepository)
	updateProfileService := application.NewUpdateProfileService(profileRepository)
	getProfileService := application.NewGetProfileService(profileRepository)
	listProfilesService := application.NewListProfileService(profileRepository)
	deleteProfileService := application.NewDeleteProfileService(profileRepository)
	useProfileService := application.NewUseProfileService(profileRepository, scmUserRepository)
	currentProfileService := application.NewCurrentProfileService(profileRepository, scmUserRepository)
	amendProfileService := application.NewAmendProfileService(profileRepository, scmCommitRepository)

	// Command
	versionCommand := command.NewVersionCommand(version, gitCommit, buildDate, profiles[0])
	createProfileCommand := command.NewSetProfileCommand(createProfileService, updateProfileService, getProfileService)
	getProfileCommand := command.NewGetProfileCommand(getProfileService)
	listProfileCommand := command.NewListProfileCommand(listProfilesService, currentProfileService)
	deleteProfileCommand := command.NewDeleteProfileCommand(getProfileService, deleteProfileService)
	useProfileCommand := command.NewUseProfileCommand(useProfileService, getProfileService, listProfilesService)
	currentProfileCommand := command.NewCurrentProfileCommand(currentProfileService)
	amendProfileCommitCommand := command.NewAmendProfileCommitCommnad(currentProfileService, amendProfileService)

	return &RootComponent{
		// Repositories
		ProfileRepository:   profileRepository,
		ScmUserRepository:   scmUserRepository,
		ScmCommitRepository: scmCommitRepository,
		// Services
		CreateProfileService:  createProfileService,
		GetProfileService:     getProfileService,
		ListProfileService:    listProfilesService,
		DeleteProfileService:  deleteProfileService,
		UseProfileSercice:     useProfileService,
		CurrentProfileService: currentProfileService,
		amendProfileService:   amendProfileService,
		// Command
		VersionCommand:        versionCommand,
		UpsertProfileCommand:  createProfileCommand,
		GetProfileCommand:     getProfileCommand,
		ListProfileCommand:    listProfileCommand,
		DeleteProfileCommand:  deleteProfileCommand,
		UseProfileCommand:     useProfileCommand,
		CurrentProfileCommand: currentProfileCommand,
		AmendProfileCommand:   amendProfileCommitCommand,
	}, nil
}

func resolveProfileLocations(workingDir string, option *RootComponentOption) ([]string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	defaultProfile := path.Join(userHomeDir, PROFILE_NAME)
	localProfile := path.Join(workingDir, PROFILE_NAME)

	profiles := []string{}
	if option != nil && option.profile != "" {
		defaultProfile = option.profile

		// Check if the profile is a directory or a file
		if info, err := os.Stat(option.profile); err == nil && info.IsDir() {
			defaultProfile = path.Join(option.profile, PROFILE_NAME)
		}
	}

	if option != nil && option.local {
		defaultProfile = localProfile
	}

	profiles = append(profiles, defaultProfile)
	if !option.local && defaultProfile != localProfile {
		profiles = append(profiles, localProfile)
	}

	return profiles, nil
}
