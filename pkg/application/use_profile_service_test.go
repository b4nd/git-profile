package application_test

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestUseProfileService_Execute(t *testing.T) {
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

	params := application.UseProfileServiceParams{
		Workspace: workspace.String(),
	}

	assert.NoError(t, err)

	t.Run("should use the profile successfully", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockGitUserRepository.On("Save", scmUser).Return(nil)

		currentProfileService := application.NewUseProfileService(mockProfileRepository, mockGitUserRepository)
		currentProfile, err := currentProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, currentProfile, profile)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockProfileRepository.On("Get", workspace).Return(profile, assert.AnError)

		currentProfileService := application.NewUseProfileService(mockProfileRepository, mockGitUserRepository)
		currentProfile, err := currentProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, currentProfile)
		assert.Equal(t, assert.AnError, err)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when saving the scm user", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockGitUserRepository := &MockUserRepository{}

		mockProfileRepository.On("Get", workspace).Return(profile, nil)
		mockGitUserRepository.On("Save", scmUser).Return(assert.AnError)

		currentProfileService := application.NewUseProfileService(mockProfileRepository, mockGitUserRepository)
		currentProfile, err := currentProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, currentProfile)
		assert.Equal(t, assert.AnError, err)

		mockGitUserRepository.AssertExpectations(t)
		mockProfileRepository.AssertExpectations(t)
	})
}
