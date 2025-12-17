package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrationsPath := "./sql"
	dbURL := "postgres://user_test:password_test@localhost:5556/postgres_test?sslmode=disable"

	m, err := migrate.New(
		"file://"+migrationsPath, // Путь к папке с миграциями
		dbURL,                    // Строка подключения к БД
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}
	log.Println("Миграции успешно применены!")
}
