package application

import (
	"errors"

	"github.com/b4nd/git-profile/pkg/domain"
)

var ErrProfileNotConfigured = errors.New("profile not configured")

type CurrentProfileService struct {
	profileRepository domain.ProfileRepository
	scmUserRepository domain.ScmUserRepository
}

func NewCurrentProfileService(
	profileRepository domain.ProfileRepository,
	scmUserRepository domain.ScmUserRepository,
) *CurrentProfileService {
	return &CurrentProfileService{profileRepository, scmUserRepository}
}

func (cp *CurrentProfileService) Execute() (*domain.Profile, error) {
	scmUser, err := cp.scmUserRepository.Get()
	if err != nil {
		return nil, err
	}

	if scmUser == nil || scmUser.Workespace == "" {
		return nil, ErrProfileNotConfigured
	}

	workspace, err := domain.NewProfileWorkspace(scmUser.Workespace)
	if err != nil {
		return nil, err
	}

	profile, err := cp.profileRepository.Get(workspace)
	if err != nil {
		return nil, ErrProfileNotExists
	}

	return profile, nil
}
