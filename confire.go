package confire

func Process(prefix string, spec interface{}, opts ...Option) error {
	return nil
}

func MustProcess(prefix string, spec interface{}, opts ...Option) {
	if err := Process(prefix, spec, opts...); err != nil {
		panic(err)
	}
}
