package emailnotify

type Option func(opt *option)

type option struct {
	EmailFrom       string
	AccountUserName string
	AccountPasswd   string
	SMTPHost        string
	SMTPPort        int
}

// WithEmailFrom
func WithEmailFrom(email string) Option {
	if !p.Match([]byte(email)) {
		panic(Err_Incorrect)
	}

	return func(opt *option) {
		opt.EmailFrom = email
	}
}

// WithAccountUserName
func WithAccountUserName(username string) Option {
	return func(opt *option) {
		opt.AccountUserName = username
	}
}

// WithAccountPasswd
func WithAccountPasswd(passwd string) Option {
	return func(opt *option) {
		opt.AccountPasswd = passwd
	}
}

// WithSMTPHost smtp
func WithSMTPHost(host string) Option {
	return func(opt *option) {
		opt.SMTPHost = host
	}
}

// WithSMTPPort smtp
func WithSMTPPort(port int) Option {
	return func(opt *option) {
		opt.SMTPPort = port
	}
}
