package application

import (
	"github.com/b4nd/git-profile/pkg/domain"
)

type UpdateProfileService struct {
	profileRepository domain.ProfileRepository
}

type UpdateProfileServiceParams struct {
	Workspace string
	Email     string
	Name      string
}

func NewUpdateProfileService(profileRepository domain.ProfileRepository) *UpdateProfileService {
	return &UpdateProfileService{profileRepository}
}

func (cp *UpdateProfileService) Execute(params UpdateProfileServiceParams) (*domain.Profile, error) {
	profile, err := domain.NewProfile(
		params.Workspace,
		params.Email,
		params.Name,
	)

	if err != nil {
		return nil, err
	}

	if _, err := cp.profileRepository.Get(profile.Workspace()); err != nil {
		return nil, ErrProfileNotExists
	}

	if err = cp.profileRepository.Save(profile); err != nil {
		return nil, err
	}

	return profile, nil
}
