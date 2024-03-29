package zapx

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Zap struct {
	prefix, format, stacktraceKey string
	level                         zapcore.Level
	encoder                       zapcore.LevelEncoder
	writeSyncer                   func(level string) (zapcore.WriteSyncer, error)
}

func NewZap(prefix, format, stacktraceKey string,
	level zapcore.Level, encoder zapcore.LevelEncoder,
	writeSyncer func(level string) (zapcore.WriteSyncer, error)) *Zap {
	z := &Zap{
		prefix:        prefix,
		format:        format,
		stacktraceKey: stacktraceKey,
		level:         level,
		encoder:       encoder,
		writeSyncer:   writeSyncer,
	}
	return z
}

func (z *Zap) GetEncoder() zapcore.Encoder {
	if z.format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

func (z *Zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  z.stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    z.encoder,
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (z *Zap) GetEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := z.writeSyncer(l.String())
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}
	return zapcore.NewCore(z.GetEncoder(), writer, level)
}

func (z *Zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(z.prefix + " " + t.Format("2006/01/02 - 15:04:05.000"))
}

func (z *Zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, zapcore.FatalLevel)
	for level := z.level; level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	return cores
}

func (z *Zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
