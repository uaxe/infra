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
	directory   string
	prefix      string
	maxAge      time.Duration
	console     bool
	fileWriter  bool
	levelPathFn func(level string) string
}

type FileRotatelogs struct {
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

func SetFileWriter(fileWriter bool) OptionFunc {
	return func(opt *FileRotateOption) error {
		opt.fileWriter = fileWriter
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

func NewFileRotateLogs(director string, opts ...OptionFunc) (*FileRotatelogs, error) {
	opt := defaultOption(director)
	for i := range opts {
		if err := opts[i](opt); err != nil {
			return nil, err
		}
	}
	f := &FileRotatelogs{
		option: *opt,
	}
	return f, nil
}

func defaultOption(directory string) *FileRotateOption {
	return &FileRotateOption{
		logExt:    ".log",
		directory: directory,
		maxAge:    time.Duration(7),
		console:   false,
	}
}

func (self *FileRotatelogs) levelPath(level string) string {
	p := fmt.Sprintf("%s%s", level, self.option.logExt)
	if self.option.prefix != "" {
		p = fmt.Sprintf("%s.%s", self.option.prefix, p)
	}
	lp := path.Join(self.option.directory, "%Y-%m-%d", p)
	if self.option.levelPathFn != nil {
		lp = self.option.levelPathFn(level)
	}
	return lp
}

func (self *FileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	writer := make([]zapcore.WriteSyncer, 0, 2)
	if self.option.fileWriter {
		fileWriter, err := rotatelogs.New(
			self.levelPath(level),
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Duration(self.option.maxAge)*24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour*24),
		)
		if err != nil {
			return nil, err
		}
		writer = append(writer, zapcore.AddSync(fileWriter))
	}
	if self.option.console || len(writer) == 0 {
		writer = append(writer, zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(writer...), nil
}
