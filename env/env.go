package env

import (
	"os"
)

func Expand(s string, mapping func(string) string) string {
	return os.Expand(s, mapping)
}

func ExpandEnv(s string) string {
	return os.Expand(s, os.Getenv)
}

func Getenv(name string, values ...string) string {
	value, found := os.LookupEnv(name)
	if !found && len(values) > 0 {
		value = values[0]
	}
	return value
}

func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

func Clearenv() {
	os.Clearenv()
}

func Environ() []string {
	return os.Environ()
}

func IsShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

func IsAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
