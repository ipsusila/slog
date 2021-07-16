package slog

import (
	"errors"
	"fmt"
	"strings"
)

// Level of the log
type Level uint32

const lvSep = "|"

// Logger level
const (
	PanicLevel Level = 1 << iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// AllLevel contains all log level
const AllLevel = PanicLevel | FatalLevel | ErrorLevel | WarnLevel | InfoLevel | DebugLevel | TraceLevel

var lvAll = []Level{PanicLevel, FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel, TraceLevel}

var lvStrMap = map[Level]string{
	PanicLevel: "panic",
	FatalLevel: "fatal",
	ErrorLevel: "error",
	WarnLevel:  "warn",
	InfoLevel:  "info",
	DebugLevel: "debug",
	TraceLevel: "trace",
	AllLevel:   "all",
}

// convert level to fixed string
var lvFixStrMap = map[Level]string{
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
	ErrorLevel: "ERROR",
	WarnLevel:  "WARNN",
	InfoLevel:  "INFOO",
	DebugLevel: "DEBUG",
	TraceLevel: "TRACE",
}

//

var strLvMap = map[string]Level{
	"panic": PanicLevel,
	"fatal": FatalLevel,
	"error": ErrorLevel,
	"warn":  WarnLevel,
	"info":  InfoLevel,
	"debug": DebugLevel,
	"trace": TraceLevel,
	"all":   AllLevel,
}

// Has specific level
func (l Level) Has(lv Level) bool {
	return l&lv != 0
}

// Set log level
func (l *Level) Set(lv Level) {
	*l |= lv
}

// Clear specific flag
func (l *Level) Clear(lv Level) {
	*l = *l &^ lv
}

// Toggle log level
func (l *Level) Toggle(lv Level) {
	*l = *l ^ lv
}

// Stringer interface
func (l Level) String() string {
	if b, err := l.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

// Levels return all level
func Levels() []Level {
	return lvAll
}

// LevelCount return number of valid log levels
func LevelsCount() int {
	return len(lvAll)
}

// LevelFixedString return fixed string for given level
func LevelFixedString(lv Level) string {
	if str, ok := lvFixStrMap[lv]; ok {
		return str
	}
	return "OTHER"
}

// ParseLevel string
func ParseLevel(level string) (Level, error) {
	lvStr := strings.ToLower(level)
	lv, ok := strLvMap[lvStr]
	if ok {
		return lv, nil
	}

	// parse combined flags
	lv = 0
	levels := strings.Split(lvStr, lvSep)
	for _, str := range levels {
		if v, ok := strLvMap[str]; ok {
			lv.Set(v)
		}
	}

	if lv == 0 {
		return 0, errors.New("unkown level: " + level)
	}

	return lv, nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (l *Level) UnmarshalText(text []byte) error {
	lv, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*l = lv

	return nil
}

// MarshalText return level as byte string
func (l Level) MarshalText() ([]byte, error) {
	sb := strings.Builder{}
	for i, lv := range lvAll {
		if l.Has(lv) {
			if i > 0 {
				sb.WriteString(lvSep)
			}
			sb.WriteString(lvStrMap[lv])
		}
	}
	if sb.Len() == 0 {
		return nil, fmt.Errorf("not a valid logger l %d", l)
	}

	// convert to byte
	lvStr := sb.String()
	return []byte(lvStr), nil
}
