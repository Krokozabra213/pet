package ssonewconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Krokozabra213/sso/newconfigs"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultHttpHost               = "localhost"
	defaultHttpPort               = "44044"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
	defaultGRPCHost               = "localhost"
	defaultGRPCPort               = "8000"
	defaultGRPCRWTimeout          = 10 * time.Second
	defaultGRPCMaxHeaderMegabytes = 0
	defaultLoggerLevel            = 4
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultPrivateKeyPath         = "./secrets/private.pem"
)

type (
	AuthConfig struct {
		JWT          JWTConfig
		AppSecretKey []byte
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		PrivateKey      []byte
	}
	Logger struct {
		Level int `mapstructure:"level"`
	}
)

type Config struct {
	Auth   AuthConfig
	Logger Logger
	PG     newconfigs.Postgres
	Redis  newconfigs.Redis
	HTTP   newconfigs.HTTPConfig
	GRPC   newconfigs.GRPCConfig
}

func newCfg() Config {
	cfg := Config{
		Auth:   AuthConfig{},
		Logger: Logger{},
		PG:     newconfigs.Postgres{},
		Redis:  newconfigs.Redis{},
		HTTP:   newconfigs.HTTPConfig{},
		GRPC:   newconfigs.GRPCConfig{},
	}
	return cfg
}

func populateDefault() {
	viper.SetDefault("http.host", defaultHttpHost)
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.maxHeaderMegabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHttpRWTimeout)

	viper.SetDefault("grpc.host", defaultGRPCHost)
	viper.SetDefault("grpc.port", defaultGRPCPort)
	viper.SetDefault("grpc.maxHeaderMegabytes", defaultGRPCMaxHeaderMegabytes)
	viper.SetDefault("grpc.readTimeout", defaultGRPCRWTimeout)
	viper.SetDefault("grpc.writeTimeout", defaultGRPCRWTimeout)

	viper.SetDefault("logger.level", defaultLoggerLevel)

	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("auth.privateKeyPath", defaultPrivateKeyPath)
}

func PrivateKeyData(keyPath string) ([]byte, error) {
	if keyPath == "" {
		return []byte(nil), errors.New("absent private key path on env file")
	}
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return []byte(nil), fmt.Errorf("failed to read private key file: %v", err)
	}

	return keyData, nil
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
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("grpc", &cfg.GRPC); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	keyPath := viper.GetString("auth.privateKeyPath")
	keyPath = filepath.Join(root, keyPath)
	data, err := PrivateKeyData(keyPath)
	if err != nil {
		return err
	}
	cfg.Auth.JWT.PrivateKey = data

	return nil
}

func setFromEnv(envpath string, cfg *Config) error {
	err := godotenv.Load(envpath)
	if err != nil {
		return err
	}

	cfg.PG.DSN = os.Getenv("DSN")

	cfg.Auth.AppSecretKey = []byte(os.Getenv("APP_SECRET"))
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	cfg.Redis.Pass = os.Getenv("REDIS_PASS")

	cache := os.Getenv("REDIS_CACHE")
	cfg.Redis.Cache, err = strconv.Atoi(cache)
	if err != nil {
		return err
	}
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
