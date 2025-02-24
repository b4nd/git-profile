package application

import (
	"github.com/b4nd/git-profile/pkg/domain"
)

type UnsetProfileService struct {
	scmUserRepository domain.ScmUserRepository
}

func NewUnsetProfileService(
	scmUserRepository domain.ScmUserRepository,
) *UnsetProfileService {
	return &UnsetProfileService{scmUserRepository}
}

func (up *UnsetProfileService) Execute() error {
	return up.scmUserRepository.Delete()
}
