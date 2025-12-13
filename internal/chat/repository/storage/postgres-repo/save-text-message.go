package postgresrepo

import (
	"context"

	chatdomain "github.com/Krokozabra213/sso/internal/chat/domain"
	"github.com/Krokozabra213/sso/internal/chat/repository/storage"
	contexthandler "github.com/Krokozabra213/sso/pkg/context-handler"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

func (p *Postgres) SaveTextMessage(
	parentCtx context.Context, textMessage *chatdomain.TextMessage,
) (*chatdomain.TextMessage, error) {

	// TODO: ADD TRANSACTION

	ctx, cancel := contexthandler.EnsureCtxTimeout(parentCtx, ctxTimeout)
	defer cancel()

	if ctx.Err() != nil {
		return nil, storage.CtxError(ctx.Err())
	}

	message := textMessage.GetMessage()

	result := p.DB.Client.WithContext(ctx).Create(message)
	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(customErr)
		return nil, err
	}

	text := textMessage.GetText()

	result = p.DB.Client.WithContext(ctx).Create(text)
	customErr = postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		err := ErrorFactory(customErr)
		return nil, err
	}

	return textMessage, nil
}
