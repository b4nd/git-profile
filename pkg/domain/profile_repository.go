package domain

type ProfileRepository interface {
	Get(workspace ProfileWorkspace) (*Profile, error)

	Save(profile *Profile) error

	Delete(workspace ProfileWorkspace) error

	List() ([]*Profile, error)
}
