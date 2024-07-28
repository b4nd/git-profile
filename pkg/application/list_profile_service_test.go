package application

import (
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func generateProfiles(t *testing.T, length uint) []*domain.Profile {
	faker := faker.New()

	profiles := make([]*domain.Profile, length)

	for i := 0; i < int(length); i++ {
		profile, err := domain.NewProfile(
			faker.Internet().User()+"-"+faker.RandomStringWithLength(10),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		assert.NoError(t, err)

		profiles[i] = profile
	}

	return profiles
}

func TestListProfileService_Execute(t *testing.T) {
	t.Run("should return all profiles", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		profiles := generateProfiles(t, 10)

		mockProfileRepository.On("List").Return(profiles, nil)

		listProfileService := NewListProfileService(mockProfileRepository)
		listProfiles, err := listProfileService.Execute()

		assert.NoError(t, err)
		assert.Equal(t, profiles, listProfiles)

		for i := 0; i < len(profiles); i++ {
			assert.Equal(t, profiles[i].Workspace(), listProfiles[i].Workspace())
			assert.Equal(t, profiles[i].Email(), listProfiles[i].Email())
			assert.Equal(t, profiles[i].Name(), listProfiles[i].Name())
		}

		mockProfileRepository.AssertExpectations(t)
	})

	t.Run("should return error when profile repository list fails", func(t *testing.T) {
		mockProfileRepository := &infrastructure.MockProfileRepository{}

		mockProfileRepository.On("List").Return([]*domain.Profile{}, assert.AnError)

		listProfileService := NewListProfileService(mockProfileRepository)
		listProfiles, err := listProfileService.Execute()

		assert.Error(t, err)
		assert.Nil(t, listProfiles)

		mockProfileRepository.AssertExpectations(t)
	})
}
