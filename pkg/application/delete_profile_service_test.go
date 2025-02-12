package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProfileExecute(t *testing.T) {
	t.Run("should delete the profile when it exists", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		profile := profiles[0]

		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)
		mockProfileRepository.On("Delete", profile.Workspace()).Return(nil)

		deleteProfileService := application.NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(application.DeleteProfileServiceParams{
			Workspace: profile.Workspace().String(),
		})

		assert.NoError(t, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		profile := profiles[0]

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)

		deleteProfileService := application.NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(application.DeleteProfileServiceParams{
			Workspace: profile.Workspace().String(),
		})

		assert.Error(t, err)
		assert.Equal(t, application.ErrProfileNotExists, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the profile is invalid", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		deleteProfileService := application.NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(application.DeleteProfileServiceParams{
			Workspace: "test invalid",
		})

		assert.Error(t, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when deleting the profile fails", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		profile := profiles[0]

		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)
		mockProfileRepository.On("Delete", profile.Workspace()).Return(assert.AnError)

		deleteProfileService := application.NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(application.DeleteProfileServiceParams{
			Workspace: profile.Workspace().String(),
		})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)

		mockProfileRepository.AssertExpectations(t)
	})
}
