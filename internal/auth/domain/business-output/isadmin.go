package businessoutput

type IsAdminOutput struct {
	access bool
}

func NewIsAdminOutput(access bool) *IsAdminOutput {
	return &IsAdminOutput{
		access: access,
	}
}

func (input *IsAdminOutput) GetAccess() bool {
	return input.access
}
