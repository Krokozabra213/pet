package businessoutput

type LogoutOutput struct {
	success bool
}

func NewLogoutOutput(success bool) *LogoutOutput {
	return &LogoutOutput{
		success: success,
	}
}

func (input *LogoutOutput) GetSuccess() bool {
	return input.success
}
