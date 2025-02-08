package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/b4nd/git-profile/pkg/domain"

	"gopkg.in/ini.v1"
)

type IniFileProfileRepository struct {
	paths []string
}

func NewIniFileProfileRepository(paths []string) (*IniFileProfileRepository, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no paths provided")
	}

	return &IniFileProfileRepository{paths}, nil
}

type iniFileSource struct {
	path string
	cfg  *ini.File
}

func (i *IniFileProfileRepository) load() ([]*iniFileSource, error) {
	var cfgs []*iniFileSource = make([]*iniFileSource, 0)
	for _, path := range i.paths {
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
	path := i.paths[0]

	// Create the file if it does not exist
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			return err
		}

		path = filepath.Clean(path)
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer file.Close()
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

	// If there is no profile to delete, return nil to indicate that the profile does not exist
	if len(sources) == 0 {
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
