package envUtils

import (
	"fmt"
	"os"
)

// GetFromEnvOrDefault return value from environment variable, return default value if not defined
func GetFromEnvOrDefault(fromEnv, defaultValue string) (value string) {
	value = os.Getenv(fromEnv)
	if value == "" {
		value = defaultValue
	}
	return
}

// GetFromEnvOrError return value from environment variable, return error if not defined
func GetFromEnvOrError(envName string) (envValue string, err error) {
	envValue = os.Getenv(envName)
	if envValue == "" {
		return envValue, fmt.Errorf("env variable %s is not defined", envName)
	}
	return
}