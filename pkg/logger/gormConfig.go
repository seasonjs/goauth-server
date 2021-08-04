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

//// Level level of a logging
//type Level uint8
//
//const (
//	// Debug logging level
//	Debug Level = 1
//	// Info logging level
//	Info Level = 2
//	// Warn logging level
//	Warn Level = 3
//	// Error logging level
//	Error Level = 4
//	// Fatal logging level
//	Fatal Level = 5
//)
//
//// Logger interface
//type Logger interface {
//	Fatal(s string, args ...interface{})
//	ErrFatal(err error)
//	Error(s string, args ...interface{})
//	ErrError(err error)
//	Warn(s string, args ...interface{})
//	ErrWarn(err error)
//	Info(s string, args ...interface{})
//	ErrInfo(err error)
//	Debug(s string, args ...interface{})
//	ErrDebug(err error)
//	CreateGormLogger() logger.Interface
//	Write(p []byte) (n int, err error)
//	GetLevel() Level
//}

// New logger constructor, creates a default logger instance.
//// Uses os.Stdout if the writer parameter is nil.
//func New(level Level, writer io.Writer) Logger {
//	if writer == nil {
//		writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
//	}
//
//	zerolog.SetGlobalLevel(getZerologLevel(level))
//	zLog := zerolog.New(writer).With().Stack().Logger()
//
//	return &DefaultLogger{
//		log:   &zLog,
//		level: level,
//	}
//}
//
//func getZerologLevel(level Level) zerolog.Level {
//	switch level {
//	case Fatal:
//		return zerolog.PanicLevel
//	case Error:
//		return zerolog.ErrorLevel
//	case Warn:
//		return zerolog.WarnLevel
//	case Info:
//		return zerolog.InfoLevel
//	default:
//		return zerolog.DebugLevel
//	}
//}
//
//// DefaultLogger default logging controller
//type DefaultLogger struct {
//	log   *zerolog.Logger
//	level Level
//}
//
//// Fatal log something at fatal level. This will panic!
//func (l *DefaultLogger) Fatal(s string, args ...interface{}) {
//	l.log.Panic().Time("time", time.Now()).Msg(fmt.Sprintf(s, args...))
//}
//
//// ErrFatal log an error at fatal level. This will panic!
//func (l *DefaultLogger) ErrFatal(err error) {
//	l.log.Panic().Time("time", time.Now()).Msg(err.Error())
//}
//
//// Error log something at error level
//func (l *DefaultLogger) Error(s string, args ...interface{}) {
//	if l.level <= Error {
//		l.log.Error().Time("time", time.Now()).Msg(fmt.Sprintf(s, args...))
//	}
//}
//
//// ErrError log an error at error level
//func (l *DefaultLogger) ErrError(err error) {
//	if l.level <= Error {
//		l.log.Error().Time("time", time.Now()).Msg(err.Error())
//	}
//}
//
//// Warn log something at warning level
//func (l *DefaultLogger) Warn(s string, args ...interface{}) {
//	if l.level <= Warn {
//		l.log.Warn().Time("time", time.Now()).Msg(fmt.Sprintf(s, args...))
//	}
//}
//
//// ErrWarn log an error at warning level
//func (l *DefaultLogger) ErrWarn(err error) {
//	if l.level <= Warn {
//		l.log.Warn().Time("time", time.Now()).Msg(err.Error())
//	}
//}
//
//// Info log something at info level
//func (l *DefaultLogger) Info(s string, args ...interface{}) {
//	if l.level <= Info {
//		l.log.Info().Time("time", time.Now()).Msg(fmt.Sprintf(s, args...))
//	}
//}
//
//// ErrInfo log an error at info level
//func (l *DefaultLogger) ErrInfo(err error) {
//	if l.level <= Info {
//		l.log.Info().Time("time", time.Now()).Msg(err.Error())
//	}
//}
//
//// Debug log something at debug level
//func (l *DefaultLogger) Debug(s string, args ...interface{}) {
//	if l.level <= Debug {
//		l.log.Debug().Time("time", time.Now()).Msg(fmt.Sprintf(s, args...))
//	}
//}
//
//// ErrDebug log an error at debug level
//func (l *DefaultLogger) ErrDebug(err error) {
//	if l.level <= Debug {
//		l.log.Debug().Time("time", time.Now()).Msg(err.Error())
//	}
//}
//
//// CreateGormLogger creates a mock of this logger for usage with gorm
//func (l *DefaultLogger) CreateGormLogger() logger.Interface {
//	return &gormLogger{
//		log: l,
//	}
//}
//
//// Write data to the logger (info level is used)
//func (l *DefaultLogger) Write(p []byte) (n int, err error) {
//	l.log.Info().Time("time", time.Now()).Msg(string(p))
//	return len(p), nil
//}
//
//// GetLevel returns the current error level
//func (l *DefaultLogger) GetLevel() Level {
//	return l.level
//}

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
