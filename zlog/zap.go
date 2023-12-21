package zlog

import (
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/uaxe/infra/zlog/zapx"
)

func Zap(directory, prefix, format, stacktraceKey, level,
	encoderLevel string, showLine bool, opts ...zapx.OptionFunc) (*zap.Logger, error) {
	if ok, _ := PathExists(directory); !ok {
		_ = os.Mkdir(directory, os.ModePerm)
	}
	fileRotate, err := zapx.NewFileRotateLogs(directory, opts...)
	if err != nil {
		return nil, err
	}
	z := zapx.NewZap(prefix, format, stacktraceKey, ZapLevel(level),
		ZapEncoderLevel(encoderLevel), fileRotate.GetWriteSyncer)
	cores := z.GetZapCores()
	logger := zap.New(zapcore.NewTee(cores...))
	if showLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger, nil
}

func ZapEncoderLevel(encoderLevel string) zapcore.LevelEncoder {
	switch encoderLevel {
	case "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func ZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
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