package businessinput

type IsAdminInput struct {
	userID int64
}

func NewIsAdminInput(userID int64) *IsAdminInput {
	return &IsAdminInput{
		userID: userID,
	}
}

func (i *IsAdminInput) GetUserID() int64 {
	return i.userID
}
