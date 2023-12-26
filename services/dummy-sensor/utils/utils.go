package utils

import (
	"os"
	"strconv"
	"time"
)

func GetEnvOrDefaultString(key string, defaultValue string) string {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		return value
	}
}

func GetEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		if x, e := strconv.Atoi(value); e != nil {
			return defaultValue
		} else {
			return x
		}
	}
}

func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		if v, e := strconv.ParseBool(value); e != nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func GetEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); len(value) == 0 {
		return defaultValue
	} else {
		if v, e := time.ParseDuration(value); e != nil {
			return defaultValue
		} else {
			return v
		}
	}
}
