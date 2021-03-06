# slog [![Go Reference](https://pkg.go.dev/badge/github.com/ipsusila/slog)](https://pkg.go.dev/github.com/ipsusila/slog)

Golang **s**imple **log**ger wrapper. It supports leveled logger with the following levels:

1. **PANIC**, log message and call `panic`
2. **FATAL**, log message and exit by calling `os.Exit(1)`
3. **ERROR**, log message in `error` level
4. **WARNING**, log message in `warning` level
5. **INFO**, log message in `info` level
6. **DEBUG**, log message in `debug` level
7. **TRACE**, log message in `trace` level

Logger implements the following interface:

```go
type Logger interface {
    HasLevel(lv Level) bool
    SetLevel(lv Level)

    // Print like methods
    Trace(args ...interface{})
    Debug(args ...interface{})
    Print(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
    Panic(args ...interface{})

    // Println like methods
    Traceln(args ...interface{})
    Debugln(args ...interface{})
    Println(args ...interface{})
    Infoln(args ...interface{})
    Warnln(args ...interface{})
    Errorln(args ...interface{})
    Fatalln(args ...interface{})
    Panicln(args ...interface{})

    // Printf like methods
    Tracef(format string, args ...interface{})
    Debugf(format string, args ...interface{})
    Printf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
    Panicf(format string, args ...interface{})

    // Log with fields (key=value)
    Tracew(msg string, keyVals ...interface{})
    Debugw(msg string, keyVals ...interface{})
    Printw(msg string, keyVals ...interface{})
    Infow(msg string, keyVals ...interface{})
    Warnw(msg string, keyVals ...interface{})
    Errorw(msg string, keyVals ...interface{})
    Fatalw(msg string, keyVals ...interface{})
    Panicw(msg string, keyVals ...interface{})
}
```

`Logger` can be initialized with the following constructor

```go
New(name string, w io.Writer, l Level) (Logger, error)
NewWithOptions(w io.Writer, level Level, op Options) (Logger, error)
```

in which

- `name` specify logger name. Currently available values `discard`, `stdlog` and `logrus`
- `w` logger output
- `l` logger level
- `op` loger options

## Options

1. `discard`, discard log ouput except `panic`. Options not supported.
2. `stdlog`, standar logger options:

    - `timestampFormat`: timestamp layout format, see [`time.Time` format](https://pkg.go.dev/time#pkg-constants)
    - `disableColor`: to disable color in log

3. `logrus`, support options for [`logrus.TextFormatter` formatter](https://pkg.go.dev/github.com/sirupsen/logrus#TextFormatter) and [`logrus.JSONFormatter` formatter](https://pkg.go.dev/github.com/sirupsen/logrus#JSONFormatter).

    - `formatter`: logrus formatter, either `text` or `json`. Default format is `logrus.TextFormatter`
    - `timestampFormat`: timestamp layout format, see [`time.Time` format](https://pkg.go.dev/time#pkg-constants)
    - `reportCaller`: if set to `true`, the calling method will be added as a field
    - `fullTimestamp`: logging the full timestamp instead of elapsed time since application started, default to `true`
    - `disableTimestamp`: disable timestamp in log
    - `fieldMap`: customize default key names
	- `dataKey`: data key for `json` formatter
	- `prettyPrint`: pretty print `json` output
	- `disableHTMLEscape`: disable HTML escape in `json` formatter
	- `forceColors`: for color in `text` formatter
	- `disableColors`: disable colors in `text` formatter
	- `forceQuote`: force quote in `text` formatter
	- `disableQuote`: disable quote in `text` formatter
	- `environmentOverrideColors`: override color based on environment config in `text` formatter
	- `disableSorting`: disable key sorting in `text` formatter
	- `disableLevelTruncation`: disable log level string truncation in `text` formatter
	- `padLevelText`: add padding in level string in `text` formatter
	- `quoteEmptyFields`: add quote for empty log entry in `text` formatter

## Credits

- Color support via [https://github.com/fatih/color](https://github.com/fatih/color)
- Logrus logger using [https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)

## License

The MIT License (MIT), see [LICENSE](LICENSE).
