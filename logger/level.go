package logger

import "go.uber.org/zap/zapcore"

// Level logger level
type Level int8

const (
	// DebugLevel 输出debug、info、warn、error级别.
	// 开发中用.
	DebugLevel Level = iota - 1
	// InfoLevel 输出info、warn、error级别.
	InfoLevel
	// WarnLevel 输出warn、error级别.
	WarnLevel
)

func (l Level) toZapcoreLevel() zapcore.LevelEnabler {
	if l < DebugLevel || l > WarnLevel {
		return zapcore.DebugLevel
	}
	return zapcore.Level(l)
}
