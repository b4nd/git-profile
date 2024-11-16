package infrastructure

import (
	"backend/git-profile/pkg/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const GIT_CONFIG_FILE = ".git/config"
const GIT_SECTION_USER = "user"

type IniFileScmUserRepository struct {
	path string
}

func NewIniFileScmUserRepository(path string) (*IniFileScmUserRepository, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	path = filepath.Join(path, GIT_CONFIG_FILE)

	return &IniFileScmUserRepository{path}, nil
}

func (i *IniFileScmUserRepository) Get() (*domain.ScmUser, error) {
	if _, err := os.Stat(i.path); errors.Is(err, os.ErrNotExist) {
		return nil, domain.ErrScmUserNotFound
	}

	cfg, err := ini.Load(i.path)
	if err != nil {
		return nil, domain.ErrScmUserNotFound
	}

	section, err := cfg.GetSection(GIT_SECTION_USER)
	if err != nil {
		return nil, domain.ErrScmUserNotFound
	}

	user := domain.NewScmUser(
		section.Key("workspace").String(),
		section.Key("email").String(),
		section.Key("name").String(),
	)

	return user, nil
}

func (i *IniFileScmUserRepository) Save(user *domain.ScmUser) error {
	if _, err := os.Stat(i.path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(i.path), os.ModePerm); err != nil {
			return err
		}

		file, err := os.Create(i.path)
		if err != nil {
			return err
		}

		file.Close()
	}

	cfg, err := ini.Load(i.path)
	if err != nil {
		return domain.ErrInvalidWorkspace
	}

	section, err := cfg.GetSection(GIT_SECTION_USER)
	if err != nil {
		section, err = cfg.NewSection(GIT_SECTION_USER)
		if err != nil {
			return err
		}
	}

	section.Key("workspace").SetValue(user.Workespace)
	section.Key("name").SetValue(user.Name)
	section.Key("email").SetValue(user.Email)

	err = cfg.SaveTo(i.path)
	if err != nil {
		return err
	}

	return nil
}
