package domain

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

type ProfileEmail struct {
	value string
}

var ErrInvalidEmail = errors.New("invalid email")

func NewProfileEmail(value string) (ProfileEmail, error) {
	email := strings.TrimSpace(value)
	email = strings.ToLower(email)

	if utf8.RuneCountInString(email) == 0 {
		return ProfileEmail{}, ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ProfileEmail{}, ErrInvalidEmail
	}

	return ProfileEmail{value: email}, nil
}

func (e ProfileEmail) Equals(email ProfileEmail) bool {
	return e.value == email.value
}

func (e ProfileEmail) String() string {
	return e.value
}
