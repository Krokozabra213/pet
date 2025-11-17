package authBusiness

import (
	"fmt"
	"log/slog"
)

func ErrorGateway(op string, log *slog.Logger, errRepo error, errorMap map[error]error, defErr error) error {

	// Прямое совпадение
	if err, exists := errorMap[errRepo]; exists {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Логирование неизвестной ошибки
	log.Warn("unmapped error in business logic",
		"operation", op,
		"error", defErr.Error())

	return fmt.Errorf("%s: %w", op, defErr)
}
