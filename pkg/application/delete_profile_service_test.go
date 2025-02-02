package application_test

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProfile_Execute(t *testing.T) {
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

}
