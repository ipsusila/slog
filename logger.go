package slog

import (
	"errors"
	"io"
	"sync"
)

// Register Logger constructor
var (
	constructorsMu sync.RWMutex
	constructors   = make(map[string]Constructor)
)

// Package logger instance
var (
	pkgLoggerMu sync.Mutex
	discardL    *discardLogger = &discardLogger{LevelLoggerBase: LevelLoggerBase{level: AllLevel}}
	pkgLogger   Logger         = discardL
)

// Unknown field name
var UnknownFieldName = "@logfield"

// Register logger constructor
func Register(name string, constructor Constructor) {
	constructorsMu.Lock()
	defer constructorsMu.Unlock()
	if constructor == nil {
		panic("logger: Register Constructor is nil")
	}
	if _, dup := constructors[name]; dup {
		panic("logger: Register called twice for Constructor " + name)
	}
	constructors[name] = constructor
}

// ConstructorFor return sql query for given name
func ConstructorFor(name string) (Constructor, bool) {
	constructorsMu.RLock()
	defer constructorsMu.RUnlock()
	for key, constructor := range constructors {
		if key == name {
			return constructor, true
		}
	}
	return nil, false
}

// Global discard logger
var Discard Logger = &discardLogger{LevelLoggerBase: LevelLoggerBase{level: AllLevel}}

// Constructor for reader creator
type Constructor interface {
	New(w io.Writer, level Level) (Logger, error)
	NewWithOptions(w io.Writer, level Level, op Options) (Logger, error)
}

// Loger interface
type Logger interface {
	HasLevel(lv Level) bool
	SetLevel(lv Level)

	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Println(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Tracew(msg string, keyVals ...interface{})
	Debugw(msg string, keyVals ...interface{})
	Printw(msg string, keyVals ...interface{})
	Infow(msg string, keyVals ...interface{})
	Warnw(msg string, keyVals ...interface{})
	Errorw(msg string, keyVals ...interface{})
	Fatalw(msg string, keyVals ...interface{})
	Panicw(msg string, keyVals ...interface{})
}

// New create new logger with given name
func New(name string, w io.Writer, l Level) (Logger, error) {
	// get constructor
	c, ok := ConstructorFor(name)
	if !ok {
		return nil, errors.New("unknown logger: " + name)
	}

	lgr, err := c.New(w, l)
	if err != nil {
		return nil, err
	}
	return lgr, nil
}

// NewWithOptions create new logger with given name, and options for spesific logger
func NewWithOptions(name string, w io.Writer, l Level, op Options) (Logger, error) {
	// get constructor
	c, ok := ConstructorFor(name)
	if !ok {
		return nil, errors.New("unknown logger: " + name)
	}

	lgr, err := c.NewWithOptions(w, l, op)
	if err != nil {
		return nil, err
	}
	return lgr, nil
}

// MustUse specific logger
func MustUse(name string, w io.Writer, l Level) {
	if err := Use(name, w, l); err != nil {
		panic(err)
	}
}

// Use specific logger
func Use(name string, w io.Writer, l Level) error {
	pkgLoggerMu.Lock()
	defer pkgLoggerMu.Unlock()

	lgr, err := New(name, w, l)
	if err != nil {
		return err
	}
	pkgLogger = lgr

	return nil
}

// may be unsafe?
func HasLevel(lv Level) bool {
	return pkgLogger.HasLevel(lv)
}
func SetLevel(lv Level) {
	pkgLogger.SetLevel(lv)
}

func Trace(args ...interface{}) {
	pkgLogger.Trace(args...)
}
func Debug(args ...interface{}) {
	pkgLogger.Debug(args...)
}
func Print(args ...interface{}) {
	pkgLogger.Print(args...)
}
func Info(args ...interface{}) {
	pkgLogger.Info(args...)
}
func Warn(args ...interface{}) {
	pkgLogger.Warn(args...)
}
func Error(args ...interface{}) {
	pkgLogger.Error(args...)
}
func Fatal(args ...interface{}) {
	pkgLogger.Fatal(args...)
}
func Panic(args ...interface{}) {
	pkgLogger.Panic(args...)
}

func Traceln(args ...interface{}) {
	pkgLogger.Traceln(args...)
}
func Debugln(args ...interface{}) {
	pkgLogger.Debugln(args...)
}
func Println(args ...interface{}) {
	pkgLogger.Println(args...)
}
func Infoln(args ...interface{}) {
	pkgLogger.Infoln(args...)
}
func Warnln(args ...interface{}) {
	pkgLogger.Warnln(args...)
}
func Errorln(args ...interface{}) {
	pkgLogger.Errorln(args...)
}
func Fatalln(args ...interface{}) {
	pkgLogger.Fatalln(args...)
}
func Panicln(args ...interface{}) {
	pkgLogger.Panicln(args...)
}

func Tracef(format string, args ...interface{}) {
	pkgLogger.Tracef(format, args...)
}
func Debugf(format string, args ...interface{}) {
	pkgLogger.Debugf(format, args...)
}
func Printf(format string, args ...interface{}) {
	pkgLogger.Printf(format, args...)
}
func Infof(format string, args ...interface{}) {
	pkgLogger.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	pkgLogger.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	pkgLogger.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	pkgLogger.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	pkgLogger.Panicf(format, args...)
}

func Tracew(msg string, keyVals ...interface{}) {
	pkgLogger.Tracew(msg, keyVals...)
}
func Debugw(msg string, keyVals ...interface{}) {
	pkgLogger.Debugw(msg, keyVals...)
}
func Printw(msg string, keyVals ...interface{}) {
	pkgLogger.Printw(msg, keyVals...)
}
func Infow(msg string, keyVals ...interface{}) {
	pkgLogger.Infow(msg, keyVals...)
}
func Warnw(msg string, keyVals ...interface{}) {
	pkgLogger.Warnw(msg, keyVals...)
}
func Errorw(msg string, keyVals ...interface{}) {
	pkgLogger.Errorw(msg, keyVals...)
}
func Fatalw(msg string, keyVals ...interface{}) {
	pkgLogger.Fatalw(msg, keyVals...)
}
func Panicw(msg string, keyVals ...interface{}) {
	pkgLogger.Panicw(msg, keyVals...)
}
