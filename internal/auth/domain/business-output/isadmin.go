package businessoutput

type IsAdminOutput struct {
	access bool
}

func NewIsAdminOutput(access bool) *IsAdminOutput {
	return &IsAdminOutput{
		access: access,
	}
}

func (i *IsAdminOutput) GetAccess() bool {
	return i.access
}
