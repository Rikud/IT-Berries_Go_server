package models

type LoginForm struct {
	email string
	password string
}

func (f LoginForm)getPassword() string {
	return f.password
}

func (f LoginForm)getEmail() string {
	return f.email
}


