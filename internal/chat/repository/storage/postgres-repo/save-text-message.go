package postgresrepo

import (
	"context"
	"errors"
	"time"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	"github.com/Krokozabra213/sso/internal/chat/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
	"gorm.io/gorm"
)

func (p *Postgres) SaveTextMessage(
	parentCtx context.Context, textMessage *chatdomain.TextMessage,
) (*chatdomain.TextMessage, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	perform := func() error {
		return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			message := textMessage.GetMessage()

			if err := tx.Create(message).Error; err != nil {
				return err
			}

			text := textMessage.GetText()
			text.SetMessageID(message.ID)

			if err := tx.Create(text).Error; err != nil {
				return err
			}
			return nil
		})
	}

	var err error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		err = perform()
		if err == nil {
			return textMessage, nil
		}

		customErr := postgrespet.ErrorWrapper(err)
		if errors.Is(customErr, postgrespet.ErrTransaction) {
			if attempt < MaxRetries-1 {
				delay := ImageTransactionTimeout * time.Duration(attempt+1)
				time.Sleep(delay)
				continue
			}
		}

		return nil, ErrorFactory(customErr)
	}

	return textMessage, nil
}
