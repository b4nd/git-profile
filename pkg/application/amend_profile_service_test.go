package application_test

import (
	"testing"
	"time"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestAmendProfileServiceExecute(t *testing.T) {
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

		headHash := domain.NewScmCommitHashHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(commit, nil).Once()
		mockScmCommitRepository.On("Save", &newAuthor).Return(nil)
		mockScmCommitRepository.On("Get", &headHash).Return(newCommit, nil).Once()

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, ammedCommit, newCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return current commit when author is the same as the profile", func(t *testing.T) {
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
			profile.Name().String(),
			profile.Email().String(),
		)
		assert.NoError(t, err)

		commit := domain.NewScmCommit(
			hash,
			author,
			time.Now(),
			faker.Lorem().Sentence(3),
		)
		assert.NoError(t, err)

		headHash := domain.NewScmCommitHashHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(commit, nil).Once()

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, ammedCommit, ammedCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return error when scm commit does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockScmCommitRepository := &MockCommitRepository{}

		profile, err := domain.NewProfile(
			workspace.String(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		headHash := domain.NewScmCommitHashHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(&domain.ScmCommit{}, assert.AnError)

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, ammedCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return error when scm commit author cannot be amended", func(t *testing.T) {
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

		newAuthor, err := domain.NewScmCommitAuthor(
			profile.Name().String(),
			profile.Email().String(),
		)
		assert.NoError(t, err)

		headHash := domain.NewScmCommitHashHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(commit, nil).Once()
		mockScmCommitRepository.On("Save", &newAuthor).Return(assert.AnError)

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, ammedCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return error when scm commit cannot be retrieved", func(t *testing.T) {
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

		newAuthor, err := domain.NewScmCommitAuthor(
			profile.Name().String(),
			profile.Email().String(),
		)
		assert.NoError(t, err)

		headHash := domain.NewScmCommitHashHead()

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockScmCommitRepository.On("Get", &headHash).Return(commit, nil).Once()
		mockScmCommitRepository.On("Save", &newAuthor).Return(nil)
		mockScmCommitRepository.On("Get", &headHash).Return(&domain.ScmCommit{}, assert.AnError).Once()

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, ammedCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockScmCommitRepository := &MockCommitRepository{}

		mockProfileRepository.On("Get", workspace).Return(&domain.Profile{}, assert.AnError)

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, ammedCommit)
		assert.Equal(t, application.ErrProfileNotExists, err)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile invalid", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockScmCommitRepository := &MockCommitRepository{}

		amendProfileService := application.NewAmendProfileService(mockProfileRepository, mockScmCommitRepository)
		ammedCommit, err := amendProfileService.Execute(application.AmendProfileServiceParams{
			Workspace: "test invalid",
		})

		assert.Error(t, err)
		assert.Nil(t, ammedCommit)

		mockProfileRepository.AssertExpectations(t)
		mockScmCommitRepository.AssertExpectations(t)
	})
}
