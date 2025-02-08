package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestCreateProfileServiceExecute(t *testing.T) {
	faker := faker.New()

	params := application.CreateProfileServiceParams{
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

	t.Run("should create it when the profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)
		mockProfileRepository.On("Save", profile).Return(nil)

		createProfileService := application.NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.NoError(t, err)
		assert.Equal(t, profile, newProfile)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile already exists", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)

		createProfileService := application.NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, application.ErrProfileAlreadyExists, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile repository save fails", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)
		mockProfileRepository.On("Save", profile).Return(assert.AnError)

		createProfileService := application.NewCreateProfileService(mockProfileRepository)
		newProfile, err := createProfileService.Execute(params)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
		assert.Equal(t, assert.AnError, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		testParams := application.CreateProfileServiceParams{
			Workspace: params.Workspace,
			Email:     params.Email,
			Name:      params.Name,
		}

		testParams.Email = faker.Internet().User()

		createProfileService := application.NewCreateProfileService(nil)
		newProfile, err := createProfileService.Execute(testParams)

		assert.Error(t, err)
		assert.Nil(t, newProfile)

		testParams.Email = ""

		newProfile, err = createProfileService.Execute(testParams)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
	})

	t.Run("should return error when name is invalid", func(t *testing.T) {
		testParams := application.CreateProfileServiceParams{
			Workspace: params.Workspace,
			Email:     params.Email,
			Name:      params.Name,
		}

		testParams.Name = ""

		createProfileService := application.NewCreateProfileService(nil)
		newProfile, err := createProfileService.Execute(testParams)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
	})

	t.Run("should return error when workspace is invalid", func(t *testing.T) {
		testParams := application.CreateProfileServiceParams{
			Workspace: params.Workspace,
			Email:     params.Email,
			Name:      params.Name,
		}

		testParams.Workspace = ""

		createProfileService := application.NewCreateProfileService(nil)
		newProfile, err := createProfileService.Execute(testParams)

		assert.Error(t, err)
		assert.Nil(t, newProfile)

		testParams.Workspace = faker.Person().Name()

		newProfile, err = createProfileService.Execute(testParams)

		assert.Error(t, err)
		assert.Nil(t, newProfile)
	})
}
