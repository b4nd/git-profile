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
	path   string
	others []string
}

func NewIniFileProfileRepository(path string, others []string) (*IniFileProfileRepository, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	return &IniFileProfileRepository{path, others}, nil
}

type iniFileSource struct {
	path string
	cfg  *ini.File
}

func (i *IniFileProfileRepository) load() ([]*iniFileSource, error) {
	var cfgs []*iniFileSource = make([]*iniFileSource, 0)
	cfg, err := ini.Load(i.path)
	if err != nil {
		return nil, domain.ErrInvalidWorkspace
	}
	cfgs = append(cfgs, &iniFileSource{i.path, cfg})

	for _, path := range i.others {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			cfg, err := ini.Load(path)
			if err != nil {
				return nil, domain.ErrInvalidWorkspace
			}

			cfgs = append(cfgs, &iniFileSource{path, cfg})
		}
	}

	return cfgs, nil
}

func (i *IniFileProfileRepository) Get(workspace domain.ProfileWorkspace) (*domain.Profile, error) {
	sources, err := i.load()
	if err != nil {
		return nil, err
	}

	for _, source := range sources {
		section, err := source.cfg.GetSection(workspace.String())
		if err != nil {
			continue
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

	return nil, domain.ErrInvalidWorkspace
}

func (i *IniFileProfileRepository) Save(profile *domain.Profile) error {
	// Create the file if it does not exist
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

	sources, err := i.load()
	if err != nil {
		return err
	}

	var source *iniFileSource = sources[0]
	var section *ini.Section = nil
	for _, c := range sources {
		s, err := c.cfg.GetSection(profile.Workspace().String())
		if err == nil {
			source = c
			section = s
		}
	}

	if section == nil {
		section, err = source.cfg.NewSection(profile.Workspace().String())
		if err != nil {
			return err
		}
	}

	section.Key("name").SetValue(profile.Name().String())
	section.Key("email").SetValue(profile.Email().String())

	err = source.cfg.SaveTo(source.path)
	if err != nil {
		return err
	}

	return nil
}

func (i *IniFileProfileRepository) Delete(workspace domain.ProfileWorkspace) error {
	sources, err := i.load()
	if err != nil {
		return nil
	}

	var source *iniFileSource = sources[0]
	var section *ini.Section = nil
	for _, c := range sources {
		s, err := c.cfg.GetSection(workspace.String())
		if err == nil {
			source = c
			section = s
		}
	}

	if section == nil {
		return nil
	}

	source.cfg.DeleteSection(workspace.String())
	err = source.cfg.SaveTo(source.path)
	if err != nil {
		return err
	}

	err = source.cfg.Reload()
	if err != nil {
		return err
	}

	return nil
}

func (i *IniFileProfileRepository) List() ([]*domain.Profile, error) {
	profiles := make([]*domain.Profile, 0)

	sources, err := i.load()
	if err != nil {
		return profiles, nil
	}

	for _, source := range sources {
		for _, section := range source.cfg.Sections() {
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
	}

	return profiles, nil
}
