package businessinput

type PublicKeyInput struct {
	appID int
}

func NewPublicKeyInput(appID int) *PublicKeyInput {
	return &PublicKeyInput{
		appID: appID,
	}
}

func (i *PublicKeyInput) GetAppID() int {
	return i.appID
}
