package postgres

import (
	"context"

	"github.com/Krokozabra213/sso/internal/auth/domain"
)

type Postgres struct {
}

func New() *Postgres {
	return &Postgres{}
}

func (p *Postgres) SaveUser(
	ctx context.Context, username string, passHash []byte,
) (uid int64, err error) {
	//TODO:...
	return 10, nil
}

func (p *Postgres) User(
	ctx context.Context, username string,
) (domain.User, error) {
	//TODO:...

	user := domain.User{
		Username: "test",
		Password: "test",
	}
	user.ID = 5
	return user, nil
}

func (p *Postgres) IsAdmin(
	ctx context.Context, userID int64,
) (bool, error) {
	//TODO:...
	return true, nil
}

func (p *Postgres) App(
	ctx context.Context, appID int,
) (domain.App, error) {
	//TODO:...
	app := domain.App{
		Name:  "app1",
		Sault: "test123",
	}
	app.ID = 5
	return app, nil
}
