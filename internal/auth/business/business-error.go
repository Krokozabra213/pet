package authBusiness

type BError struct {
	Entity string
	Err    error
}

func BusinessError(entity string, err error) *BError {
	return &BError{
		Entity: entity,
		Err:    err,
	}
}

func (e *BError) Error() string {
	return e.Entity + ": " + e.Err.Error()
}

func (e *BError) Is(target error) bool {
	return e.Error() == target.Error() || e.Err.Error() == target.Error()
}
