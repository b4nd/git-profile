package infrastructure_test

import (
	"archive/zip"
	"io"
	"os"
	"testing"
	"time"

	"github.com/b4nd/git-profile/pkg/domain"
	"github.com/b4nd/git-profile/pkg/infrastructure"
	"github.com/jaswdr/faker"

	"github.com/stretchr/testify/assert"
)

const (
	DataGitRepo  = "data_git_repo.zip"
	EmptyGitRepo = "empty_git_repo.zip"
)

// unzipFile unzips a file into a destination folder path using the zip package.
func unzipFile(t *testing.T, src string, dest string) {
	r, err := zip.OpenReader(src)
	assert.NoError(t, err)
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		assert.NoError(t, err)

		fpath := dest + "/" + f.Name

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, 0750)
			assert.NoError(t, err)
			continue
		}

		err = os.MkdirAll(dest+"/"+f.Name[:len(f.Name)-len(f.FileInfo().Name())], 0750)
		assert.NoError(t, err)

		file, err := os.Create(fpath)
		assert.NoError(t, err)

		_, err = io.Copy(file, rc)
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)

		err = rc.Close()
		assert.NoError(t, err)
	}
}

// initializateZipGitRepository creates a temporary directory and unzips a git repository into it.
// It returns the path of the created directory.
func initializateZipGitRepository(t *testing.T, zipName string) string {
	gitPath, err := os.MkdirTemp("", "git")
	assert.NoError(t, err)
	assert.NotEmpty(t, gitPath)

	err = os.Mkdir(gitPath+"/.git", os.ModePerm)
	assert.NoError(t, err)

	unzipFile(t, "../../tests/"+zipName, gitPath+"/.git")
	return gitPath
}

// newGitCommit creates a new domain.ScmCommit instance.
func newGitCommit(t *testing.T, hash string, authorName string, authorEmail string, date time.Time, message string) *domain.ScmCommit {
	author, err := domain.NewScmCommitAuthor(authorName, authorEmail)
	assert.NoError(t, err)

	hashValue, err := domain.NewScmCommitHash(hash)
	assert.NoError(t, err)

	commit := domain.NewScmCommit(hashValue, author, date, message)
	assert.NotNil(t, commit)

	return commit
}

func TestGitCommitRepositoryNewGitCommitRepository(t *testing.T) {
	// The git repository used in the tests is a zip file that contains the following commits:
	// commit 4094389632c66e23559d46b1899110e9368d79e7 (develop)
	// Author: Dev Name <dev@example.com>
	// Date:   Sat Feb 1 19:05:59 2025 +0000
	//
	//     Develop Commit
	//
	// commit 518836dd1dcc766d8f5a972583b253db856cc4dd (HEAD -> master)
	// Author: Your Name <you@example.com>
	// Date:   Sat Feb 1 18:57:25 2025 +0000
	//
	//     Second Commit
	//
	// commit fc8d711d866b6fac0e4dce8cbe8209f035cda82d
	// Author: Your Name <you@example.com>
	// Date:   Sat Feb 1 18:56:36 2025 +0000
	//
	//     Initial Commit

	// The git repository used in the tests is a zip file that contains the following commits:
	var commits = []*domain.ScmCommit{
		newGitCommit(t, "fc8d711d866b6fac0e4dce8cbe8209f035cda82d", "Your Name", "you@example.com", time.Date(2025, 2, 1, 18, 56, 36, 0, time.UTC), "Initial Commit"),
		newGitCommit(t, "518836dd1dcc766d8f5a972583b253db856cc4dd", "Your Name", "you@example.com", time.Date(2025, 2, 1, 18, 57, 25, 0, time.UTC), "Second Commit"),
		newGitCommit(t, "4094389632c66e23559d46b1899110e9368d79e7", "Dev Name", "dev@example.com", time.Date(2025, 2, 1, 19, 5, 59, 0, time.UTC), "Develop Commit"),
	}

	var commitHead = commits[1]
	var commitMaster = commits[1]
	var commitInitial = commits[0]
	var commitDevelop = commits[2]

	faker := faker.New()

	t.Run("should return error when path is empty", func(t *testing.T) {
		repo, err := infrastructure.NewGitCommitRepository("")
		assert.Error(t, err)
		assert.Nil(t, repo)
	})

	t.Run("should return HEAD commit", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash := domain.NewScmCommitHashHead()

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitHead, commit)
	})

	t.Run("should return commit from 'master' branch", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("master")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitMaster, commit)
	})

	t.Run("should return initial commit by full hash", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("fc8d711d866b6fac0e4dce8cbe8209f035cda82d")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitInitial, commit)
	})

	t.Run("should return initial commit by short hash", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("fc8d711d")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitInitial, commit)
	})

	t.Run("should return commit from 'develop' branch", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("develop")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitDevelop, commit)
	})

	t.Run("should return an error when hash is not found", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash(faker.Hash().SHA256())
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.Error(t, err)

		assert.Nil(t, commit)
	})

	t.Run("should return an error when path is empty", func(t *testing.T) {
		path := initializateZipGitRepository(t, EmptyGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash(faker.Hash().SHA256())
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.Error(t, err)

		assert.Nil(t, commit)
	})

	t.Run("should return an error when path is not a git repository", func(t *testing.T) {
		path, err := os.MkdirTemp("", "git")
		assert.NoError(t, err)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("master")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.Error(t, err)
		assert.Nil(t, commit)
	})

	t.Run("should return last commit when upadate author", func(t *testing.T) {
		path := initializateZipGitRepository(t, DataGitRepo)

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash := domain.NewScmCommitHashHead()

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, commitHead, commit)

		email := faker.Internet().Email()
		name := faker.Person().Name()

		author, err := domain.NewScmCommitAuthor(name, email)
		assert.NoError(t, err)

		err = repo.Save(&author)
		assert.NoError(t, err)

		newCommit, err := repo.Get(&hash)
		assert.NoError(t, err)

		assert.Equal(t, author.Name(), newCommit.Author.Name())
		assert.Equal(t, author.Email(), newCommit.Author.Email())

		assert.NotEqual(t, newCommit.Hash.String(), commit.Hash.String())
		assert.Equal(t, newCommit.Message, commit.Message)
		assert.Equal(t, newCommit.Date, commit.Date)
	})
}
