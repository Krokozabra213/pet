package configs

import (
	"fmt"
	"os"
	"path/filepath"
)

type DB struct {
	DSN string
}

type Server struct {
	Host    string
	Port    string
	TimeOut int
}

type RedisDB struct {
	Addr  string
	Pass  string
	Cache int
}

// func GetProjectRoot() string {
// 	// Получаем путь к текущему файлу
// 	_, filename, _, _ := runtime.Caller(0)

// 	// Поднимаемся до корня проекта (где лежит go.mod)
// 	root := filepath.Join(filepath.Dir(filename), "..", "..")
// 	absRoot, _ := filepath.Abs(root)

// 	return absRoot
// }

func FindProjectRoot() (string, error) {
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
