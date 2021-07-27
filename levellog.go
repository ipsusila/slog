package slog

// LevelLogger base
type LevelLoggerBase struct {
	level Level
}

// NewLevelLoggerBase return new instance of level logger base
func NewLevelLoggerBase(lv Level) *LevelLoggerBase {
	bl := &LevelLoggerBase{}
	bl.SetLevel(lv)
	return bl
}

// HasLevel return current logger level
func (b *LevelLoggerBase) HasLevel(lv Level) bool {
	return b.level.Has(lv)
}

// SetLevel set logger level using flags
func (b *LevelLoggerBase) SetLevel(lv Level) {
	lvlIdx := 0
	for i, l := range lvAll {
		if lv.Has(l) {
			lvlIdx = i
		}
	}
	for i := 0; i < lvlIdx; i++ {
		lv.Set(lvAll[i])
	}
	b.level = 0
	b.level.Set(lv)
}
