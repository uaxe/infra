package zapx

import (
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

type FileRotateOption struct {
	logExt      string
	director    string
	prefix      string
	maxAge      time.Duration
	console     bool
	levelPathFn func(level string) string
}

type fileRotatelogs struct {
	option FileRotateOption
}

type OptionFunc func(opt *FileRotateOption) error

func SetMaxAge(maxAge time.Duration) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.maxAge = maxAge
		return nil
	}
}

func SetConsole(console bool) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.console = console
		return nil
	}
}

func SetPrefix(prefix string) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.prefix = prefix
		return nil
	}
}

func SetLogExt(logExt string) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.logExt = logExt
		return nil
	}
}

func SetLevelPath(f func(level string) string) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.levelPathFn = f
		return nil
	}
}

func NewFileRotateLogs(director string, opts ...OptionFunc) (*fileRotatelogs, error) {
	opt := defaultOption(director)
	for i := range opts {
		if err := opts[i](opt); err != nil {
			return nil, err
		}
	}
	f := &fileRotatelogs{
		option: *opt,
	}
	return f, nil
}

func defaultOption(director string) *FileRotateOption {
	return &FileRotateOption{
		logExt:   ".log",
		director: director,
		maxAge:   time.Duration(7),
		console:  false,
	}
}

func (self *fileRotatelogs) levelPath(level string) string {
	p := fmt.Sprintf("%s%s", level, self.option.logExt)
	if self.option.prefix != "" {
		p = fmt.Sprintf("%s.%s", self.option.prefix, p)
	}
	lp := path.Join(self.option.director, "%Y-%m-%d", p)
	if self.option.levelPathFn != nil {
		lp = self.option.levelPathFn(level)
	}
	return lp
}

func (self *fileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	fileWriter, err := rotatelogs.New(
		self.levelPath(level),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(self.option.maxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if self.option.console {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
