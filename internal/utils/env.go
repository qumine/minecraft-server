package utils

import (
	"os"
	"strconv"
	"strings"
)

// GetEnvString get key environment variable if exist otherwise return defalutValue
func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// GetEnvStringList get key environment variable as string list if exist otherwise return defalutValue
func GetEnvStringList(key string, defaultValue string) []string {
	return strings.Split(GetEnvString(key, defaultValue), ",")
}

// GetEnvInt get key environment variable if exist otherwise return defalutValue
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	int32Value, err := strconv.Atoi(value)
	if err != nil {
		return int32Value
	}
	return int32Value
}

// GetEnvBool get key environment variable if exist otherwise return defalutValue
func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return boolValue
	}
	return boolValue
}
