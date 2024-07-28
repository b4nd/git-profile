package command

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
)

var errorMessages = map[error]string{
	application.ErrProfileAlreadyExists:  "Profile \"%s\" already exists.",
	application.ErrProfileNotExists:      "Profile \"%s\" does not exist.",
	domain.ErrInvalidEmail:               "The email is invalid.",
	domain.ErrInvalidName:                "The name is invalid.",
	domain.ErrInvalidWorkspace:           "The workspace is invalid.",
	domain.ErrInvalidWorkspaceCharacters: "The workspace must contain only alphanumeric characters.",
}
