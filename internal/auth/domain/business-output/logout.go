package businessoutput

type LogoutOutput struct {
	success bool
}

func NewLogoutOutput(success bool) *LogoutOutput {
	return &LogoutOutput{
		success: success,
	}
}

func (i *LogoutOutput) GetSuccess() bool {
	return i.success
}
