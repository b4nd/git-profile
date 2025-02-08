package application_test

import (
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/stretchr/testify/mock"
)

type MockCommitRepository struct {
	mock.Mock
}

func (m *MockCommitRepository) Get(hash *domain.ScmCommitHash) (*domain.ScmCommit, error) {
	args := m.Called(hash)
	return args.Get(0).(*domain.ScmCommit), args.Error(1)
}

func (m *MockCommitRepository) AmendAuthor(author *domain.ScmCommitAuthor) error {
	args := m.Called(author)
	return args.Error(0)
}
