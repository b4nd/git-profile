package application_test

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile_Execute(t *testing.T) {
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
}
