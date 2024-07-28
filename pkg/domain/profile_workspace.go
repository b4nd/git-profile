package domain

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ProfileWorkspace struct {
	value string
}

var ErrInvalidWorkspace = errors.New("invalid workspace")
var ErrInvalidWorkspaceCharacters = errors.New("invalid workspace characters")

const regexWorkspaceCharacters = `^[a-zA-Z0-9_\.-]+$` // only letters, numbers and underscore

func NewWorkspace(value string) (ProfileWorkspace, error) {
	workspace := strings.TrimSpace(value)
	workspace = strings.ToLower(workspace)

	if utf8.RuneCountInString(workspace) == 0 {
		return ProfileWorkspace{}, ErrInvalidWorkspace
	}

	if mateh, _ := regexp.MatchString(regexWorkspaceCharacters, workspace); !mateh {
		return ProfileWorkspace{}, ErrInvalidWorkspaceCharacters
	}

	return ProfileWorkspace{value: workspace}, nil
}

func (w ProfileWorkspace) Value() string {
	return w.value
}

func (w ProfileWorkspace) Equals(workspace ProfileWorkspace) bool {
	return w.value == workspace.value
}

func (w ProfileWorkspace) String() string {
	return w.value
}
