package domain

type Profile struct {
	workspace ProfileWorkspace
	email     ProfileEmail
	name      ProfileName
}

func NewProfile(workspace string, email string, name string) (*Profile, error) {
	w, err := NewWorkspace(workspace)
	if err != nil {
		return nil, err
	}

	e, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	n, err := NewName(name)
	if err != nil {
		return nil, err
	}

	return &Profile{
		workspace: w,
		email:     e,
		name:      n,
	}, nil
}

func (p Profile) Workspace() ProfileWorkspace {
	return p.workspace
}

func (p Profile) Email() ProfileEmail {
	return p.email
}

func (p Profile) Name() ProfileName {
	return p.name
}

func (p Profile) Equals(profile *Profile) bool {
	return p.workspace.Equals(profile.workspace) &&
		p.email.Equals(profile.email) &&
		p.name.Equals(profile.name)
}
