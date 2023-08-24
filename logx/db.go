package logx

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"runtime"
	"strings"
	"time"
)

type dbLog struct {
	logger.Config
	Log *zap.Logger
}

func NewDbLog(log *zap.Logger, config logger.Config) logger.Interface {
	return &dbLog{Log: log, Config: config}
}
func (db *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	db.LogLevel = level
	return db
}
func (db *dbLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if db.LogLevel < logger.Info {
		return
	}
	db.Log.Info("db info:", zap.Any("data", data))
}

func (db *dbLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if db.LogLevel < logger.Warn {
		return
	}
	db.Log.Warn("db ware:"+msg, zap.Any("data", data))
}

func (db *dbLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if db.LogLevel < logger.Error {
		return
	}
	db.Log.Error("db err:"+msg, zap.Any("data", data))
}

func (db *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if db.LogLevel <= logger.Silent {
		return
	}
	pc, file, line, _ := runtime.Caller(3)
	funcName := runtime.FuncForPC(pc).Name()
	elapsed := time.Since(begin)
	switch {
	case err != nil && db.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !db.IgnoreRecordNotFoundError):
		sql, rows := fc()
		sql = formatSql(sql)
		db.Log.Error("db trace err:", zap.String("file", file), zap.Int("line", line), zap.String("func", funcName), zap.Error(err), zap.Float64("ms", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("sql", sql))
	case elapsed > db.SlowThreshold && db.SlowThreshold != 0 && db.LogLevel >= logger.Warn:
		sql, rows := fc()
		sql = formatSql(sql)
		db.Log.Info("db trace ware:", zap.String("file", file), zap.Int("line", line), zap.String("func", funcName), zap.Duration("SLOW SQL>=", db.SlowThreshold), zap.Int64("rows", rows), zap.String("sql", sql))
	case db.LogLevel == logger.Info:
		sql, rows := fc()
		sql = formatSql(sql)
		db.Log.Info("db trace info:", zap.String("file", file), zap.Int("line", line), zap.String("func", funcName), zap.Float64("ms", float64(elapsed.Nanoseconds())/1e6), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func formatSql(sql string) string {
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	sql = strings.ReplaceAll(sql, "\\", "")
	return sql
}
