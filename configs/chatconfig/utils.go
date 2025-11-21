package chatconfig

import (
	"log"
	"os"
	"strconv"

	"github.com/Krokozabra213/sso/configs"
)

func getEnvAtoiOrFatal(key string) int {
	stringVal := getEnvOrFatal(key)
	intVal, err := strconv.Atoi(stringVal)
	if err != nil {
		log.Fatalf("%s%s", op, err.Error())
	}
	return intVal
}

func getEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != emptyValue {
		return value
	}
	return defaultValue
}

func getEnvDefaultInt(key string, defaultValue int) int {
	var intValue int
	var err error

	if value := os.Getenv(key); value != emptyValue {
		if intValue, err = strconv.Atoi(value); err != nil {
			log.Printf("getEnvDefaultInt config err")
			return defaultValue
		}
	}
	return intValue
}

func getEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == emptyValue {
		log.Fatalf("%s%s", op, configs.ErrEmptyValue.Error())
	}
	return value
}
