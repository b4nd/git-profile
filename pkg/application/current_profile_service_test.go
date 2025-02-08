package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCurrentProfileServiceExecute(t *testing.T) {
	faker := faker.New()

	workspace, err := domain.NewProfileWorkspace(faker.Internet().User())
	assert.NoError(t, err)

	scmUser := domain.NewScmUser(
		workspace.String(),
		faker.Internet().Email(),
		faker.Person().Name(),
	)

	profile, err := domain.NewProfile(
		workspace.String(),
		scmUser.Email,
		scmUser.Name,
	)

	assert.NoError(t, err)

	t.Run("should return the current profile successfully", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockGitUserRepository.On("Get").Return(scmUser, nil)
		mockProfileRepository.On("Get", workspace).Return(profile, nil)

		currentProfileService := application.NewCurrentProfileService(mockProfileRepository, mockGitUserRepository)
		currentProfile, err := currentProfileService.Execute()

		assert.NoError(t, err)
		assert.Equal(t, currentProfile, profile)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockGitUserRepository.On("Get").Return(scmUser, nil)
		mockProfileRepository.On("Get", workspace).Return(&domain.Profile{}, assert.AnError)

		currentProfileService := application.NewCurrentProfileService(mockProfileRepository, mockGitUserRepository)
		currentProfile, err := currentProfileService.Execute()

		assert.Error(t, err)
		assert.Nil(t, currentProfile)
		assert.Equal(t, application.ErrProfileNotExists, err)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the scm user does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockGitUserRepository.On("Get").Return(&domain.ScmUser{}, assert.AnError)

		currentProfileService := application.NewCurrentProfileService(mockProfileRepository, mockGitUserRepository)
		profile, err := currentProfileService.Execute()

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Equal(t, assert.AnError, err)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when current profile is not configured", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockGitUserRepository.On("Get").Return(&domain.ScmUser{}, nil)

		currentProfileService := application.NewCurrentProfileService(mockProfileRepository, mockGitUserRepository)
		profile, err := currentProfileService.Execute()

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.Equal(t, application.ErrProfileNotConfigured, err)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})
}
