package postgrespet

import (
	"errors"
	"strings"
)

type CustomError struct {
	Message string
	Err     error
}

func NewError(message string, err error) *CustomError {
	return &CustomError{
		Message: message,
		Err:     err,
	}
}

func (e *CustomError) Error() string {
	if e == nil {
		return ""
	}

	if e.Message != "" && e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	if e.Message != "" {
		return e.Message
	}

	if e.Err != nil {
		return e.Err.Error()
	}

	return ErrInternal.Error()
}

func (e *CustomError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *CustomError) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	if e.Message == target.Error() {
		return true
	}

	if targetCustom, ok := target.(*CustomError); ok {
		if e.Message == targetCustom.Message {
			return true
		}

		if e.Message != "" && targetCustom.Message != "" {
			return strings.Contains(e.Message, targetCustom.Message) ||
				strings.Contains(targetCustom.Message, e.Message)
		}
	}

	if e.Err != nil {
		return errors.Is(e.Err, target)
	}

	return e.Error() == target.Error()
}
