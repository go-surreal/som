package testutil

import (
	"os"
	"testing"
)

func SkipWithoutEnv(t *testing.T, key string) string {
	env := os.Getenv(key)
	if env == "" {
		t.Skip("set " + key + " env variable to run this test")
	}

	return env
}
