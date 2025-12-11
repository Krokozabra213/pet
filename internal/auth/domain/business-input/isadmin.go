package businessinput

type IsAdminInput struct {
	userID int64
}

func NewIsAdminInput(userID int64) *IsAdminInput {
	return &IsAdminInput{
		userID: userID,
	}
}

func (input *IsAdminInput) GetUserID() int64 {
	return input.userID
}
