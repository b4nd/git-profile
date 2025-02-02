package infrastructure_test

import (
	"archive/zip"
	"backend/git-profile/pkg/domain"
	"backend/git-profile/pkg/infrastructure"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			err = os.MkdirAll(fpath, os.ModePerm)
			assert.NoError(t, err)
			continue
		}

		err = os.MkdirAll(dest+"/"+f.Name[:len(f.Name)-len(f.FileInfo().Name())], os.ModePerm)
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

func TestGitCommitRepository_NewGitCommitRepository(t *testing.T) {
	t.Run("should return HEAD commit", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash := domain.NewScmCommitHead()

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "518836dd1dcc766d8f5a972583b253db856cc4dd", commit.Hash.Value())
		assert.Equal(t, "Your Name", commit.Author.Name())
		assert.Equal(t, "you@example.com", commit.Author.Email())
		assert.Equal(t, "Second Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 18, 57, 25, 0, time.UTC), commit.Date)
	})

	t.Run("should return commit from 'master' branch", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("master")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "518836dd1dcc766d8f5a972583b253db856cc4dd", commit.Hash.Value())
		assert.Equal(t, "Your Name", commit.Author.Name())
		assert.Equal(t, "you@example.com", commit.Author.Email())
		assert.Equal(t, "Second Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 18, 57, 25, 0, time.UTC), commit.Date)
	})

	t.Run("should return initial commit by full hash", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("fc8d711d866b6fac0e4dce8cbe8209f035cda82d")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "fc8d711d866b6fac0e4dce8cbe8209f035cda82d", commit.Hash.Value())
		assert.Equal(t, "Your Name", commit.Author.Name())
		assert.Equal(t, "you@example.com", commit.Author.Email())
		assert.Equal(t, "Initial Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 18, 56, 36, 0, time.UTC), commit.Date)
	})

	t.Run("should return initial commit by short hash", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("fc8d711d")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "fc8d711d866b6fac0e4dce8cbe8209f035cda82d", commit.Hash.Value())
		assert.Equal(t, "Your Name", commit.Author.Name())
		assert.Equal(t, "you@example.com", commit.Author.Email())
		assert.Equal(t, "Initial Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 18, 56, 36, 0, time.UTC), commit.Date)
	})

	t.Run("should return commit from 'develop' branch", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("develop")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "4094389632c66e23559d46b1899110e9368d79e7", commit.Hash.Value())
		assert.Equal(t, "Dev Name", commit.Author.Name())
		assert.Equal(t, "dev@example.com", commit.Author.Email())
		assert.Equal(t, "Develop Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 19, 5, 59, 0, time.UTC), commit.Date)
	})

	t.Run("should return an error when hash is not found", func(t *testing.T) {
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("0000089632c66e23559d46b1899110e9368d79e7")
		assert.NoError(t, err)

		commit, err := repo.Get(&hash)
		assert.Error(t, err)

		assert.Nil(t, commit)
	})

	t.Run("should return an error when path is empty", func(t *testing.T) {
		path := initializateZipGitRepository(t, "empty_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash, err := domain.NewScmCommitHash("0000089632c66e23559d46b1899110e9368d79e7")
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
		path := initializateZipGitRepository(t, "data_git_repo.zip")

		repo, err := infrastructure.NewGitCommitRepository(path)
		assert.NoError(t, err)
		assert.NotNil(t, repo)

		hash := domain.NewScmCommitHead()

		commit, err := repo.Get(&hash)
		assert.NoError(t, err)
		assert.NotNil(t, commit)

		assert.Equal(t, "518836dd1dcc766d8f5a972583b253db856cc4dd", commit.Hash.Value())
		assert.Equal(t, "Your Name", commit.Author.Name())
		assert.Equal(t, "you@example.com", commit.Author.Email())
		assert.Equal(t, "Second Commit", commit.Message)
		assert.Equal(t, time.Date(2025, 2, 1, 18, 57, 25, 0, time.UTC), commit.Date)

		author, err := domain.NewScmCommitAuthor("New Name", "new@example.com")
		assert.NoError(t, err)

		err = repo.AmendAuthor(&author)
		assert.NoError(t, err)

		newCommit, err := repo.Get(&hash)
		assert.NoError(t, err)

		assert.Equal(t, author.Name(), newCommit.Author.Name())
		assert.Equal(t, author.Email(), newCommit.Author.Email())

		assert.NotEqual(t, newCommit.Hash.Value(), commit.Hash.Value())
		assert.Equal(t, newCommit.Message, commit.Message)
		assert.Equal(t, newCommit.Date, commit.Date)
	})
}
