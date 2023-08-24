package conf

type (
	Options struct {
		Env string
	}

	OptionFunc func(*Options)
)

func UseEnv(env string) OptionFunc {
	return func(opts *Options) {
		opts.Env = env
	}
}
