package ssoconfig

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Krokozabra213/sso/configs"
)

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

func getEnvPairTokens(env string) (int, int) {
	var accessTokenTTL int
	var refreshTokenTTL int
	if env == envLocal || env == envDev {
		accessTokenTTL = getEnvAtoiOrFatal(DevAccessTokenTTL)
		refreshTokenTTL = getEnvAtoiOrFatal(DevRefreshTokenTTL)
	} else {
		accessTokenTTL = getEnvAtoiOrFatal(ProdAccessTokenTTL)
		refreshTokenTTL = getEnvAtoiOrFatal(ProdRefreshTokenTTL)
	}
	return accessTokenTTL, refreshTokenTTL
}

func getEnvAtoiOrFatal(key string) int {
	stringVal := getEnvOrFatal(key)
	intVal, err := strconv.Atoi(stringVal)
	if err != nil {
		log.Fatalf("%s%s", op, err.Error())
	}
	return intVal
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == emptyValue {
		log.Fatalf("%s%s", op, configs.ErrEmptyValue.Error())
	}
	return value
}
