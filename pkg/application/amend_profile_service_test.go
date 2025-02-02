package application_test

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestAmendProfileService_Execute(t *testing.T) {
	faker := faker.New()

	params := application.AmendProfileServiceParams{
		Workspace: faker.Internet().User(),
	}

	workspace, err := domain.NewProfileWorkspace(params.Workspace)

	assert.NoError(t, err)

	t.Run("should return new commit with the same author as the profile", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockScmCommitRepository := &MockCommitRepository{}

		profile, err := domain.NewProfile(
			workspace.String(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		hash, err := domain.NewScmCommitHash(faker.Hash().SHA256())
		assert.NoError(t, err)

		author, err := domain.NewScmCommitAuthor(
			faker.Person().Name(),
			faker.Internet().Email(),
		)
		assert.NoError(t, err)

		commit := domain.NewScmCommit(
			hash,
			author,
			time.Now(),
			faker.Lorem().Sentence(3),
		)
		assert.NoError(t, err)

		newHash, err := domain.NewScmCommitHash(faker.Hash().SHA256())
		assert.NoError(t, err)

		newAuthor, err := domain.NewScmCommitAuthor(
			profile.Name().String(),
			profile.Email().String(),
		)
		assert.NoError(t, err)

		newCommit := domain.NewScmCommit(
			newHash,
			newAuthor,
			commit.Date,
			commit.Message,
		)
		assert.NoError(t, err)

		headHash := domain.NewScmCommitHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(commit, nil).Once()
		mockScmCommitRepository.On("AmendAuthor", &newAuthor).Return(nil)
		mockScmCommitRepository.On("Get", &headHash).Return(newCommit, nil).Once()

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, ammedCommit, newCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

}
