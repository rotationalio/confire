package confire

type Option func(opts *options) error

var NoDefaults = func(opts *options) error {
	opts.noDefaults = true
	return nil
}

var NoEnv = func(opts *options) error {
	opts.noEnv = true
	return nil
}

var NoValidate = func(opts *options) error {
	opts.noValidate = true
	return nil
}

type options struct {
	noDefaults bool
	noEnv      bool
	noValidate bool
}

func makeOptions(opts ...Option) (*options, error) {
	conf := &options{}
	for _, opt := range opts {
		if err := opt(conf); err != nil {
			return nil, err
		}
	}
	return conf, nil
}
