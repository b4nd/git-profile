package domain

type Profile struct {
	workspace ProfileWorkspace
	email     ProfileEmail
	name      ProfileName
}

const NotConfiguredWorkspace = "(not configured)"

func NewProfile(workspace string, email string, name string) (*Profile, error) {
	w, err := NewProfileWorkspace(workspace)
	if err != nil {
		return nil, err
	}

	e, err := NewProfileEmail(email)
	if err != nil {
		return nil, err
	}

	n, err := NewProfileName(name)
	if err != nil {
		return nil, err
	}

	return &Profile{
		workspace: w,
		email:     e,
		name:      n,
	}, nil
}

func NewProfileWithoutWorkspace(email string, name string) (*Profile, error) {
	e, err := NewProfileEmail(email)
	if err != nil {
		return nil, err
	}

	n, err := NewProfileName(name)
	if err != nil {
		return nil, err
	}

	return &Profile{
		workspace: ProfileWorkspace{
			value: NotConfiguredWorkspace,
		},
		email: e,
		name:  n,
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
