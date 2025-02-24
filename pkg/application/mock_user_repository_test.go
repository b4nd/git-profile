package application_test

import (
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get() (*domain.ScmUser, error) {
	args := m.Called()
	return args.Get(0).(*domain.ScmUser), args.Error(1)
}

func (m *MockUserRepository) Save(user *domain.ScmUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete() error {
	args := m.Called()
	return args.Error(0)
}
