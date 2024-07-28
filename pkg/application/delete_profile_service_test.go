package application

import (
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProfile_Execute(t *testing.T) {
	t.Run("should delete the profile when it exists", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		profile := profiles[0]

		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)
		mockProfileRepository.On("Delete", profile.Workspace()).Return(nil)

		deleteProfileService := NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(DeleteProfileServiceParams{
			Workspace: profile.Workspace().String(),
		})

		assert.NoError(t, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		profile := profiles[0]

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)

		deleteProfileService := NewDeleteProfileService(mockProfileRepository)
		err := deleteProfileService.Execute(DeleteProfileServiceParams{
			Workspace: profile.Workspace().String(),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrProfileNotExists, err)

		mockProfileRepository.AssertExpectations(t)
	})

}
