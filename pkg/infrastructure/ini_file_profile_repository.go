package infrastructure

import (
	"backend/git-profile/pkg/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type IniFileProfileRepository struct {
	path string
}

func NewIniFileProfileRepository(path string) (*IniFileProfileRepository, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	return &IniFileProfileRepository{path}, nil
}

func (i *IniFileProfileRepository) Get(workspace domain.ProfileWorkspace) (*domain.Profile, error) {
	if _, err := os.Stat(i.path); errors.Is(err, os.ErrNotExist) {
		return nil, domain.ErrInvalidWorkspace
	}

	cfg, err := ini.Load(i.path)
	if err != nil {
		return nil, domain.ErrInvalidWorkspace
	}

	section, err := cfg.GetSection(workspace.String())
	if err != nil {
		return nil, domain.ErrInvalidWorkspace
	}

	profile, err := domain.NewProfile(
		workspace.String(),
		section.Key("email").String(),
		section.Key("name").String(),
	)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (i *IniFileProfileRepository) Save(profile *domain.Profile) error {
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

	section, err := cfg.GetSection(profile.Workspace().String())
	if err != nil {
		section, err = cfg.NewSection(profile.Workspace().String())
		if err != nil {
			return err
		}
	}

	section.Key("name").SetValue(profile.Name().String())
	section.Key("email").SetValue(profile.Email().String())

	err = cfg.SaveTo(i.path)
	if err != nil {
		return err
	}

	return nil
}

func (i *IniFileProfileRepository) Delete(workspace domain.ProfileWorkspace) error {
	if _, err := os.Stat(i.path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	cfg, err := ini.Load(i.path)
	if err != nil {
		return err
	}

	cfg.DeleteSection(workspace.String())

	err = cfg.SaveTo(i.path)
	if err != nil {
		return err
	}

	return nil
}

func (i *IniFileProfileRepository) List() ([]*domain.Profile, error) {
	profiles := make([]*domain.Profile, 0)

	if _, err := os.Stat(i.path); errors.Is(err, os.ErrNotExist) {
		return profiles, nil
	}

	cfg, err := ini.Load(i.path)
	if err != nil {
		return nil, domain.ErrInvalidWorkspace
	}

	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}

		profile, err := domain.NewProfile(
			section.Name(),
			section.Key("email").String(),
			section.Key("name").String(),
		)

		if err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
