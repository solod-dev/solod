package os

import "solod.dev/so/c"

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// The returned string is a view into the static environment buffer;
// the caller must not modify or free it.
func Getenv(key string) string {
	ptr := getenv(key).(*byte)
	if ptr == nil {
		return ""
	}
	return c.String(ptr)
}

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
//
// The returned string is a view into the static environment buffer;
// the caller must not modify or free it.
func LookupEnv(key string) (string, bool) {
	ptr := getenv(key).(*byte)
	if ptr == nil {
		return "", false
	}
	return c.String(ptr), true
}

// Setenv sets the value of the environment variable named by the key.
// It returns an error, if any.
func Setenv(key, value string) error {
	if setenv(key, value, 1) != 0 {
		return mapError()
	}
	return nil
}

// Unsetenv unsets a single environment variable.
func Unsetenv(key string) error {
	if unsetenv(key) != 0 {
		return mapError()
	}
	return nil
}
