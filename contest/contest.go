package contest

import "os"

// Env allows the user to set environment variables for tests and restore the original
// environment variables after the test is complete.
type Env map[string]string

type Cleanup func()

// Sets the environment variables from the Vars map, tracking the original values and
// returning a cleanup function that will restore the environment to its original state
// to ensure subsequent tests are not affected.
//
// Specify a list of keys to set, or no keys to set all environment variables.
//
// Usage: t.Cleanup(env.Set())
func (e Env) Set(keys ...string) Cleanup {
	orig := make(map[string]*string)

	if len(keys) > 0 {
		for _, key := range keys {
			// Skip if the key is not in the Env map
			val, ok := e[key]
			if !ok {
				continue
			}

			if oval, ok := os.LookupEnv(key); ok {
				orig[key] = &oval
			} else {
				orig[key] = nil
			}
			os.Setenv(key, val)
		}
	} else {
		for key, val := range e {
			if oval, ok := os.LookupEnv(key); ok {
				orig[key] = &oval
			} else {
				orig[key] = nil
			}
			os.Setenv(key, val)
		}
	}

	return func() {
		for key, val := range orig {
			if val != nil {
				os.Setenv(key, *val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}

// Clear the environment variables defined by the Vars map, tracking the original
// values and returning a cleanup function that will restore the environment to its
// original state to ensure subsequent tests are not affected.
//
// Specify a list of keys to clear, or no keys to clear all environment variables.
//
// Usage: t.Cleanup(env.Clear())
func (e Env) Clear(keys ...string) Cleanup {
	orig := make(map[string]*string)
	if len(keys) > 0 {
		for _, key := range keys {
			// Skip if the key is not in the Env map
			if _, ok := e[key]; !ok {
				continue
			}

			if oval, ok := os.LookupEnv(key); ok {
				orig[key] = &oval
			} else {
				orig[key] = nil
			}
			os.Unsetenv(key)
		}
	} else {
		for key := range e {
			if oval, ok := os.LookupEnv(key); ok {
				orig[key] = &oval
			} else {
				orig[key] = nil
			}
			os.Unsetenv(key)
		}
	}

	return func() {
		for key, val := range orig {
			if val != nil {
				os.Setenv(key, *val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
