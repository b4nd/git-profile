package application

import (
	"github.com/b4nd/git-profile/pkg/domain"
)

type GetProfileService struct {
	profileRepository domain.ProfileRepository
}

type GetProfileServiceParams struct {
	Workspace string
}

func NewGetProfileService(profileRepository domain.ProfileRepository) *GetProfileService {
	return &GetProfileService{profileRepository}
}

func (cp *GetProfileService) Execute(params GetProfileServiceParams) (*domain.Profile, error) {
	workspace, err := domain.NewProfileWorkspace(params.Workspace)
	if err != nil {
		return nil, err
	}

	profile, err := cp.profileRepository.Get(workspace)
	if err != nil {
		return nil, ErrProfileNotExists
	}

	return profile, nil
}
