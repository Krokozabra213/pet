package authBusiness

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IBusinessError interface {
	ToGRPC() error
	error
}

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

func (e *BError) ToGRPC() error {
	return status.Error(e.GRPCCode(), e.Error())
}

func (e *BError) GRPCCode() codes.Code {
	switch {
	case errors.Is(e.Err, ErrTimeout):
		return codes.DeadlineExceeded // Операция не завершилась за отведенное время.

	case errors.Is(e.Err, ErrExists):
		return codes.AlreadyExists // Попытка создать сущность, которая уже существует.

	case errors.Is(e.Err, ErrNotFound):
		return codes.NotFound // Запрошенная сущность (например, пользователь) не найдена.

	case errors.Is(e.Err, ErrInternal) || errors.Is(e.Err, ErrHashPassword) || errors.Is(e.Err, ErrTokenGenerate):
		return codes.Internal // Внутренняя ошибка сервера, не зависящая от клиента.

	case errors.Is(e.Err, ErrInvalidCredentials) || errors.Is(e.Err, ErrTokenExpired) || errors.Is(e.Err, ErrTokenRevoked):
		return codes.Unauthenticated

	case errors.Is(e.Err, ErrParse):
		return codes.InvalidArgument // Ошибка парсинга данных от клиента (напр., невалидный JWT).

	case errors.Is(e.Err, ErrPermission):
		return codes.PermissionDenied // Пользователь аутентифицирован, но у него нет прав.

	default:
		// Для любых неизвестных ошибок
		return codes.Internal
	}
}
