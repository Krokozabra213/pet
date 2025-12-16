package businessinput

type RefreshInput struct {
	refreshToken string
}

func NewRefreshInput(token string) *RefreshInput {
	return &RefreshInput{
		refreshToken: token,
	}
}

func (i *RefreshInput) GetRefreshToken() string {
	return i.refreshToken
}
