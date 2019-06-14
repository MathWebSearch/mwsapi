package utils

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Getenv gets an environment variable as a string with a fallback
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// GetenvInt is the same as Getenv, except that it turns values into int64
// it panics iff the variable exists but is not a valid int64
func GetenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		err = errors.Wrapf(err, "Unable to parse %q environment variable", key)
		panic(err)
	}
	return int(i)
}
