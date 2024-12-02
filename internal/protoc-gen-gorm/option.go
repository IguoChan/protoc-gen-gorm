package protoc_gen_gorm

var (
	defaultOptions = &Options{
		version: "v0.0.0",
	}
)

type Options struct {
	version        string
	withGormOption bool
}

type Option func(*Options)

func WithVersion(version string) Option {
	return func(o *Options) {
		o.version = version
	}
}

func WithGormOption() Option {
	return func(o *Options) {
		o.withGormOption = true
	}
}
