package businessinput

type PublicKeyInput struct {
	appID int
}

func NewPublicKeyInput(appID int) *PublicKeyInput {
	return &PublicKeyInput{
		appID: appID,
	}
}

func (input *PublicKeyInput) GetAppID() int {
	return input.appID
}
