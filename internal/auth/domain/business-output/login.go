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

func (i *LoginOutput) GetAccess() string {
	return i.accessT
}

func (i *LoginOutput) GetRefresh() string {
	return i.refreshT
}
