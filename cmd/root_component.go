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
	ProfileRepository       domain.ProfileRepository
	ScmUserRepository       domain.ScmUserRepository
	ScmGlobalUserRepository domain.ScmUserRepository
	ScmCommitRepository     domain.ScmCommitRepository

	CreateProfileService        *application.CreateProfileService
	UpdateProfileService        *application.UpdateProfileService
	GetProfileService           *application.GetProfileService
	ListProfileService          *application.ListProfileService
	DeleteProfileService        *application.DeleteProfileService
	SetProfileService           *application.SetProfileService
	SetProfileGlobalService     *application.SetProfileService
	UnsetProfileService         *application.UnsetProfileService
	UnsetProfileGlobalService   *application.UnsetProfileService
	CurrentProfileService       *application.CurrentProfileService
	CurrentProfileGlobalService *application.CurrentProfileService
	AmendProfileService         *application.AmendProfileService

	VersionCommand        *command.VersionCommand
	UpsertProfileCommand  *command.CreateProfileCommand
	GetProfileCommand     *command.GetProfileCommand
	ListProfileCommand    *command.ListProfileCommand
	DeleteProfileCommand  *command.DeleteProfileCommand
	SetProfileCommand     *command.SetProfileCommand
	UnsetProfileCommand   *command.UnsetProfileCommand
	CurrentProfileCommand *command.CurrentProfileCommand
	AmendProfileCommand   *command.AmendProfileCommitCommand
}

type RootComponentOption struct {
	// profile flag is used to set the profile file path (default is $HOME/.gitprofile)
	profile string
	// local flag is used to set the local profile (default is .gitprofile in the current directory)
	local bool
	// workingDir flag is used to set the current working directory (default is the current directory)
	workingDir string
	// userHomeDir is used to set the user home directory (default is the user home directory)
	userHomeDir string
}

func NewRootComponent(option *RootComponentOption) (*RootComponent, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if option != nil && option.userHomeDir != "" {
		userHomeDir = option.userHomeDir
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if option != nil && option.workingDir != "" {
		workingDir = option.workingDir
	}

	// Set default profile file path to $HOME/.gitprofile if not provided by user
	profiles, err := resolveProfileLocations(workingDir, userHomeDir, option)
	if err != nil {
		return nil, err
	}

	// Repositories
	profileRepository, err := infrastructure.NewIniFileProfileRepository(profiles)
	if err != nil {
		return nil, err
	}

	scmUserRepository, err := infrastructure.NewGitUserRepository(path.Join(workingDir, infrastructure.GIT_LOCAL_CONFIG_FILE))
	if err != nil {
		return nil, err
	}

	scmGlobalUserRepository, err := infrastructure.NewGitUserRepository(path.Join(userHomeDir, infrastructure.GIT_GLOBAL_CONFIG_FILE))
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
	setProfileService := application.NewSetProfileService(profileRepository, scmUserRepository)
	setProfileGlobalService := application.NewSetProfileService(profileRepository, scmGlobalUserRepository)
	usetProfileService := application.NewUnsetProfileService(scmUserRepository)
	unsetProfileGlobalService := application.NewUnsetProfileService(scmGlobalUserRepository)
	currentProfileService := application.NewCurrentProfileService(profileRepository, scmUserRepository)
	currentProfileGlobalService := application.NewCurrentProfileService(profileRepository, scmGlobalUserRepository)
	amendProfileService := application.NewAmendProfileService(profileRepository, scmCommitRepository)

	// Command
	versionCommand := command.NewVersionCommand(version, gitCommit, buildDate, profiles[0])
	createProfileCommand := command.NewCreateProfileCommand(createProfileService, updateProfileService, getProfileService)
	getProfileCommand := command.NewGetProfileCommand(getProfileService)
	listProfileCommand := command.NewListProfileCommand(listProfilesService, currentProfileService)
	deleteProfileCommand := command.NewDeleteProfileCommand(getProfileService, deleteProfileService)
	SetProfileCommand := command.NewSetProfileCommand(setProfileService, setProfileGlobalService, getProfileService, listProfilesService)
	unsetProfileCommand := command.NewUnsetProfileCommand(usetProfileService, unsetProfileGlobalService, currentProfileService, currentProfileGlobalService)
	currentProfileCommand := command.NewCurrentProfileCommand(currentProfileService, currentProfileGlobalService)
	amendProfileCommitCommand := command.NewAmendProfileCommitCommnad(currentProfileService, amendProfileService)

	return &RootComponent{
		// Repositories
		ProfileRepository:       profileRepository,
		ScmUserRepository:       scmUserRepository,
		ScmGlobalUserRepository: scmGlobalUserRepository,
		ScmCommitRepository:     scmCommitRepository,
		// Services
		CreateProfileService:        createProfileService,
		GetProfileService:           getProfileService,
		ListProfileService:          listProfilesService,
		DeleteProfileService:        deleteProfileService,
		SetProfileService:           setProfileService,
		SetProfileGlobalService:     setProfileGlobalService,
		UnsetProfileService:         usetProfileService,
		UnsetProfileGlobalService:   unsetProfileGlobalService,
		CurrentProfileService:       currentProfileService,
		CurrentProfileGlobalService: currentProfileGlobalService,
		AmendProfileService:         amendProfileService,
		// Command
		VersionCommand:        versionCommand,
		UpsertProfileCommand:  createProfileCommand,
		GetProfileCommand:     getProfileCommand,
		ListProfileCommand:    listProfileCommand,
		DeleteProfileCommand:  deleteProfileCommand,
		SetProfileCommand:     SetProfileCommand,
		UnsetProfileCommand:   unsetProfileCommand,
		CurrentProfileCommand: currentProfileCommand,
		AmendProfileCommand:   amendProfileCommitCommand,
	}, nil
}

func resolveProfileLocations(workingDir string, userHomeDir string, option *RootComponentOption) ([]string, error) {
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
