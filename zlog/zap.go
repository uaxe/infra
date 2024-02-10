package zlog

import (
	"strings"

	"github.com/uaxe/infra/zlog/zapx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap(directory, prefix, format, stacktraceKey, level,
	encoderLevel string, showLine bool, opts ...zapx.OptionFunc) (*zap.Logger, error) {
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

func ZapAndWriter(directory, prefix, format, stacktraceKey, level,
	encoderLevel string, showLine bool, opts ...zapx.OptionFunc) (*zap.Logger,
	func(level string) (zapcore.WriteSyncer, error), error) {
	fileRotate, err := zapx.NewFileRotateLogs(directory, opts...)
	if err != nil {
		return nil, nil, err
	}
	z := zapx.NewZap(prefix, format, stacktraceKey, ZapLevel(level),
		ZapEncoderLevel(encoderLevel), fileRotate.GetWriteSyncer)
	cores := z.GetZapCores()
	logger := zap.New(zapcore.NewTee(cores...))
	if showLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger, fileRotate.GetWriteSyncer, nil
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
