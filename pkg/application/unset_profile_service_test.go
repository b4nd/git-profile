package application_test

import (
	"testing"

	"github.com/b4nd/git-profile/pkg/application"

	"github.com/stretchr/testify/assert"
)

func TestUnsetProfileServiceExecute(t *testing.T) {
	t.Run("should update the profile when it exists", func(t *testing.T) {
		mockUserRepository := &MockUserRepository{}
		mockUserRepository.On("Delete").Return(nil)

		UnsetProfileService := application.NewUnsetProfileService(mockUserRepository)
		err := UnsetProfileService.Execute()

		assert.NoError(t, err)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("should return an error when the profile does not exist", func(t *testing.T) {
		mockUserRepository := &MockUserRepository{}
		mockUserRepository.On("Delete").Return(assert.AnError)

		UnsetProfileService := application.NewUnsetProfileService(mockUserRepository)
		err := UnsetProfileService.Execute()

		assert.Error(t, err)

		mockUserRepository.AssertExpectations(t)
	})
}
