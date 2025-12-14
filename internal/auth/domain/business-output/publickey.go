package businessoutput

type PublicKeyOutput struct {
	publicKey string
}

func NewPublicKeyOutput(publicKey string) *PublicKeyOutput {
	return &PublicKeyOutput{
		publicKey: publicKey,
	}
}

func (i *PublicKeyOutput) GetPublicKey() string {
	return i.publicKey
}
