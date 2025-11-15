package postgrespet

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	ErrorTypeDuplicateKey = "DUPLICATE_KEY"
	ErrorTypeNotFound     = "NOT_FOUND"
	ErrorTypeValidation   = "VALIDATION_ERROR"
	ErrorTypeInternal     = "INTERNAL_ERROR"
)

type CustomError struct {
	Type    string
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return e.Err.Error()
}

func (e *CustomError) GetType() string {
	return e.Type
}

func (e *CustomError) GetMessage() string {
	return e.Message
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func ErrorHandler(err error) *CustomError {
	if err == nil {
		return nil
	}

	// Проверяем типы ошибок в порядке приоритета
	switch {
	case isDuplicateKey(err):
		return &CustomError{
			Type:    ErrorTypeDuplicateKey,
			Message: "Запись уже существует",
			Err:     err,
		}
	case isNotFound(err):
		return &CustomError{
			Type:    ErrorTypeNotFound,
			Message: "Запись не найдена",
			Err:     err,
		}
	default:
		return &CustomError{
			Type:    ErrorTypeInternal,
			Message: "Внутренняя ошибка сервера",
			Err:     err,
		}
	}
}

// Вспомогательные функции
func isDuplicateKey(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
