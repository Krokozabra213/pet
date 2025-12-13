package businessoutput

type LoginOutput struct {
	accessT  string
	refreshT string
}

func NewLoginOutput(accessToken string, refreshToken string) *LoginOutput {
	return &LoginOutput{
		accessT:  accessToken,
		refreshT: refreshToken,
	}
}

func (input *LoginOutput) GetAccess() string {
	return input.accessT
}

func (input *LoginOutput) GetRefresh() string {
	return input.refreshT
}
