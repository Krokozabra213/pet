package postgrespet

import "errors"

const (
	// logic errors
	CodeUniqueViolation           = "23505"
	CodeNotNullViolation          = "23502"
	CodeForeignKeyViolation       = "23503"
	CodeCheckViolation            = "23514"
	CodeStringTooLong             = "22001"
	CodeNumericOutOfRange         = "22003"
	CodeInvalidDatetimeFormat     = "22007"
	CodeDatetimeOverflow          = "22008"
	CodeDivisionByZero            = "22012"
	CodeInvalidTextRepresentation = "22P02"

	// transaction errors
	CodeSerializationFailure = "40001" // Serialization failure
	CodeDeadlockDetected     = "40P01" // Deadlock detected
	CodeLockNotAvailable     = "55P03" // Lock not available / timeout
)

var (
	ErrDuplicateKey = errors.New("DUPLICATE_KEY_ERROR")
	ErrNotFound     = errors.New("NOT_FOUND_ERROR")
	ErrValidation   = errors.New("VALIDATION_ERROR")
	ErrCtxCancelled = errors.New("CTX_CANCELLED_ERROR")
	ErrCtxDeadline  = errors.New("CTX_DEADLINE_ERROR")
	ErrInternal     = errors.New("INTERNAL_ERROR")
	ErrTransaction  = errors.New("TRANSACTION_ERROR")
)
