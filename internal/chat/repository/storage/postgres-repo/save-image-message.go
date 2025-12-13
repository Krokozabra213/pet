package postgresrepo

import (
	"context"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	"github.com/Krokozabra213/sso/internal/chat/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func (p *Postgres) SaveImageMessage(
	parentCtx context.Context, imageMessage *chatdomain.ImageMessage,
) (*chatdomain.ImageMessage, error) {

	// TODO: ADD TRANSACTION

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	message := imageMessage.GetMessage()

	result := p.DB.Client.WithContext(ctx).Create(message)
	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(customErr)
		return nil, err
	}

	image := imageMessage.GetImage()

	result = p.DB.Client.WithContext(ctx).Create(image)
	customErr = postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(customErr)
		return nil, err
	}

	return imageMessage, nil
}
