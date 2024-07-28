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

func (g *ScmUser) Equals(user *ScmUser) bool {
	return g.Name == user.Name &&
		g.Email == user.Email &&
		g.Workespace == user.Workespace
}
