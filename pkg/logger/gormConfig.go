package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	ormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type gormLogger struct {
	log zerolog.Logger
}

func GormConfig() *gormLogger {
	return &gormLogger{log.Logger}
}
func (l *gormLogger) LogMode(ormLogger.LogLevel) ormLogger.Interface {
	return l // do nothing here
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.InfoLevel {
		l.log.Info().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.WarnLevel {
		l.log.Warn().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.ErrorLevel {
		l.log.Error().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		l.log.Error().Time("time", time.Now()).Msg(err.Error())
		return
	}

	level := l.log.GetLevel()

	if level <= zerolog.ErrorLevel {
		elapsed := time.Since(begin)

		sql, rows := fc()
		if rows == -1 {
			l.log.Debug().Time("time", time.Now()).Msg(fmt.Sprintf("%s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			l.log.Debug().Time("time", time.Now()).Msg(fmt.Sprintf("%s\n[%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	}
}
