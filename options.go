package confire

type Option func(opts *options) error

type options struct{}
