package infrastructure_test

import (
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/b4nd/git-profile/pkg/domain"
	"github.com/b4nd/git-profile/pkg/infrastructure"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

const WorkspaceInvalid = "non-existing-workspace"

const ProfileTemplate string = `{{range .}}
[{{.Workspace}}]
name = {{.Name}}
email = {{.Email}}
{{end}}
`

type Profile struct {
	Workspace string
	Name      string
	Email     string
}

// generateTempFileAndProfiles generates a temporary file with profiles
// and returns the file, the profiles and a function to close and remove the file
func generateTempFileAndProfiles(t *testing.T, length uint) (*os.File, []Profile, func()) {
	faker := faker.New()

	file, err := os.CreateTemp("", ".gitprofile")
	assert.NoError(t, err)

	template, err := template.New("ini").Parse(ProfileTemplate)
	assert.NoError(t, err)

	profiles := []Profile{}

	for i := 0; i < int(length); i++ {
		profile := Profile{
			Workspace: faker.Internet().User() + "-" + faker.RandomStringWithLength(10),
			Email:     faker.Internet().Email(),
			Name:      faker.Person().Name(),
		}
		profiles = append(profiles, profile)
	}

	err = template.Execute(file, profiles)
	assert.NoError(t, err)

	return file, profiles, func() {
		file.Close()
		os.Remove(file.Name())
	}
}

func TestIniFileProfileRepository(t *testing.T) {
	faker := faker.New()

	t.Run("should return error when file is not a valid ini file", func(t *testing.T) {
		repo, err := infrastructure.NewIniFileProfileRepository([]string{})
		assert.Error(t, err)
		assert.Nil(t, repo)
	})

	t.Run("should return the profile when it exists", func(t *testing.T) {
		file, profiles, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		for _, profile := range profiles {
			profile, err := domain.NewProfile(
				profile.Workspace,
				profile.Email,
				profile.Name,
			)
			assert.NoError(t, err)

			gettedProfile, err := iniFileProfileRepository.Get(profile.Workspace())
			assert.NoError(t, err)
			assert.True(t, profile.Equals(gettedProfile))
		}
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		workspace, err := domain.NewProfileWorkspace(WorkspaceInvalid)
		assert.NoError(t, err)

		_, err = iniFileProfileRepository.Get(workspace)
		assert.ErrorIs(t, err, domain.ErrInvalidWorkspace)
	})

	t.Run("should return error when file is empty", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 0)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		workspace, err := domain.NewProfileWorkspace(WorkspaceInvalid)
		assert.NoError(t, err)

		_, err = iniFileProfileRepository.Get(workspace)
		assert.ErrorIs(t, err, domain.ErrInvalidWorkspace)
	})

	t.Run("should return error when file does not exist", func(t *testing.T) {
		path := path.Join(os.TempDir(), faker.RandomStringWithLength(10))
		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{path})
		assert.NoError(t, err)
		assert.NoFileExists(t, path)

		workspace, err := domain.NewProfileWorkspace(WorkspaceInvalid)
		assert.NoError(t, err)

		_, err = iniFileProfileRepository.Get(workspace)
		assert.ErrorIs(t, err, domain.ErrInvalidWorkspace)
	})

	t.Run("should return all profiles", func(t *testing.T) {
		file, profiles, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, len(profiles))

		for i, profile := range profiles {
			profile, err := domain.NewProfile(
				profile.Workspace,
				profile.Email,
				profile.Name,
			)
			assert.NoError(t, err)

			assert.True(t, profile.Equals(gettedProfiles[i]))
		}
	})

	t.Run("should return array empty when file is empty", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 0)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, 0)
	})

	t.Run("should return array empty when file does not exist", func(t *testing.T) {
		path := path.Join(os.TempDir(), faker.RandomStringWithLength(10))
		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{path})
		assert.NoError(t, err)
		assert.NoFileExists(t, path)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, 0)
	})

	t.Run("should create profile", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		profile, err := domain.NewProfile(
			faker.Internet().User()+"-new",
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Save(profile)
		assert.NoError(t, err)

		gettedProfile, err := iniFileProfileRepository.Get(profile.Workspace())
		assert.NoError(t, err)

		assert.True(t, profile.Equals(gettedProfile))
	})

	t.Run("should update existing profile", func(t *testing.T) {
		file, profiles, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		profile, err := domain.NewProfile(
			profiles[0].Workspace,
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Save(profile)
		assert.NoError(t, err)

		gettedProfile, err := iniFileProfileRepository.Get(profile.Workspace())
		assert.NoError(t, err)

		assert.True(t, profile.Equals(gettedProfile))
	})

	t.Run("should create profile to empty file", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 0)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		profile, err := domain.NewProfile(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Save(profile)
		assert.NoError(t, err)

		gettedProfile, err := iniFileProfileRepository.Get(profile.Workspace())
		assert.NoError(t, err)

		assert.True(t, profile.Equals(gettedProfile))
	})

	t.Run("should create profile to non existing file", func(t *testing.T) {
		path := path.Join(os.TempDir(), faker.RandomStringWithLength(10))
		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{path})
		assert.NoError(t, err)
		assert.NoFileExists(t, path)

		profile, err := domain.NewProfile(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Save(profile)
		assert.NoError(t, err)
		assert.FileExists(t, path)

		gettedProfile, err := iniFileProfileRepository.Get(profile.Workspace())
		assert.NoError(t, err)

		assert.True(t, profile.Equals(gettedProfile))
	})

	t.Run("should delete profile", func(t *testing.T) {
		file, profiles, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		profile, err := domain.NewProfile(
			profiles[0].Workspace,
			profiles[0].Email,
			profiles[0].Name,
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(profile.Workspace())
		assert.NoError(t, err)

		_, err = iniFileProfileRepository.Get(profile.Workspace())
		assert.ErrorIs(t, err, domain.ErrInvalidWorkspace)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, len(profiles)-1)
	})

	t.Run("should not return error when profile does not exist", func(t *testing.T) {
		file, profiles, closeAndRemoveFile := generateTempFileAndProfiles(t, 100)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		workspace, err := domain.NewProfileWorkspace(WorkspaceInvalid)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(workspace)
		assert.NoError(t, err)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, len(profiles))
	})

	t.Run("should not return error when file is empty", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 0)
		defer closeAndRemoveFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name()})
		assert.NoError(t, err)

		profile, err := domain.NewProfile(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(profile.Workspace())
		assert.NoError(t, err)

		gettedProfiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)

		assert.Len(t, gettedProfiles, 0)
	})

	t.Run("should not return error when file does not exist", func(t *testing.T) {
		path := path.Join(os.TempDir(), faker.RandomStringWithLength(10))
		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{path})
		assert.NoError(t, err)
		assert.NoFileExists(t, path)

		profile, err := domain.NewProfile(
			faker.Internet().User(),
			faker.Internet().Email(),
			faker.Person().Name(),
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(profile.Workspace())
		assert.NoError(t, err)
	})

	t.Run("should delete profile from other file", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 2)
		defer closeAndRemoveFile()

		otherFile, otherProfiles, closeAndRemoveOtherFile := generateTempFileAndProfiles(t, 2)
		defer closeAndRemoveOtherFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name(), otherFile.Name()})
		assert.NoError(t, err)

		otherProfile, err := domain.NewProfile(
			otherProfiles[0].Workspace,
			otherProfiles[0].Email,
			otherProfiles[0].Name,
		)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(otherProfile.Workspace())
		assert.NoError(t, err)

		profiles, err := iniFileProfileRepository.List()
		assert.NoError(t, err)
		assert.Len(t, profiles, 3)
	})

	t.Run("should return no error when profile does not exist in other file", func(t *testing.T) {
		file, _, closeAndRemoveFile := generateTempFileAndProfiles(t, 2)
		defer closeAndRemoveFile()

		otherFile, _, closeAndRemoveOtherFile := generateTempFileAndProfiles(t, 2)
		defer closeAndRemoveOtherFile()

		iniFileProfileRepository, err := infrastructure.NewIniFileProfileRepository([]string{file.Name(), otherFile.Name()})
		assert.NoError(t, err)

		workspace, err := domain.NewProfileWorkspace(WorkspaceInvalid)
		assert.NoError(t, err)

		err = iniFileProfileRepository.Delete(workspace)
		assert.NoError(t, err)
	})
}
