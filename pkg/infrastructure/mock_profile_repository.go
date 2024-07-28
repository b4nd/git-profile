package infrastructure

import (
	"backend/git-profile/pkg/domain"

	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) Get(workspace domain.ProfileWorkspace) (*domain.Profile, error) {
	args := m.Called(workspace)
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) Save(profile *domain.Profile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockProfileRepository) Delete(workspace domain.ProfileWorkspace) error {
	args := m.Called(workspace)
	return args.Error(0)
}

func (m *MockProfileRepository) List() ([]*domain.Profile, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Profile), args.Error(1)
}
