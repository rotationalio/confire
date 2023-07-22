package confire

import (
	"github.com/rotationalio/confire/defaults"
	"github.com/rotationalio/confire/env"
	"github.com/rotationalio/confire/validate"
)

// Process is the main entry point to configuring and validating a struct from defaults
// and the environment. Pass in a prefix for environment variables and a pointer to the
// configuration struct you want processed (as well as any options). The processor will
// first populate the struct with defaults, then load any values found in the
// environment, finally validating the struct based on struct tags and the validate
// interface. A ParseError or a ValidationError may be returned if not successful.
func Process(prefix string, spec interface{}, opts ...Option) (err error) {
	var opt *options
	if opt, err = makeOptions(opts...); err != nil {
		return err
	}

	if !opt.noDefaults {
		if err = defaults.Process(spec); err != nil {
			return err
		}
	}

	if !opt.noEnv {
		if err = env.Process(prefix, spec); err != nil {
			return err
		}
	}

	if !opt.noValidate {
		if err = validate.Validate(spec); err != nil {
			return err
		}
	}

	return nil
}

// MustProcess panics if processing the specification results in an error.
func MustProcess(prefix string, spec interface{}, opts ...Option) {
	if err := Process(prefix, spec, opts...); err != nil {
		panic(err)
	}
}
