package helpers

import (
	"context"
	"errors"
	"fmt"
	"wms-server/constants"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Options ...
type Options struct {
	Logger                    *log.Logger
	LogLevel                  gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	FileWithLineNumField      string
}

// Logger ...
type Logger struct {
	Options
}

// GormLogger ...
func GormLogger(opts Options) *Logger {
	l := &Logger{Options: opts}
	if l.LogLevel == 0 {
		l.LogLevel = gormlogger.Silent
	}

	return l
}

func formatGormLog(format string, v ...interface{}) map[string]interface{} {
	logsSQL := make(map[string]interface{})

	refFile := constants.REF_FILE
	duration := constants.DURATION
	row := constants.ROW_AFFECTED
	query := constants.QUERY

	logsSQL["WARNING"] = "-"

	if len(v) > 4 {
		duration++
		row++
		query++
		logsSQL["WARNING"] = fmt.Sprintf("%v", v[1])
	}

	logsSQL["REF_FILE"] = v[refFile]
	logsSQL["DURATION"] = fmt.Sprintf("%.3f ms", v[duration])
	logsSQL["ROW_AFFECTED"] = fmt.Sprintf("%v", v[row])
	logsSQL["QUERY"] = strings.ReplaceAll(fmt.Sprintf("%v", v[query]), `\"`, ``)

	return logsSQL
}

// LogMode ...
func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info ...
func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Logger.WithContext(ctx).WithFields(formatGormLog(s, args...)).Info()
	}
}

// Warn ...
func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Logger.WithContext(ctx).WithFields(formatGormLog(s, args...)).Warn()
	}
}

// Error ...
func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Logger.WithContext(ctx).WithFields(formatGormLog(s, args...)).Error()
	}
}

// Trace ...
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	fields := log.Fields{}
	// if l.FileWithLineNumField != "" {
	// 	fields[l.FileWithLineNumField] = utils.FileWithLineNum()
	// }

	fields["REF_FILE"] = utils.FileWithLineNum()

	sql, rows := fc()
	if rows == -1 {
		fields["ROW_AFFECTED"] = "-"
	} else {
		fields["ROW_AFFECTED"] = rows
	}

	elapsed := time.Since(begin)
	fields["DURATION"] = elapsed.Milliseconds()
	fields["QUERY"] = strings.ReplaceAll(sql, `"`, ``)
	fields["WARNING"] = "-"

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) && l.LogLevel >= gormlogger.Error:
		fields["WARNING"] = err
		l.Logger.WithContext(ctx).WithFields(fields).Error(constants.LOG_QUERY)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		fields["WARNING"] = "SLOW SQL"
		l.Logger.WithContext(ctx).WithFields(fields).Warn(constants.LOG_QUERY)
	case l.LogLevel == gormlogger.Info:
		l.Logger.WithContext(ctx).WithFields(fields).Info(constants.LOG_QUERY)
	}
}
