package entities

type User struct {
	email    string
	password string
	username string
	avatar   string
}

func (user *User) GetEmailPoint() *string {
	return &(user.email)
}

func (user *User) GetPasswordPoint() *string {
	return &(user.password)
}

func (user *User) GetUsernamePoint() *string {
	return &(user.username)
}

func (user *User) GetAvatarPoint() *string {
	return &(user.email)
}

func (user *User) GetEmail() string {
	return user.email
}

func (user *User) GetPassword() string {
	return user.password
}

func (user *User) GetUsername() string {
	return user.username
}

func (user *User) GetAvatar() string {
	return user.email
}
