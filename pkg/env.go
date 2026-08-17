package pkg

import (
	"log"
	"os"
	"strconv"
)

func MustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("переменная окружения %s не задана", key)
	}
	return val
}

func EnvOr(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func EnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("переменная %s должна быть числом, получили: %s", key, val)
	}
	return n
}
