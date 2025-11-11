package chatconfig

import (
	"log"
	"path/filepath"

	"github.com/Krokozabra213/sso/configs"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

var (
	emptyValue        = ""
	op                = "chatconfig: "
	defaultHost       = "localhost"
	defaultPort       = "44045"
	defaultCtxTimeout = 5000
)

type Config struct {
	Security *Security
	Server   *configs.Server
	DB       *configs.DB
	Redis    *configs.RedisDB
}

type Security struct {
	Secret []byte
}

func Load(env string, test bool) *Config {
	var root string
	var err error

	if root, err = configs.FindProjectRoot(); err != nil {
		log.Fatalf("Project root not found: %v", err)
	}

	if test {
		envPath := filepath.Join(root, ".test.env")
		err = godotenv.Load(envPath)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatalf("%s%s", op, configs.ErrLoadConfig.Error())
	}

	return &Config{
		Security: &Security{
			Secret: []byte(getEnvOrFatal(Secret)),
		},
		Server: &configs.Server{
			Host:    getEnvDefault(HOST, defaultHost),
			Port:    getEnvDefault(PORT, defaultPort),
			TimeOut: getEnvDefaultInt(ContextTimeout, defaultCtxTimeout),
		},
		DB: &configs.DB{
			DSN: getEnvOrFatal(DSN),
		},
		Redis: &configs.RedisDB{
			Addr:  getEnvOrFatal(RedisAddr),
			Pass:  getEnvOrFatal(RedisPass),
			Cache: getEnvAtoiOrFatal(RedisCache),
		},
	}
}
