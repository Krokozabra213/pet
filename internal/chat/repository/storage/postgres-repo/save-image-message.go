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

const (
	ImageTransactionTimeout = 15 * time.Millisecond
)

func (p *Postgres) SaveImageMessage(
	parentCtx context.Context, imageMessage *chatdomain.ImageMessage,
) (*chatdomain.ImageMessage, error) {

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	perform := func() error {
		return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			message := imageMessage.GetMessage()

			if err := tx.Create(message).Error; err != nil {
				return err
			}

			image := imageMessage.GetImage()
			image.SetMessageID(message.ID)

			if err := tx.Create(image).Error; err != nil {
				return err
			}
			return nil
		})
	}

	var err error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		err = perform()
		if err == nil {
			return imageMessage, nil
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

	return imageMessage, nil
}
