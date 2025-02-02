package application

import (
	"backend/git-profile/pkg/domain"
)

type AmendProfileService struct {
	profileRepository   domain.ProfileRepository
	scmCommitRepository domain.ScmCommitRepository
}

type AmendProfileServiceParams struct {
	Workspace string
}

func NewAmendProfileService(
	profileRepository domain.ProfileRepository,
	scmCommitRepository domain.ScmCommitRepository,
) *AmendProfileService {
	return &AmendProfileService{
		profileRepository,
		scmCommitRepository,
	}
}

func (cp *AmendProfileService) Execute(params AmendProfileServiceParams) (*domain.ScmCommit, error) {
	workspace, err := domain.NewProfileWorkspace(params.Workspace)
	if err != nil {
		return nil, err
	}

	profile, err := cp.profileRepository.Get(workspace)
	if err != nil {
		return nil, ErrProfileAlreadyExists
	}

	scmHash := domain.NewScmCommitHead()
	scmCommit, err := cp.scmCommitRepository.Get(&scmHash)
	if err != nil {
		return nil, err
	}

	// If the author of the commit is the same as the profile, return the commit as the profile
	if scmCommit.Author.Name() == profile.Email().String() && scmCommit.Author.Email() == profile.Email().String() {
		return scmCommit, nil
	}

	// If the author of the commit is different from the profile, amend the commit
	newScmAuthor, err := domain.NewScmCommitAuthor(profile.Name().String(), profile.Email().String())
	if err != nil {
		return nil, err
	}

	err = cp.scmCommitRepository.AmendAuthor(&newScmAuthor)
	if err != nil {
		return nil, err
	}

	// Get the amended commit and return it
	scmCommit, err = cp.scmCommitRepository.Get(&scmHash)
	if err != nil {
		return nil, err
	}

	return scmCommit, nil
}
