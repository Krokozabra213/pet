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

func (i *LoginInput) GetUsername() string {
	return i.username
}

func (i *LoginInput) GetPassword() string {
	return i.password
}

func (i *LoginInput) GetAppID() int {
	return i.appID
}
