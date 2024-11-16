package main

import (
	"backend/git-profile/cmd/command"
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"os"
)

type RootComponent struct {
	ProfileRepository domain.ProfileRepository
	ScmUserRepository domain.ScmUserRepository

	CreateProfileService  *application.CreateProfileService
	UpdateProfileService  *application.UpdateProfileService
	GetProfileService     *application.GetProfileService
	ListProfileService    *application.ListProfileService
	DeleteProfileService  *application.DeleteProfileService
	UseProfileSercice     *application.UseProfileService
	CurrentProfileService *application.CurrentProfileService

	VersionCommand        *command.VersionCommand
	UpsertProfileCommand  *command.SetProfileCommand
	GetProfileCommand     *command.GetProfileCommand
	ListProfileCommand    *command.ListProfileCommand
	DeleteProfileCommand  *command.DeleteProfileCommand
	UseProfileCommand     *command.UseProfileCommand
	CurrentProfileCommand *command.CurrentProfileCommand
}

type RootComponentOption struct {
	profile string
}

func NewRootComponent(
	option *RootComponentOption,
) (*RootComponent, error) {
	// Set default profile file path to $HOME/.gitprofile if not provided by user
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	defaultProfile := userHomeDir + "/.gitprofile"
	if option != nil && option.profile != "" {
		defaultProfile = option.profile
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Repositories
	profileRepository, err := infrastructure.NewIniFileProfileRepository(defaultProfile)
	if err != nil {
		return nil, err
	}

	scmUserRepository, err := infrastructure.NewIniFileScmUserRepository(workingDir)
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

	// Command
	versionCommand := command.NewVersionCommand(gitVersion, gitCommit, buildDate)
	createProfileCommand := command.NewSetProfileCommand(createProfileService, updateProfileService, getProfileService)
	getProfileCommand := command.NewGetProfileCommand(getProfileService)
	listProfileCommand := command.NewListProfileCommand(listProfilesService, currentProfileService)
	deleteProfileCommand := command.NewDeleteProfileCommand(getProfileService, deleteProfileService)
	useProfileCommand := command.NewUseProfileCommand(useProfileService, getProfileService, listProfilesService)
	currentProfileCommand := command.NewCurrentProfileCommand(currentProfileService)

	return &RootComponent{
		// Repositories
		ProfileRepository: profileRepository,
		ScmUserRepository: scmUserRepository,
		// Services
		CreateProfileService:  createProfileService,
		GetProfileService:     getProfileService,
		ListProfileService:    listProfilesService,
		DeleteProfileService:  deleteProfileService,
		UseProfileSercice:     useProfileService,
		CurrentProfileService: currentProfileService,
		// Command
		VersionCommand:        versionCommand,
		UpsertProfileCommand:  createProfileCommand,
		GetProfileCommand:     getProfileCommand,
		ListProfileCommand:    listProfileCommand,
		DeleteProfileCommand:  deleteProfileCommand,
		UseProfileCommand:     useProfileCommand,
		CurrentProfileCommand: currentProfileCommand,
	}, nil
}
