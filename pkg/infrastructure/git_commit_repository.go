package infrastructure

import (
	"backend/git-profile/pkg/domain"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type GitCommitRepository struct {
	path string
}

func NewGitCommitRepository(path string) (*GitCommitRepository, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	return &GitCommitRepository{path}, nil
}

func (r *GitCommitRepository) Get(hash *domain.ScmCommitHash) (*domain.ScmCommit, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%ai,%H,%an,%ae,%s", hash.Value())
	cmd.Dir = r.path

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(output), ",", 5)
	if len(parts) != 5 {
		return nil, domain.ErrScmCommitNotFound
	}

	scmHash, err := domain.NewScmCommitHash(parts[1])
	if err != nil {
		return nil, err
	}

	scmAuthor, err := domain.NewScmCommitAuthor(parts[2], parts[3])
	if err != nil {
		return nil, err
	}

	scmDate, err := time.Parse("2006-01-02 15:04:05 -0700", parts[0])
	if err != nil {
		return nil, err
	}

	scmMessage := strings.TrimSpace(parts[4])
	scmMessage = strings.ReplaceAll(scmMessage, "\n", " ")

	return domain.NewScmCommit(scmHash, scmAuthor, scmDate.UTC(), scmMessage), nil
}

func (r *GitCommitRepository) AmendAuthor(author *domain.ScmCommitAuthor) error {
	cmd := exec.Command("git", "commit", "--amend", "--author=\""+author.String()+"\"", "--no-edit", "--allow-empty")
	cmd.Dir = r.path

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
