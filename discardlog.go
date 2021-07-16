package slog

import (
	"fmt"
	"io"
	"os"
)

// Name of discard logger
const DiscardLoggerName = "discard"

type discardConstructor struct{}

// logger without output, except for panic
type discardLogger struct {
	LevelLoggerBase
}

func init() {
	Register(DiscardLoggerName, &discardConstructor{})
}

// NewDiscardLogger creates discard logger
func NewDiscardLogger(l Level) Logger {
	return &discardLogger{LevelLoggerBase: LevelLoggerBase{level: l}}
}

func (c *discardConstructor) New(w io.Writer, l Level) (Logger, error) {
	return &discardLogger{LevelLoggerBase: LevelLoggerBase{level: l}}, nil
}

func (c *discardConstructor) NewWithOptions(w io.Writer, l Level, op Options) (Logger, error) {
	return &discardLogger{LevelLoggerBase: LevelLoggerBase{level: l}}, nil
}

func (d *discardLogger) Trace(args ...interface{}) {
}
func (d *discardLogger) Debug(args ...interface{}) {
}
func (d *discardLogger) Print(args ...interface{}) {
}
func (d *discardLogger) Info(args ...interface{}) {
}
func (d *discardLogger) Warn(args ...interface{}) {
}
func (d *discardLogger) Error(args ...interface{}) {
}
func (d *discardLogger) Fatal(args ...interface{}) {
	os.Exit(1)
}
func (d *discardLogger) Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	panic(s)
}

func (d *discardLogger) Traceln(args ...interface{}) {
}
func (d *discardLogger) Debugln(args ...interface{}) {
}
func (d *discardLogger) Println(args ...interface{}) {
}
func (d *discardLogger) Infoln(args ...interface{}) {
}
func (d *discardLogger) Warnln(args ...interface{}) {
}
func (d *discardLogger) Errorln(args ...interface{}) {
}
func (d *discardLogger) Fatalln(args ...interface{}) {
	os.Exit(1)
}
func (d *discardLogger) Panicln(args ...interface{}) {
	s := fmt.Sprintln(args...)
	panic(s)
}

func (d *discardLogger) Tracef(format string, args ...interface{}) {
}
func (d *discardLogger) Debugf(format string, args ...interface{}) {
}
func (d *discardLogger) Printf(format string, args ...interface{}) {
}
func (d *discardLogger) Infof(format string, args ...interface{}) {
}
func (d *discardLogger) Warnf(format string, args ...interface{}) {
}
func (d *discardLogger) Errorf(format string, args ...interface{}) {
}
func (d *discardLogger) Fatalf(format string, args ...interface{}) {
	os.Exit(1)
}
func (d *discardLogger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	panic(s)
}

func (d *discardLogger) Tracew(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Debugw(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Printw(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Infow(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Warnw(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Errorw(msg string, keyVals ...interface{}) {
}
func (d *discardLogger) Fatalw(msg string, keyVals ...interface{}) {
	os.Exit(1)
}
func (d *discardLogger) Panicw(msg string, keyVals ...interface{}) {
	panic(SimpleFormatter(msg, keyVals, "="))
}
