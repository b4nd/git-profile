package application

import "github.com/b4nd/git-profile/pkg/domain"

type ListProfileService struct {
	profileRepository domain.ProfileRepository
}

func NewListProfileService(
	profileRepository domain.ProfileRepository,
) *ListProfileService {
	return &ListProfileService{profileRepository}
}

func (lp *ListProfileService) Execute() ([]*domain.Profile, error) {
	profiles, err := lp.profileRepository.List()
	if err != nil {
		return nil, err
	}

	return profiles, nil
}
