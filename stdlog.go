package slog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Name of the standard logger
const StdLoggerName = "stdlog"

const (
	defaultTimestampFormat = "2006/01/02 15:04:05 MST"
	fieldTimestampFormat   = "timestampFormat"
	fieldDisableColor      = "disableColor"
)

type stdLoggerConstructor struct{}

// color mapper
var colorMapper = map[Level]func(string, ...interface{}) string{
	PanicLevel: color.HiRedString,
	FatalLevel: color.HiMagentaString,
	ErrorLevel: color.RedString,
	WarnLevel:  color.YellowString,
	InfoLevel:  color.GreenString,
	DebugLevel: color.BlueString,
	TraceLevel: color.CyanString,
}

// logger without output, except for panic
type stdLogger struct {
	LevelLoggerBase
	mu           sync.Mutex
	out          io.Writer
	buf          bytes.Buffer
	prefixes     map[Level]string
	tsFormat     string
	disableColor bool
}

func init() {
	Register(StdLoggerName, &stdLoggerConstructor{})
}

func (c *stdLoggerConstructor) New(w io.Writer, l Level) (Logger, error) {
	return NewStdLogger(w, l, nil)
}

// create logger with options
func (c *stdLoggerConstructor) NewWithOptions(w io.Writer, l Level, op Options) (Logger, error) {
	return NewStdLogger(w, l, op)
}

// NewStdLogger creates new logger with given parameters
func NewStdLogger(w io.Writer, l Level, op Options) (Logger, error) {
	bl := NewLevelLoggerBase(l)
	sl := stdLogger{
		LevelLoggerBase: *bl,
		out:             w,
		prefixes:        make(map[Level]string),
		tsFormat:        defaultTimestampFormat,
		disableColor:    false,
	}

	// customized options
	if len(op) != 0 {
		sl.tsFormat = op.GetString(fieldTimestampFormat, defaultTimestampFormat)
		sl.disableColor = op.GetBool(fieldDisableColor, false)
	}

	// create logger for each level
	all := Levels()
	for _, lv := range all {
		prefix := LevelFixedString(lv)
		if prefix != "" {
			prefix += " "
		}

		if sl.disableColor {
			if fn, ok := colorMapper[lv]; ok {
				prefix = fn(prefix)
			}
		}
		sl.prefixes[lv] = prefix
	}

	return &sl, nil
}

func (sl *stdLogger) writeHeader(prefix string) {
	sl.buf.WriteString(prefix)
	sl.buf.WriteRune('[')
	sl.buf.WriteString(time.Now().Format(sl.tsFormat))
	sl.buf.WriteString("] ")
}

func (sl *stdLogger) outputFields(lv Level, msg string, keyVals []interface{}, seps ...rune) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sep := '='
	if len(seps) > 0 {
		sep = seps[0]
	}

	prefix, ok := sl.prefixes[lv]
	if !ok {
		prefix = "OTHER"
	}
	sl.writeHeader(prefix)
	sl.buf.WriteString(msg)
	sl.buf.WriteRune('\t')

	// write fields
	n := len(keyVals)
	if n > 0 {
		// number of pair
		nkv := (n + 1) / 2

		// number of args is odd
		if n%2 != 0 {
			n--
		}

		// format field and values
		j := 0
		for i := 0; i < n; i += 2 {
			field, _ := AsString(keyVals[i])
			if !sl.disableColor {
				if fn, ok := colorMapper[lv]; ok {
					field = fn(field)
				}
			}
			sl.buf.WriteString(field)
			sl.buf.WriteRune(sep)
			sl.buf.WriteString(AsStringQ(keyVals[i+1]))
			sl.buf.WriteRune(' ')
			j++
		}

		// number of args is odd
		if n != len(keyVals) {
			field := fmt.Sprintf("%s-%02d", UnknownFieldName, nkv)
			if !sl.disableColor {
				if fn, ok := colorMapper[lv]; ok {
					field = fn(field)
				}
			}
			sl.buf.WriteString(field)
			sl.buf.WriteRune(sep)
			sl.buf.WriteString(AsStringQ(keyVals[n]))
		}
	}
	// write LF
	n = sl.buf.Len()
	if n == 0 || sl.buf.Bytes()[n-1] != '\n' {
		sl.buf.WriteRune('\n')
	}

	// copy content
	io.Copy(sl.out, &sl.buf)
}

func (sl *stdLogger) output(lv Level, str string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	prefix, ok := sl.prefixes[lv]
	if !ok {
		prefix = "OTHER"
	}
	sl.writeHeader(prefix)
	sl.buf.WriteString(str)

	// write LF
	n := sl.buf.Len()
	if n == 0 || sl.buf.Bytes()[n-1] != '\n' {
		sl.buf.WriteRune('\n')
	}

	io.Copy(sl.out, &sl.buf)
}

