package businessoutput

type PublicKeyOutput struct {
	publicKey string
}

func NewPublicKeyOutput(publicKey string) *PublicKeyOutput {
	return &PublicKeyOutput{
		publicKey: publicKey,
	}
}

func (input *PublicKeyOutput) GetPublicKey() string {
	return input.publicKey
}
