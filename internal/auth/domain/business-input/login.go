package businessinput

type LoginInput struct {
	username string
	password string
	appID    int
}

func NewLoginInput(username string, password string, appID int) *LoginInput {
	return &LoginInput{
		username: username,
		password: password,
		appID:    appID,
	}
}

func (input *LoginInput) GetUsername() string {
	return input.username
}

func (input *LoginInput) GetPassword() string {
	return input.password
}

func (input *LoginInput) GetAppID() int {
	return input.appID
}
