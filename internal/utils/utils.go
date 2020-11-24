package utils

import "os"

// GetEnv get key environment variable if exist otherwise return defalutValue
func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
