package zapx

import (
	"errors"
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

func (f *FileRotatelogs) levelPath(level string) string {
	p := fmt.Sprintf("%s%s", level, f.option.logExt)
	if f.option.prefix != "" {
		p = fmt.Sprintf("%s.%s", f.option.prefix, p)
	}
	lp := path.Join(f.option.directory, "%Y-%m-%d", p)
	if f.option.levelPathFn != nil {
		lp = f.option.levelPathFn(level)
	}
	return lp
}

func (f *FileRotatelogs) GetWriteSyncer(level string) (zapcore.WriteSyncer, error) {
	writer := make([]zapcore.WriteSyncer, 0, 2)
	if f.option.fileWriter {
		if ok, _ := PathExists(f.option.directory); !ok {
			_ = os.Mkdir(f.option.directory, os.ModePerm)
		}

		fileWriter, err := rotatelogs.New(
			f.levelPath(level),
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Hour*24*7),
			rotatelogs.WithRotationTime(time.Hour*24),
		)
		if err != nil {
			return nil, err
		}
		writer = append(writer, zapcore.AddSync(fileWriter))
	}
	if f.option.console || len(writer) == 0 {
		writer = append(writer, zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(writer...), nil
}

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("same name exists")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
