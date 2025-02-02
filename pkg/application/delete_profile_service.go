package application

import (
	"backend/git-profile/pkg/domain"
)

type DeleteProfileService struct {
	profileRepository domain.ProfileRepository
}

type DeleteProfileServiceParams struct {
	Workspace string
}

func NewDeleteProfileService(profileRepository domain.ProfileRepository) *DeleteProfileService {
	return &DeleteProfileService{profileRepository}
}

func (cp *DeleteProfileService) Execute(params DeleteProfileServiceParams) error {
	workspace, err := domain.NewProfileWorkspace(params.Workspace)
	if err != nil {
		return err
	}

	if _, err := cp.profileRepository.Get(workspace); err != nil {
		return ErrProfileNotExists
	}

	if err = cp.profileRepository.Delete(workspace); err != nil {
		return err
	}

	return nil
}
