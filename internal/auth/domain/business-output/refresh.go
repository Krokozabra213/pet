package businessoutput

type RefreshOutput struct {
	accessT  string
	refreshT string
}

func NewRefreshOutput(accessToken string, refreshToken string) *RefreshOutput {
	return &RefreshOutput{
		accessT:  accessToken,
		refreshT: refreshToken,
	}
}

func (input *RefreshOutput) GetAccess() string {
	return input.accessT
}

func (input *RefreshOutput) GetRefresh() string {
	return input.refreshT
}
