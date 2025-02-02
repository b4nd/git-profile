package application

import (
	"backend/git-profile/pkg/domain"
	"errors"
)

var ErrProfileNotExists = errors.New("profile not exists")

type UseProfileService struct {
	profileRepository domain.ProfileRepository
	scmUserRepository domain.ScmUserRepository
}

type UseProfileServiceParams struct {
	Workspace string
}

func NewUseProfileService(
	profileRepository domain.ProfileRepository,
	scmUserRepository domain.ScmUserRepository,
) *UseProfileService {
	return &UseProfileService{profileRepository, scmUserRepository}
}

func (up *UseProfileService) Execute(params UseProfileServiceParams) (*domain.Profile, error) {
	workspace, err := domain.NewProfileWorkspace(params.Workspace)
	if err != nil {
		return nil, err
	}

	profile, err := up.profileRepository.Get(workspace)
	if err != nil {
		return nil, err
	}

	scmUser := domain.NewScmUser(
		profile.Workspace().String(),
		profile.Email().String(),
		profile.Name().String(),
	)

	err = up.scmUserRepository.Save(scmUser)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
