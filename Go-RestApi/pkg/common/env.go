package common

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return valueInt
}
