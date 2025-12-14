package businessoutput

type RegisterOutput struct {
	userID uint64
}

func NewRegisterOutput(userID uint64) *RegisterOutput {
	return &RegisterOutput{
		userID: userID,
	}
}

func (i *RegisterOutput) GetUserID() uint64 {
	return i.userID
}
