package businessinput

type RegisterInput struct {
	username string
	password string
}

func NewRegisterInput(username string, password string) *RegisterInput {
	return &RegisterInput{
		username: username,
		password: password,
	}
}

func (i *RegisterInput) GetUsername() string {
	return i.username
}

func (i *RegisterInput) GetPassword() string {
	return i.password
}
