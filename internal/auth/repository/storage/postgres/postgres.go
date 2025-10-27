package postgres

import (
	"context"
	"time"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"github.com/Krokozabra213/sso/internal/auth/repository/storage"
	"github.com/Krokozabra213/sso/pkg/db"
)

type Postgres struct {
	DB *db.Db
}

func New(db *db.Db) *Postgres {
	return &Postgres{DB: db}
}

func (p *Postgres) SaveUser(
	ctx context.Context, username string, pass string,
) (uid uint, err error) {

	user := &domain.User{
		Username: username,
		Password: pass,
	}
	result := p.DB.DB.Create(user)
	if result.Error != nil {
		if duplicateKey(result.Error) {
			return 0, storage.ErrUserExist
		}

		return 0, storage.ErrUnknown
	}
	return user.ID, nil
}

func (p *Postgres) User(
	ctx context.Context, username string,
) (*domain.User, error) {
	var user domain.User
	result := p.DB.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		if notFound(result.Error) {
			return nil, storage.ErrUserNotFound
		}

		return nil, storage.ErrUnknown
	}
	return &user, nil
}

func (p *Postgres) IsAdmin(
	ctx context.Context, userID int64,
) (bool, error) {
	var admin domain.Admin
	result := p.DB.DB.First(&admin, "user_id = ?", userID)
	if result.Error != nil {
		if notFound(result.Error) {
			return false, nil
		}

		return false, storage.ErrUnknown
	}

	return true, nil
}

func (p *Postgres) App(
	ctx context.Context, appID int,
) (*domain.App, error) {

	var app domain.App
	result := p.DB.DB.First(&app, "id = ?", appID)
	if result.Error != nil {
		if notFound(result.Error) {
			return nil, storage.ErrAppNotFound
		}

		return nil, storage.ErrUnknown
	}
	return &app, nil
}

func (p *Postgres) SaveToken(
	ctx context.Context, hashToken string, exp time.Time,
) (err error) {

	blackToken := &domain.BlackToken{
		Token: hashToken,
		Exp:   exp,
	}

	result := p.DB.DB.Create(blackToken)
	if result.Error != nil {
		if duplicateKey(result.Error) {
			return storage.ErrTokenRevoked
		}
		return storage.ErrUnknown
	}
	return nil
}

func (p *Postgres) CheckToken(
	ctx context.Context, hashToken string, exp time.Time,
) (err error) {

	var token domain.BlackToken

	result := p.DB.DB.Last(&token, "token = ? AND exp = ?", hashToken, exp)
	if result.Error != nil {
		if notFound(result.Error) {
			return nil
		}
		return storage.ErrUnknown
	}
	return storage.ErrTokenRevoked
}
