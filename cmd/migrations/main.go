package main

import (
	"flag"
	"fmt"

	"github.com/Krokozabra213/sso/internal/auth/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = flag.String("host", "", "host for connect db")
	user     = flag.String("user", "", "username for connect db")
	password = flag.String("password", "", "password for connect db")
	dbname   = flag.String("dbname", "", "dbname for connect db")
	port     = flag.String("port", "", "port for connect db")
	sslmode  = flag.String("sslmode", "", "sslmode for connect db")
)

func main() {
	flag.Parse()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		*host, *user, *password, *dbname, *port, *sslmode,
	)
	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.App{}, &domain.User{})
}
