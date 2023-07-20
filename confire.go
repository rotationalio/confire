package confire

func Process(prefix string, spec interface{}, opts ...Option) error {
	return nil
}

// MustProcess panics if processing the specification results in an error.
func MustProcess(prefix string, spec interface{}, opts ...Option) {
	if err := Process(prefix, spec, opts...); err != nil {
		panic(err)
	}
}
