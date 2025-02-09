package domain

type ScmUser struct {
	Workespace string
	Email      string
	Name       string
}

func NewScmUser(workspace string, email string, name string) *ScmUser {
	return &ScmUser{
		Workespace: workspace,
		Email:      email,
		Name:       name,
	}
}
