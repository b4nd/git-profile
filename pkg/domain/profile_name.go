package domain

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type ProfileName struct {
	value string
}

var ErrInvalidName = errors.New("invalid name")

func NewProfileName(value string) (ProfileName, error) {
	name := strings.TrimSpace(value)

	if utf8.RuneCountInString(name) == 0 {
		return ProfileName{}, ErrInvalidName
	}

	return ProfileName{value: name}, nil
}

func (n ProfileName) Value() string {
	return n.value
}

func (n ProfileName) Equals(name ProfileName) bool {
	return n.value == name.value
}

func (n ProfileName) String() string {
	return n.value
}
