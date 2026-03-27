package contest_test

import (
	"os"
	"testing"

	"go.rtnl.ai/confire/assert"
	"go.rtnl.ai/confire/contest"
)

func TestEnv(t *testing.T) {
	// Set some environment variables for the "original" environment
	// Debug, Port, and User are values set in the "original" environment
	// Rate and Timeout are not set in the "original" environment
	//
	// NOTE: CONFIRE_FOO is not managed by the contest.Env as part of the tests.
	env := contest.Env{
		"CONFIRE_DEBUG":   "true",
		"CONFIRE_PORT":    "8888",
		"CONFIRE_RATE":    "0.25",
		"CONFIRE_USER":    "werebear",
		"CONFIRE_TIMEOUT": "5m",
	}

	// Ensure the "original" environment variables are as expected
	for key := range env {
		// Ensure the environment variable is not already set prior to testing
		if _, ok := os.LookupEnv(key); ok {
			t.Fatalf("environment variable %s is already set", key)
		}

		switch key {
		case "CONFIRE_DEBUG", "CONFIRE_PORT", "CONFIRE_USER":
			os.Setenv(key, "1")
		case "CONFIRE_RATE", "CONFIRE_TIMEOUT", "CONFIRE_FOO":
			os.Unsetenv(key)
		default:
			t.Fatalf("unhandled environment variable: %s", key)
		}
	}

	assertIsOriginal := func(t *testing.T) {
		t.Helper()
		for key := range env {
			switch key {
			case "CONFIRE_DEBUG", "CONFIRE_PORT", "CONFIRE_USER":
				assert.EnvEquals(t, key, "1")
			case "CONFIRE_RATE", "CONFIRE_TIMEOUT", "CONFIRE_FOO":
				assert.EnvUnset(t, key)
			default:
				t.Fatalf("unhandled environment variable: %s", key)
			}
		}
	}

	t.Run("Set", func(t *testing.T) {
		cleanup := env.Set()
		for key, val := range env {
			assert.EnvEquals(t, key, val)
		}

		cleanup()
		assertIsOriginal(t)
	})

	t.Run("SetKeys", func(t *testing.T) {
		keys := []string{"CONFIRE_DEBUG", "CONFIRE_RATE", "CONFIRE_FOO", "CONFIRE_USER"}
		cleanup := env.Set(keys...)

		for _, key := range keys {
			if val, ok := env[key]; ok {
				assert.EnvEquals(t, key, val)
			} else {
				assert.EnvUnset(t, key)
			}
		}

		// Omitted keys should remain unchanged
		assert.EnvEquals(t, "CONFIRE_PORT", "1")
		assert.EnvUnset(t, "CONFIRE_TIMEOUT")

		cleanup()
		assertIsOriginal(t)
	})

	t.Run("Clear", func(t *testing.T) {
		cleanup := env.Clear()
		for key := range env {
			assert.EnvUnset(t, key)
		}

		cleanup()
		assertIsOriginal(t)
	})

	t.Run("ClearKeys", func(t *testing.T) {
		keys := []string{"CONFIRE_DEBUG", "CONFIRE_RATE", "CONFIRE_FOO", "CONFIRE_USER"}
		cleanup := env.Clear(keys...)

		for _, key := range keys {
			assert.EnvUnset(t, key)
		}

		// Omitted keys should remain unchanged
		assert.EnvEquals(t, "CONFIRE_PORT", "1")
		assert.EnvUnset(t, "CONFIRE_TIMEOUT")

		cleanup()
		assertIsOriginal(t)
	})
}
