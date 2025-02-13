package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type ScmCommitHash struct {
	value string
}

var ErrInvalidHash = errors.New("invalid hash or branch name")

func NewScmCommitHash(value string) (ScmCommitHash, error) {
	name := strings.TrimSpace(value)

	if utf8.RuneCountInString(name) == 0 {
		return ScmCommitHash{}, ErrInvalidHash
	}

	if strings.Contains(name, " ") {
		return ScmCommitHash{}, ErrInvalidHash
	}

	return ScmCommitHash{value: name}, nil
}

func NewScmCommitHashHead() ScmCommitHash {
	hash, _ := NewScmCommitHash("HEAD")
	return hash
}

func (n ScmCommitHash) String() string {
	return n.value
}
