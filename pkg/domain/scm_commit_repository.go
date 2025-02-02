package domain

import "errors"

var ErrScmCommitNotFound = errors.New("scm commit not found")

type ScmCommitRepository interface {
	Get(hash *ScmCommitHash) (*ScmCommit, error)

	AmendAuthor(author *ScmCommitAuthor) error
}
