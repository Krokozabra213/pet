package postgres

import (
	"context"
	"log/slog"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	postgrespet "github.com/Krokozabra213/sso/pkg/db/postgres-pet"
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
	)

	user := &domain.User{
		Username: username,
		Password: pass,
	}
	result := p.DB.Client.Create(user)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return 0, err
	}
	return user.ID, nil
}

func (p *Postgres) User(
	ctx context.Context, username string,
) (*domain.User, error) {

	const op = "postgres.User"
	log := p.log.With(
		slog.String("op", op),
	)

	var user domain.User
	result := p.DB.Client.First(&user, "username = ?", username)

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

	var user domain.User
	result := p.DB.Client.First(&user, "id = ?", id)

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

	var exists bool
	result := p.DB.Client.Raw(
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

	var app domain.App
	result := p.DB.Client.First(&app, "id = ?", appID)

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

	result := p.DB.Client.Create(token)

	customErr := postgrespet.ErrorWrapper(result.Error)
	if customErr != nil {
		log.Error("postgres error", "err", customErr.Error())
		err := ErrorFactory(customErr)
		return err
	}
	return nil
}
