package command

import (
	"backend/git-profile/pkg/application"
	"backend/git-profile/pkg/domain"
)

var errorMessages = map[error]string{
	application.ErrProfileAlreadyExists:  "Profile \"%s\" already exists.\n",
	application.ErrProfileNotExists:      "Profile \"%s\" does not exist.\n",
	domain.ErrInvalidEmail:               "The email is invalid.\n",
	domain.ErrInvalidName:                "The name is invalid.\n",
	domain.ErrInvalidWorkspace:           "Profile \"%s\" does not exist.\n",
	domain.ErrInvalidWorkspaceCharacters: "The workspace must contain only alphanumeric characters.\n",
}
