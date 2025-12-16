package businessinput

type LogoutInput struct {
	refreshToken string
}

func NewLogoutInput(token string) *LogoutInput {
	return &LogoutInput{
		refreshToken: token,
	}
}

func (i *LogoutInput) GetRefreshToken() string {
	return i.refreshToken
}
