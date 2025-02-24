package infrastructure_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/b4nd/git-profile/pkg/domain"
	"github.com/b4nd/git-profile/pkg/infrastructure"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

const (
	GitConfigFile = "/.git/config"
)

func initializateGitRepository(t *testing.T) string {
	path, err := os.MkdirTemp("", "git")
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	cmd := exec.Command("git", "init")
	cmd.Dir = path
	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	return path
}

func TestGitUserRepository(t *testing.T) {
	faker := faker.New()

	t.Run("should return an error when the path is empty", func(t *testing.T) {
		repository, err := infrastructure.NewGitUserRepository("")
		assert.Error(t, err)
		assert.Nil(t, repository)
	})

	t.Run("should return an error when the profile does not exist", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return profile when set new profile", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		user := domain.NewScmUser(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		err = repository.Save(user)
		assert.NoError(t, err)

		currentUser, err := repository.Get()
		assert.NoError(t, err)
		assert.Equal(t, user, currentUser)
	})

	t.Run("should return profile when set new profile and update", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		user1 := domain.NewScmUser(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		err = repository.Save(user1)
		assert.NoError(t, err)

		currentUser1, err := repository.Get()
		assert.NoError(t, err)
		assert.Equal(t, user1, currentUser1)

		user2 := domain.NewScmUser(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		err = repository.Save(user2)
		assert.NoError(t, err)

		currentUser2, err := repository.Get()
		assert.NoError(t, err)
		assert.Equal(t, user2, currentUser2)
	})

	t.Run("should return an error when the path is not a git repository", func(t *testing.T) {
		path, err := os.MkdirTemp("", "git")
		assert.NoError(t, err)
		assert.NotEmpty(t, path)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error when the git repository does not have a config file", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.Remove(path + GitConfigFile)
		assert.NoError(t, err)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error when the git repository does not have a user section", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.WriteFile(path+GitConfigFile, []byte(`
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
`), 0644)

		assert.NoError(t, err)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return profile when set new profile and create file is not exist", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.Remove(path + GitConfigFile)
		assert.NoError(t, err)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		user := domain.NewScmUser(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		err = repository.Save(user)
		assert.NoError(t, err)

		currentUser, err := repository.Get()
		assert.NoError(t, err)
		assert.Equal(t, user, currentUser)
	})

	t.Run("should not return an error when delete profile", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		user := domain.NewScmUser(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)

		err = repository.Save(user)
		assert.NoError(t, err)

		currentUser, err := repository.Get()
		assert.NoError(t, err)
		assert.Equal(t, user, currentUser)

		err = repository.Delete()
		assert.NoError(t, err)

		currentUser, err = repository.Get()
		assert.Error(t, err)
		assert.Nil(t, currentUser)
	})

	t.Run("should not return an error when delete profile and create file is not exist", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.Remove(path + GitConfigFile)
		assert.NoError(t, err)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		err = repository.Delete()
		assert.NoError(t, err)
	})

	t.Run("should return an error when delete profile and user section is not exist", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.WriteFile(path+GitConfigFile, []byte(`
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
`), 0644)

		assert.NoError(t, err)

		repository, err := infrastructure.NewGitUserRepository(path + GitConfigFile)
		assert.NoError(t, err)

		err = repository.Delete()
		assert.Error(t, err)
	})

}
