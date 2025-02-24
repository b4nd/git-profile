package application

import (
	"errors"

	"github.com/b4nd/git-profile/pkg/domain"
)

var ErrProfileNotExists = errors.New("profile not exists")

type SetProfileService struct {
	profileRepository domain.ProfileRepository
	scmUserRepository domain.ScmUserRepository
}

type SetProfileServiceParams struct {
	Workspace string
}

func NewSetProfileService(
	profileRepository domain.ProfileRepository,
	scmUserRepository domain.ScmUserRepository,
) *SetProfileService {
	return &SetProfileService{profileRepository, scmUserRepository}
}

func (up *SetProfileService) Execute(params SetProfileServiceParams) (*domain.Profile, error) {
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
