package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestGetProfileExecute(t *testing.T) {
	faker := faker.New()

	params := application.GetProfileServiceParams{
		Workspace: faker.Internet().User(),
	}

	profile, err := domain.NewProfile(
		params.Workspace,
		faker.Internet().Email(),
		faker.Person().Name(),
	)

	assert.NoError(t, err)

	t.Run("should return the profile when it exists", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)

		getProfileService := application.NewGetProfileService(mockProfileRepository)
		newProfile, err := getProfileService.Execute(params)

		assert.NoError(t, err)
		assert.NotNil(t, newProfile)
		assert.Equal(t, profile, newProfile)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)

		getProfileService := application.NewGetProfileService(mockProfileRepository)
		_, err := getProfileService.Execute(params)

		assert.Error(t, err)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the profile is invalid", func(t *testing.T) {
		mockProfileRepository := &MockProfileRepository{}

		getProfileService := application.NewGetProfileService(mockProfileRepository)
		profile, err := getProfileService.Execute(application.GetProfileServiceParams{
			Workspace: "test invalid",
		})

		assert.Error(t, err)
		assert.Nil(t, profile)

		mockProfileRepository.AssertExpectations(t)
	})
}
