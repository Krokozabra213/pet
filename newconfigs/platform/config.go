package platformconfig

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
	defaultHttpHost               = "localhost"
	defaultHttpPort               = "44050"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1

	defaultLimiterRPS   = 10
	defaultLimiterBurst = 2
	defaultLimiterTTL   = 10 * time.Minute

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	AppConfig struct {
		AppSecretKey []byte
		Environment  string
	}

	MongoConfig struct {
		URI      string
		Host     string
		Port     string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	FileStorageConfig struct {
		Endpoint  string
		Bucket    string
		AccessKey string
		SecretKey string
	}
)

type Config struct {
	App         AppConfig
	Mongo       MongoConfig
	FileStorage FileStorageConfig
	HTTP        newconfigs.HTTPConfig
	Limiter     newconfigs.LimiterConfig
}

func newCfg() Config {
	cfg := Config{
		App:         AppConfig{},
		Mongo:       MongoConfig{},
		HTTP:        newconfigs.HTTPConfig{},
		FileStorage: FileStorageConfig{},
		Limiter:     newconfigs.LimiterConfig{},
	}
	return cfg
}

func populateDefault() {
	viper.SetDefault("http.host", defaultHttpHost)
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.maxHeaderMegabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHttpRWTimeout)

	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
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

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}

	return nil
}

func setFromEnv(envpath string, cfg *Config) error {
	err := godotenv.Load(envpath)
	if err != nil {
		return err
	}

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.Host = os.Getenv("MONGO_HOST")
	cfg.Mongo.Port = os.Getenv("MONGO_PORT")
	cfg.Mongo.User = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.Mongo.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	cfg.App.Environment = os.Getenv("APP_ENV")
	cfg.App.AppSecretKey = []byte(os.Getenv("SECRET_KEY"))

	cfg.HTTP.Host = os.Getenv("HOST")

	cfg.FileStorage.Endpoint = os.Getenv("FILESTORAGE_ENDPOINT")
	cfg.FileStorage.Bucket = os.Getenv("FILESTORAGE_BUCKET")
	cfg.FileStorage.AccessKey = os.Getenv("FILESTORAGE_ACCESS_KEY")
	cfg.FileStorage.SecretKey = os.Getenv("FILESTORAGE_SECRET_KEY")

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
