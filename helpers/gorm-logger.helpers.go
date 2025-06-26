package helpers

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// LoggerInterface ...
type (
	logger struct {
		SlowThreshold         time.Duration
		SourceField           string
		SkipErrRecordNotFound bool
		Debug                 bool
	}

	LoggerInterface interface {
		LogMode(gormlogger.LogLevel) gormlogger.Interface
		Info(ctx context.Context, s string, args ...interface{})
		Warn(ctx context.Context, s string, args ...interface{})
		Error(ctx context.Context, s string, args ...interface{})
		Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
	}
)

// GormLoggerNew ...
func GormLoggerNew() LoggerInterface {
	return &logger{
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

// LogMode ...
func (l *logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info ...
func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Infof(s, args...)
}

// Warn ...
func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Warnf(s, args...)
}

// Error ...
func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Errorf(s, args...)
}

// Trace ...
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := log.Fields{}
	formatError := "%s [%s]"
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[log.ErrorKey] = err
		log.WithContext(ctx).WithFields(fields).Errorf(formatError, sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.WithContext(ctx).WithFields(fields).Warnf(formatError, sql, elapsed)
		return
	}

	if l.Debug {
		log.WithContext(ctx).WithFields(fields).Debugf(formatError, sql, elapsed)
	}
}
