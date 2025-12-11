package businessinput

type RefreshInput struct {
	refreshToken string
}

func NewRefreshInput(token string) *RefreshInput {
	return &RefreshInput{
		refreshToken: token,
	}
}

func (input *RefreshInput) GetRefreshToken() string {
	return input.refreshToken
}
