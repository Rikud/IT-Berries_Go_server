package models

type LoginForm struct {
	email    string
	password string
}

func (f LoginForm) GetPassword() string {
	return f.password
}

func (f LoginForm) GetEmail() string {
	return f.email
}
