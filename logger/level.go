package logger

import "go.uber.org/zap/zapcore"

// Level logger level
type Level int8

const (
	// DebugLevel 最低输出debug级别.
	// 开发中用.
	DebugLevel Level = iota - 1
	// InfoLevel 最低输出InfoLevel级别.
	InfoLevel
	// WarnLevel 最低输出WarnLevel级别.
	WarnLevel
)

func (l Level) toZapcoreLevel() zapcore.LevelEnabler {
	if l < DebugLevel || l > WarnLevel {
		return zapcore.DebugLevel
	}
	return zapcore.Level(l)
}
