package confire

import (
	"github.com/rotationalio/confire/defaults"
	"github.com/rotationalio/confire/env"
	"github.com/rotationalio/confire/validate"
)

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
