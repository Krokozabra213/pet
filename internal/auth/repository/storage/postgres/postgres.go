package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
)

const (
	ctxTimeout = 5 * time.Second
)

type Postgres struct {
	DB  *postgrespet.PGDB
	log *slog.Logger
}

func New(db *postgrespet.PGDB, log *slog.Logger) *Postgres {
	return &Postgres{
		DB:  db,
		log: log,
	}
}

func (p *Postgres) SaveUser(
	ctx context.Context, username string, pass string,
) (uid uint, err error) {

	const op = "postgres.SaveUser"
	log := p.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return 0, ErrContext
	}

	user := &domain.User{
		Username: username,
		Password: pass,
	}
	result := p.DB.Client.WithContext(ctx).Create(user)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return 0, err
	}

	log.Info("user saved successfully", "user_id", user.ID)
	return user.ID, nil
}

func (p *Postgres) User(
	ctx context.Context, username string,
) (*domain.User, error) {

	const op = "postgres.User"
	log := p.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return nil, ErrContext
	}

	var user domain.User
	result := p.DB.Client.WithContext(ctx).First(&user, "username = ?", username)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) UserByID(
	ctx context.Context, id int,
) (*domain.User, error) {

	const op = "postgres.UserByID"
	log := p.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return nil, ErrContext
	}

	var user domain.User
	result := p.DB.Client.WithContext(ctx).First(&user, "id = ?", id)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) IsAdmin(
	ctx context.Context, userID int64,
) (bool, error) {

	const op = "postgres.IsAdmin"
	log := p.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return false, ErrContext
	}

	var exists bool
	result := p.DB.Client.WithContext(ctx).Raw(
		"SELECT EXISTS(SELECT 1 FROM admins WHERE user_id = ?)",
		userID,
	).Scan(&exists)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return false, err
	}
	return exists, nil
}

func (p *Postgres) AppByID(
	ctx context.Context, appID int,
) (*domain.App, error) {

	const op = "postgres.AppByID"
	log := p.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return nil, ErrContext
	}

	var app domain.App
	result := p.DB.Client.WithContext(ctx).First(&app, "id = ?", appID)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return nil, err
	}
	return &app, nil
}

func (p *Postgres) SaveToken(
	ctx context.Context, token *domain.BlackToken,
) (err error) {

	const op = "postgres.SaveToken"
	log := p.log.With(
		slog.String("op", op),
	)

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ctxTimeout)
		defer cancel()
	}

	if ctx.Err() != nil {
		log.Error("context error", "err", ctx.Err())
		return ErrContext
	}

	result := p.DB.Client.WithContext(ctx).Create(token)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return err
	}
	return nil
}
