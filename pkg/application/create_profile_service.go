package application

import (
	"errors"

	"github.com/b4nd/git-profile/pkg/domain"
)

var ErrProfileAlreadyExists = errors.New("profile already exists")

type CreateProfileService struct {
	profileRepository domain.ProfileRepository
}

type CreateProfileServiceParams struct {
	Workspace string
	Email     string
	Name      string
}

func NewCreateProfileService(profileRepository domain.ProfileRepository) *CreateProfileService {
	return &CreateProfileService{profileRepository}
}

func (cp *CreateProfileService) Execute(params CreateProfileServiceParams) (*domain.Profile, error) {
	profile, err := domain.NewProfile(
		params.Workspace,
		params.Email,
		params.Name,
	)

	if err != nil {
		return nil, err
	}

	if _, err := cp.profileRepository.Get(profile.Workspace()); err == nil {
		return nil, ErrProfileAlreadyExists
	}

	if err = cp.profileRepository.Save(profile); err != nil {
		return nil, err
	}

	return profile, nil
}
