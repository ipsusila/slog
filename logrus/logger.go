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
	defaultTimestampFormat         = "2006/01/02 15:04:05 MST"
	fieldFormatter                 = "formatter"
	fieldTimestampFormat           = "timestampFormat"
	fieldDisableTimestamp          = "disbleTimestamp"
	fieldReportCaller              = "reportCaller"
	fieldFulltimeStamp             = "fullTimestamp"
	fieldMapper                    = "fieldMap"
	fieldDataKey                   = "dataKey"
	fieldPrettyPrint               = "prettyPrint"
	fieldDisableHTMLEscape         = "disableHTMLEscape"
	fieldForceColors               = "forceColors"
	fieldDisableColors             = "disableColors"
	fieldForceQuote                = "forceQuote"
	fieldDisableQuote              = "disableQuote"
	fieldEnvironmentOverrideColors = "environmentOverrideColors"
	fieldDisableSorting            = "disableSorting"
	fieldDisableLevelTruncation    = "disableLevelTruncation"
	fieldPadLevelText              = "padLevelText"
	fieldQuoteEmptyFields          = "quoteEmptyFields"
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
func New(w io.Writer, l slog.Level, op slog.Options) (slog.Logger, error) {
	ll, ok := toLogrusLevel(l)
	if !ok {
		return nil, fmt.Errorf("unknown logger level: %v", l)
	}

	// customize standard fields
	var fieldMap = log.FieldMap{
		log.FieldKeyMsg:         "@msg",
		log.FieldKeyLevel:       "@level",
		log.FieldKeyTime:        "@time",
		log.FieldKeyLogrusError: "@logrus_error",
		log.FieldKeyFunc:        "@func",
		log.FieldKeyFile:        "@file",
	}

	// get various options
	reportCaller := false
	var formatter log.Formatter = &log.TextFormatter{
		TimestampFormat: defaultTimestampFormat,
		FullTimestamp:   true,
		FieldMap:        fieldMap,
	}
	if len(op) != 0 {
		// override report caller options
		reportCaller = op.GetBool(fieldReportCaller, false)

		// customized field mapper
		fm := op.GetOptions(fieldMapper)
		for key, val := range fieldMap {
			fieldMap[key] = fm.GetString(string(key), val)
		}

		disableTimestamp := op.GetBool(fieldDisableTimestamp, false)

		// formatter options
		txtF := op.GetString(fieldFormatter, "text")
		switch txtF {
		case "json":
			formatter = &log.JSONFormatter{
				TimestampFormat:   op.GetString(fieldTimestampFormat, time.RFC3339),
				DataKey:           op.GetString(fieldDataKey, ""),
				PrettyPrint:       op.GetBool(fieldPrettyPrint, false),
				DisableHTMLEscape: op.GetBool(fieldDisableHTMLEscape, false),
				DisableTimestamp:  disableTimestamp,
				FieldMap:          fieldMap,
			}
		default:
			formatter = &log.TextFormatter{
				ForceColors:               op.GetBool(fieldForceColors, false),
				DisableColors:             op.GetBool(fieldDisableColors, false),
				ForceQuote:                op.GetBool(fieldForceQuote, false),
				DisableQuote:              op.GetBool(fieldDisableQuote, false),
				EnvironmentOverrideColors: op.GetBool(fieldEnvironmentOverrideColors, false),
				DisableSorting:            op.GetBool(fieldDisableSorting, false),
				DisableLevelTruncation:    op.GetBool(fieldDisableLevelTruncation, false),
				PadLevelText:              op.GetBool(fieldPadLevelText, false),
				QuoteEmptyFields:          op.GetBool(fieldQuoteEmptyFields, false),
				TimestampFormat:           op.GetString(fieldTimestampFormat, defaultTimestampFormat),
				FullTimestamp:             op.GetBool(fieldFulltimeStamp, true),
				DisableTimestamp:          disableTimestamp,
				FieldMap:                  fieldMap,
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
	return New(w, l, op)
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
