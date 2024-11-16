package infrastructure

import (
	"backend/git-profile/pkg/domain"
	"os"
	"os/exec"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
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

func TestIniFileScmUserRepository(t *testing.T) {
	faker := faker.New()

	t.Run("should return an error when the profile does not exist", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := NewIniFileScmUserRepository(path)
		assert.NoError(t, err)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return profile when set new profile", func(t *testing.T) {
		path := initializateGitRepository(t)

		repository, err := NewIniFileScmUserRepository(path)
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

		repository, err := NewIniFileScmUserRepository(path)
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

		repository, err := NewIniFileScmUserRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error when the git repository does not have a config file", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.Remove(path + "/.git/config")
		assert.NoError(t, err)

		repository, err := NewIniFileScmUserRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return an error when the git repository does not have a user section", func(t *testing.T) {
		path := initializateGitRepository(t)
		err := os.WriteFile(path+"/.git/config", []byte(`
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
`), 0644)

		assert.NoError(t, err)

		repository, err := NewIniFileScmUserRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repository)

		user, err := repository.Get()
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
