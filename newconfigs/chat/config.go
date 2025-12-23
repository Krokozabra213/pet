package chatnewconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Krokozabra213/sso/newconfigs"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultGRPCHost               = "localhost"
	defaultGRPCPort               = "44045"
	defaultGRPCRWTimeout          = 10 * time.Second
	defaultGRPCMaxHeaderMegabytes = 1
	defaultLoggerLevel            = 0
)

type (
	Security struct {
		AppSecretKey []byte
	}
)

type Config struct {
	PG       newconfigs.Postgres
	Redis    newconfigs.Redis
	GRPC     newconfigs.GRPCConfig
	Security Security
}

func newCfg() Config {
	cfg := Config{
		PG:       newconfigs.Postgres{},
		Redis:    newconfigs.Redis{},
		GRPC:     newconfigs.GRPCConfig{},
		Security: Security{},
	}
	return cfg
}

func populateDefault() {
	viper.SetDefault("grpc.host", defaultGRPCHost)
	viper.SetDefault("grpc.port", defaultGRPCPort)
	viper.SetDefault("grpc.maxHeaderMegabytes", defaultGRPCMaxHeaderMegabytes)
	viper.SetDefault("grpc.readTimeout", defaultGRPCRWTimeout)
	viper.SetDefault("grpc.writeTimeout", defaultGRPCRWTimeout)
	viper.SetDefault("logger.level", defaultLoggerLevel)
}

func parseConfigFile(configPath string) error {

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func Init(configfile, envfile string) (*Config, error) {

	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}
	configpath := filepath.Join(root, configfile)
	envpath := filepath.Join(root, envfile)

	populateDefault()

	if err := parseConfigFile(configpath); err != nil {
		return nil, err
	}

	cfg := newCfg()

	err = unmarshal(&cfg, root)
	if err != nil {
		return nil, err
	}

	err = setFromEnv(envpath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func unmarshal(cfg *Config, root string) error {

	if err := viper.UnmarshalKey("grpc", &cfg.GRPC); err != nil {
		return err
	}

	return nil
}

func setFromEnv(envpath string, cfg *Config) error {
	err := godotenv.Load(envpath)
	if err != nil {
		return err
	}

	cfg.PG.DSN = os.Getenv("DSN")

	cfg.Security.AppSecretKey = []byte(os.Getenv("APP_SECRET"))
	// cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	// cfg.Redis.Pass = os.Getenv("REDIS_PASS")

	// cache := os.Getenv("REDIS_CACHE")
	// cfg.Redis.Cache, err = strconv.Atoi(cache)
	// if err != nil {
	// 	return err
	// }

	cfg.PG.User = os.Getenv("POSTGRES_USER")
	cfg.PG.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.PG.DB = os.Getenv("POSTGRES_DB")
	cfg.PG.LocalPort = os.Getenv("POSTGRES_PORT")
	return nil
}

func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Проверяем, есть ли go.mod в текущей директории
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		// Поднимаемся на уровень выше
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Достигли корня файловой системы
			return "", fmt.Errorf("go.mod not found")
		}
		currentDir = parent
	}
}
