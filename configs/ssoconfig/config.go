package ssoconfig

import (
	"log"

	"github.com/Krokozabra213/sso/configs"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

var (
	emptyValue = ""
	op         = "ssoconfig: "
)

type Config struct {
	Security *Security
	Server   *configs.Server
	DB       *configs.DB
}

type Security struct {
	PrivateKey      []byte
	Secret          []byte
	AccessTokenTTL  int
	RefreshTokenTTL int
}

func Load(env string) *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%s%s", op, configs.ErrLoadConfig.Error())
	}

	privateKeyPath := getEnvOrFatal(PrivateKey)
	privateKey, err := PrivateKeyData(privateKeyPath)
	if err != nil {
		log.Fatalf("%s%s", op, err.Error())
	}

	accessTokenTTL, refreshTokenTTL := getEnvPairTokens(env)

	return &Config{
		Security: &Security{
			PrivateKey:      privateKey,
			AccessTokenTTL:  accessTokenTTL,
			RefreshTokenTTL: refreshTokenTTL,
			Secret:          []byte(getEnvOrFatal(Secret)),
		},
		Server: &configs.Server{
			Host: getEnvDefault(HOST, "localhost"),
			Port: getEnvDefault(PORT, "44044"),
		},
		DB: &configs.DB{
			DSN: getEnvOrFatal(DSN),
		},
	}
}
