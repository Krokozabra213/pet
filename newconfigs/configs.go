package newconfigs

import (
	"errors"
	"time"
)

var (
	ErrLoadConfig = errors.New("error load config")
	ErrEmptyValue = errors.New("error empty value")
)

type (
	Postgres struct {
		DSN string
	}

	Redis struct {
		Addr  string
		Pass  string
		Cache int
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	GRPCConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LoggerConfig struct {
		AppSecretKey int `mapstructure:"level"`
	}
)
