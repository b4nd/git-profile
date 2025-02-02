package domain

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type ScmCommitAuthor struct {
	email string
	name  string
}

var ErrInvalidAuthorEmail = errors.New("invalid email")
var ErrInvalidAuthorName = errors.New("invalid name")

func NewScmCommitAuthor(name string, email string) (ScmCommitAuthor, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	if utf8.RuneCountInString(email) == 0 {
		return ScmCommitAuthor{}, ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ScmCommitAuthor{}, ErrInvalidEmail
	}

	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) == 0 {
		return ScmCommitAuthor{}, ErrInvalidAuthorName
	}

	return ScmCommitAuthor{email, name}, nil
}

func (n ScmCommitAuthor) Name() string {
	return n.name
}

func (n ScmCommitAuthor) Email() string {
	return n.email
}

func (n ScmCommitAuthor) Equals(name ScmCommitAuthor) bool {
	return n.name == name.name && n.email == name.email
}

func (n ScmCommitAuthor) String() string {
	return n.name + " <" + n.email + ">"
}
