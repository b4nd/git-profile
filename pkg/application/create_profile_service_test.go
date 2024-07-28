package application

import (
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCreateProfileService_Execute(t *testing.T) {
	faker := faker.New()

	params := CreateProfileServiceParams{
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

	t.Run("Should create it when the profile does not exist", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)
		mockProfileRepository.On("Save", profile).Return(nil)

		createProfileService := NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, profile, newProfile)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile already exists", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)

		createProfileService := NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, ErrProfileAlreadyExists, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile repository save fails", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)
		mockProfileRepository.On("Save", profile).Return(assert.AnError)

		createProfileService := NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, assert.AnError, err)

		mockProfileRepository.AssertExpectations(t)
	})
}
