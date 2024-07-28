package infrastructure

import (
	"backend/git-profile/pkg/domain"

	"github.com/stretchr/testify/mock"
)

type MockScmUserRepository struct {
	mock.Mock
}

func (m *MockScmUserRepository) Get() (*domain.ScmUser, error) {
	args := m.Called()
	return args.Get(0).(*domain.ScmUser), args.Error(1)
}

func (m *MockScmUserRepository) Save(user *domain.ScmUser) error {
	args := m.Called(user)
	return args.Error(0)
}
