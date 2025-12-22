package main

import (
	"fmt"
	"log"

	chatnewconfig "github.com/Krokozabra213/sso/newconfigs/chat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	path    = "./sql/chat"
	envfile = "chat.env"
	cfgpath = "settings/chat_main.yml"
)

func main() {

	cfg, err := chatnewconfig.Init(cfgpath, envfile)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		cfg.PG.User, cfg.PG.Password, cfg.PG.LocalPort, cfg.PG.DB,
	)

	m, err := migrate.New(
		"file://"+path, // Путь к папке с миграциями
		dbURL,          // Строка подключения к БД
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}
	log.Println("Миграции успешно применены!")
}
