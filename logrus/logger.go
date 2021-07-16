package logrus

import (
	"fmt"
	"io"
	"time"

	"github.com/ipsusila/slog"
	log "github.com/sirupsen/logrus"
)

// Name of logrus logger
const Name = "logrus"

type logrusConstructor struct{}

const (
	fieldFormatter       = "formatter"
	fieldTimestampFormat = "timestampFormat"
	fieldReportCaller    = "reportCaller"
	fieldFulltimeStamp   = "fullTimestamp"
)

// logger without output, except for panic
type logrusLogger struct {
	log.Logger
	level slog.Level
}

func init() {
	slog.Register(Name, &logrusConstructor{})
}

func toLogrusLevel(l slog.Level) (log.Level, bool) {
	all := slog.Levels()
	maxIdx := -1
	for idx, lv := range all {
		if l.Has(lv) {
			maxIdx = idx
		}
	}
	if maxIdx >= 0 {
		ui32 := uint32(maxIdx) + uint32(log.PanicLevel)
		return log.Level(ui32), true
	}
	return 0, false
}

// New creates discard logger
func New(w io.Writer, l slog.Level, op *slog.Options) (slog.Logger, error) {
	ll, ok := toLogrusLevel(l)
	if !ok {
		return nil, fmt.Errorf("unknown logger level: %v", l)
	}

	// get various options
	reportCaller := false
	fullTimestamp := true
	var formatter log.Formatter = &log.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05 MST",
		FullTimestamp:   fullTimestamp,
	}
	if op != nil {
		// report caller here?
		reportCaller = op.GetBool(fieldReportCaller, false)
		txtF := op.GetString(fieldFormatter, "text")
		switch txtF {
		case "json":
			formatter = &log.JSONFormatter{
				TimestampFormat: op.GetString(fieldTimestampFormat, time.RFC3339),
			}
		default:
			formatter = &log.TextFormatter{
				TimestampFormat: op.GetString(fieldTimestampFormat, "2006/01/02 15:04:05 MST"),
				FullTimestamp:   op.GetBool(fieldTimestampFormat, true),
			}
		}
	}
	// end options

	lr := log.Logger{
		Out:          w,
		Formatter:    formatter,
		Level:        ll,
		ReportCaller: reportCaller,
	}
	lg := logrusLogger{
		Logger: lr,
		level:  l,
	}
	return &lg, nil
}

func (c *logrusConstructor) New(w io.Writer, l slog.Level) (slog.Logger, error) {
	return New(w, l, nil)
}
func (c *logrusConstructor) NewWithOptions(w io.Writer, l slog.Level, op slog.Options) (slog.Logger, error) {
	return New(w, l, &op)
}

func (l *logrusLogger) HasLevel(lv slog.Level) bool {
	return l.level.Has(lv)
}
func (l *logrusLogger) SetLevel(lv slog.Level) {
	if ll, ok := toLogrusLevel(lv); ok {
		l.level = lv
		l.Logger.SetLevel(ll)
	}
}

func (l *logrusLogger) Tracew(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Trace(msg)
}
func (l *logrusLogger) Debugw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Debug(msg)
}
func (l *logrusLogger) Printw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Print(msg)
}
func (l *logrusLogger) Infow(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Info(msg)
}
func (l *logrusLogger) Warnw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Warn(msg)
}
func (l *logrusLogger) Errorw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Error(msg)
}
func (l *logrusLogger) Fatalw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Fatal(msg)
}
func (l *logrusLogger) Panicw(msg string, keyVals ...interface{}) {
	fields := slog.FieldsToMap(keyVals)
	l.WithFields(log.Fields(fields)).Panic(msg)
}
