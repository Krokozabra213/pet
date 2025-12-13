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

func (input *RegisterInput) GetUsername() string {
	return input.username
}

func (input *RegisterInput) GetPassword() string {
	return input.password
}
