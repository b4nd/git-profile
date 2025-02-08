package command

import (
	"github.com/b4nd/git-profile/pkg/application"
	"github.com/b4nd/git-profile/pkg/domain"
)

var errorMessages = map[error]string{
	application.ErrProfileAlreadyExists:  "Profile \"%s\" already exists.\n",
	application.ErrProfileNotExists:      "Profile \"%s\" does not exist.\n",
	domain.ErrInvalidEmail:               "The email is invalid.\n",
	domain.ErrInvalidName:                "The name is invalid.\n",
	domain.ErrInvalidWorkspace:           "Profile \"%s\" does not exist.\n",
	domain.ErrInvalidWorkspaceCharacters: "The workspace must contain only alphanumeric characters.\n",
}