func (sl *stdLogger) Trace(args ...interface{}) {
	if sl.HasLevel(TraceLevel) {
		sl.output(TraceLevel, fmt.Sprint(args...))
	}
}
func (sl *stdLogger) Debug(args ...interface{}) {
	if sl.HasLevel(DebugLevel) {
		sl.output(DebugLevel, fmt.Sprint(args...))
	}
}
func (sl *stdLogger) Print(args ...interface{}) {
	sl.Info(args...)
}
func (sl *stdLogger) Info(args ...interface{}) {
	if sl.HasLevel(InfoLevel) {
		sl.output(InfoLevel, fmt.Sprint(args...))
	}
}
func (sl *stdLogger) Warn(args ...interface{}) {
	if sl.HasLevel(WarnLevel) {
		sl.output(WarnLevel, fmt.Sprint(args...))
	}
}
func (sl *stdLogger) Error(args ...interface{}) {
	if sl.HasLevel(ErrorLevel) {
		sl.output(ErrorLevel, fmt.Sprint(args...))
	}
}
func (sl *stdLogger) Fatal(args ...interface{}) {
	if sl.HasLevel(FatalLevel) {
		sl.output(FatalLevel, fmt.Sprint(args...))
	}
	os.Exit(1)

}
func (sl *stdLogger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	if sl.HasLevel(PanicLevel) {
		sl.output(PanicLevel, s)
	}
	panic(s)
}

func (sl *stdLogger) Traceln(args ...interface{}) {
	if sl.HasLevel(TraceLevel) {
		sl.output(TraceLevel, fmt.Sprintln(args...))
	}
}
func (sl *stdLogger) Debugln(args ...interface{}) {
	if sl.HasLevel(DebugLevel) {
		sl.output(DebugLevel, fmt.Sprintln(args...))
	}
}
func (sl *stdLogger) Println(args ...interface{}) {
	sl.Infoln(args...)
}
func (sl *stdLogger) Infoln(args ...interface{}) {
	if sl.HasLevel(InfoLevel) {
		sl.output(InfoLevel, fmt.Sprintln(args...))
	}
}
func (sl *stdLogger) Warnln(args ...interface{}) {
	if sl.HasLevel(WarnLevel) {
		sl.output(WarnLevel, fmt.Sprintln(args...))
	}
}
func (sl *stdLogger) Errorln(args ...interface{}) {
	if sl.HasLevel(ErrorLevel) {
		sl.output(ErrorLevel, fmt.Sprintln(args...))
	}
}
func (sl *stdLogger) Fatalln(args ...interface{}) {
	if sl.HasLevel(FatalLevel) {
		sl.output(FatalLevel, fmt.Sprintln(args...))
	}
	os.Exit(1)
}
func (sl *stdLogger) Panicln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	if sl.HasLevel(PanicLevel) {
		sl.output(PanicLevel, s)
	}
	panic(s)
}

func (sl *stdLogger) Tracef(format string, args ...interface{}) {
	if sl.HasLevel(TraceLevel) {
		sl.output(TraceLevel, fmt.Sprintf(format, args...))
	}
}
func (sl *stdLogger) Debugf(format string, args ...interface{}) {
	if sl.HasLevel(DebugLevel) {
		sl.output(DebugLevel, fmt.Sprintf(format, args...))
	}
}
func (sl *stdLogger) Printf(format string, args ...interface{}) {
	sl.Infof(format, args...)
}
func (sl *stdLogger) Infof(format string, args ...interface{}) {
	if sl.HasLevel(InfoLevel) {
		sl.output(InfoLevel, fmt.Sprintf(format, args...))
	}
}
func (sl *stdLogger) Warnf(format string, args ...interface{}) {
	if sl.HasLevel(WarnLevel) {
		sl.output(WarnLevel, fmt.Sprintf(format, args...))
	}
}
func (sl *stdLogger) Errorf(format string, args ...interface{}) {
	if sl.HasLevel(ErrorLevel) {
		sl.output(ErrorLevel, fmt.Sprintf(format, args...))
	}
}
func (sl *stdLogger) Fatalf(format string, args ...interface{}) {
	if sl.HasLevel(FatalLevel) {
		sl.output(FatalLevel, fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}
func (sl *stdLogger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if sl.HasLevel(PanicLevel) {
		sl.output(PanicLevel, s)
	}
	panic(s)
}

// with fields

func (sl *stdLogger) Tracew(msg string, keyVals ...interface{}) {
	if sl.HasLevel(TraceLevel) {
		sl.outputFields(TraceLevel, msg, keyVals)
	}
}
func (sl *stdLogger) Debugw(msg string, keyVals ...interface{}) {
	if sl.HasLevel(DebugLevel) {
		sl.outputFields(DebugLevel, msg, keyVals)
	}
}
func (sl *stdLogger) Printw(msg string, keyVals ...interface{}) {
	sl.Infow(msg, keyVals...)
}
func (sl *stdLogger) Infow(msg string, keyVals ...interface{}) {
	if sl.HasLevel(InfoLevel) {
		sl.outputFields(InfoLevel, msg, keyVals)
	}
}
func (sl *stdLogger) Warnw(msg string, keyVals ...interface{}) {
	if sl.HasLevel(WarnLevel) {
		sl.outputFields(WarnLevel, msg, keyVals)
	}
}
func (sl *stdLogger) Errorw(msg string, keyVals ...interface{}) {
	if sl.HasLevel(ErrorLevel) {
		sl.outputFields(ErrorLevel, msg, keyVals)
	}
}
func (sl *stdLogger) Fatalw(msg string, keyVals ...interface{}) {
	if sl.HasLevel(FatalLevel) {
		sl.outputFields(FatalLevel, msg, keyVals)
	}
	os.Exit(1)
}
func (sl *stdLogger) Panicw(msg string, keyVals ...interface{}) {
	if sl.HasLevel(PanicLevel) {
		sl.outputFields(PanicLevel, msg, keyVals)
	}
	panic(SimpleFormatter(msg, keyVals, "="))
}
