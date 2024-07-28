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

	CreateProfileService *application.CreateProfileService
	UpdateProfileService *application.UpdateProfileService
	GetProfileService    *application.GetProfileService
	ListProfileService   *application.ListProfileService
	DeleteProfileService *application.DeleteProfileService

	VersionCommand       *command.VersionCommand
	UpsertProfileCommand *command.SetProfileCommand
	GetProfileCommand    *command.GetProfileCommand
	ListProfileCommand   *command.ListProfileCommand
	DeleteProfileCommand *command.DeleteProfileCommand
}

func NewRootComponent() (*RootComponent, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Repository

	profileRepository, err := infrastructure.NewIniFileProfileRepository(userHomeDir + "/.gitprofile")
	if err != nil {
		return nil, err
	}

	// Services
	createProfileService := application.NewCreateProfileService(profileRepository)
	updateProfileService := application.NewUpdateProfileService(profileRepository)
	getProfileService := application.NewGetProfileService(profileRepository)
	listProfilesService := application.NewListProfileService(profileRepository)
	deleteProfileService := application.NewDeleteProfileService(profileRepository)

	// Command
	versionCommand := command.NewVersionCommand(gitVersion, gitCommit, buildDate)
	createProfileCommand := command.NewSetProfileCommand(createProfileService, updateProfileService, getProfileService)
	getProfileCommand := command.NewGetProfileCommand(getProfileService)
	listProfileCommand := command.NewListProfileCommand(listProfilesService)
	deleteProfileCommand := command.NewDeleteProfileCommand(getProfileService, deleteProfileService)

	return &RootComponent{
		ProfileRepository: profileRepository,
		// Services
		CreateProfileService: createProfileService,
		GetProfileService:    getProfileService,
		ListProfileService:   listProfilesService,
		DeleteProfileService: deleteProfileService,
		// Command
		VersionCommand:       versionCommand,
		UpsertProfileCommand: createProfileCommand,
		GetProfileCommand:    getProfileCommand,
		ListProfileCommand:   listProfileCommand,
		DeleteProfileCommand: deleteProfileCommand,
	}, nil
}
