package businessinput

type LogoutInput struct {
	refreshToken string
}

func NewLogoutInput(token string) *LogoutInput {
	return &LogoutInput{
		refreshToken: token,
	}
}

func (input *LogoutInput) GetRefreshToken() string {
	return input.refreshToken
}
