package domain

import "errors"

var ErrScmUserNotFound = errors.New("scm user not found")

type ScmUserRepository interface {
	Get() (*ScmUser, error)

	Save(user *ScmUser) error

	Delete() error
}
