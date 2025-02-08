package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProfileServiceExecute(t *testing.T) {
	faker := faker.New()

	params := application.UpdateProfileServiceParams{
		Workspace: faker.Internet().User(),
		Email:     faker.Internet().Email(),
		Name:      faker.Person().Name(),
	}

	profile, err := domain.NewProfile(
		params.Workspace,
		params.Email,
		params.Name,
	)

	assert.NoError(t, err)

	t.Run("should update the profile when it exists", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)
		mockProfileRepository.On("Save", profile).Return(nil)

		updateProfileService := application.NewUpdateProfileService(mockProfileRepository)
		newProfile, err := updateProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, profile, newProfile)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)

		updateProfileService := application.NewUpdateProfileService(mockProfileRepository)
		newProfile, err := updateProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, application.ErrProfileNotExists, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile repository save fails", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)
		mockProfileRepository.On("Save", profile).Return(assert.AnError)

		updateProfileService := application.NewUpdateProfileService(mockProfileRepository)
		newProfile, err := updateProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, assert.AnError, err)

		mockProfileRepository.AssertExpectations(t)
	})
}
