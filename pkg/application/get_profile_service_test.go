package application

import (
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile_Execute(t *testing.T) {
	faker := faker.New()

	params := GetProfileServiceParams{
		Workspace: faker.Internet().User(),
	}

	profile, err := domain.NewProfile(
		params.Workspace,
		faker.Internet().Email(),
		faker.Person().Name(),
	)

	assert.NoError(t, err)

	t.Run("should return the profile when it exists", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}
		mockProfileRepository.On("Get", profile.Workspace()).Return(profile, nil)

		getProfileService := NewGetProfileService(mockProfileRepository)
		newProfile, err := getProfileService.Execute(params)

		assert.NoError(t, err)
		assert.NotNil(t, newProfile)
		assert.Equal(t, profile, newProfile)

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		mockProfileRepository.On("Get", profile.Workspace()).Return(&domain.Profile{}, assert.AnError)

		getProfileService := NewGetProfileService(mockProfileRepository)
		_, err := getProfileService.Execute(params)

		assert.Error(t, err)

		mockProfileRepository.AssertExpectations(t)
	})
}
