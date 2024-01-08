package utils

import (
	"os"
)

func Getenv(name string, values ...string) string {
	value, found := os.LookupEnv(name)
	if !found && len(values) > 0 {
		value = values[0]
	}
	return value
}
