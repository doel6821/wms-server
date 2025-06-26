package helpers

import (
	"fmt"
	"runtime"
	"wms-server/constants"
	"wms-server/helpers/models"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"

	// elastic "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	// elogrus "gopkg.in/sohlich/elogrus.v7"
)

// InitializeNewLogs ...
func InitializeNewLogs() *log.Logger {
	l := log.New()
	l.SetFormatter(&log.JSONFormatter{
		PrettyPrint: false,
	})
	l.SetReportCaller(true)

	size, _ := strconv.Atoi(GetEnv("LOG_SIZE"))
	bu, _ := strconv.Atoi(GetEnv("LOG_BACKUP"))
	age, _ := strconv.Atoi(GetEnv("LOG_AGE"))
	compress, _ := strconv.ParseBool(GetEnv("LOG_COMPRESS"))

	l.SetOutput(&lumberjack.Logger{
		Filename:   GetEnv("LOG_NAME"),
		MaxSize:    size,     // Max size in MB
		MaxBackups: bu,       // Max number of old log files to keep
		MaxAge:     age,      // Max age in days to keep a log file
		Compress:   compress, // Compress old log files
	})

	// l.Out = os.Stdout
	// file, err := os.OpenFile(GetEnv("LOG_NAME"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	l.Out = file
	// } else {
	// 	l.Warn("Failed to log to file, using default stderr")
	// }

	return l
}

// TraceIDHook ...
type TraceIDHook struct {
	TraceID string
}

// NewTraceIDHook ...
func NewTraceIDHook(traceID string) log.Hook {
	hook := TraceIDHook{
		TraceID: traceID,
	}
	return &hook
}

// Fire ...
func (hook *TraceIDHook) Fire(entry *log.Entry) error {
	entry.Data["TRACE_ID"] = entry.Context.Value(constants.TRANSACTION_ID)
	entry.Data["SERVICE_NAME"] = constants.SERVICE_NAME
	return nil
}

// Levels ...
func (hook *TraceIDHook) Levels() []log.Level {
	return log.AllLevels
}

// GetCaller ..
func GetCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// GettingResponseLog ...
func GettingResponseLog(g *gin.Context, req, res interface{}, duration time.Duration) log.Fields {
	logsField := models.LogModels{
		METHOD:   g.Request.Method,
		PATH:     g.Request.URL.Path,
		HEADER:   g.Request.Header,
		CLIENTIP: g.Request.RemoteAddr,
		REQUEST:  req,
		RESPONSE: res,
		DURATION: duration.Milliseconds(),
	}
	return structs.Map(logsField)
}

// GettingDetaultLog ...
func GettingDetaultLog(message interface{}) log.Fields {
	return log.Fields{
		"MESSAGE": message,
	}
}

// MyWriter ...
type MyWriter struct {
	log *log.Logger
}

// Printf ...
func (m *MyWriter) Printf(format string, v ...interface{}) {
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

	m.log.WithFields(logsSQL).Info("Query logs")
}

// NewMyWriter ...
func NewMyWriter(log *log.Logger) *MyWriter {
	return &MyWriter{log: log}
}
