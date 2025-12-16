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

func (i *RefreshOutput) GetAccess() string {
	return i.accessT
}

func (i *RefreshOutput) GetRefresh() string {
	return i.refreshT
}
