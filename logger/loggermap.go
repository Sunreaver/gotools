package logger

import (
	"path"
	"sync"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type instance struct {
	logger *zap.Logger
	writer *lumberjack.Logger
}

type loggerMap struct {
	lock      *sync.RWMutex
	instances map[string]instance
}

var (
	loggers = loggerMap{
		new(sync.RWMutex),
		make(map[string]instance),
	}
	config Config

	// LoggerByDay 按照天来划分的logger
	LoggerByDay *zap.SugaredLogger
)

const (
	loggerByDayFormat = "2006-01-02.log"
)

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.Local()
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *loggerMap) Close(name string) error {
	l.lock.RLock()
	_, ok := l.instances[name]
	l.lock.RUnlock()

	if ok {
		l.lock.Lock()
		defer l.lock.Unlock()
		i, ok := l.instances[name]
		if ok {
			if e := i.logger.Sync(); e != nil {
				return e
			}
			if e := i.writer.Close(); e != nil {
				return e
			}
			delete(l.instances, name)
		}
	}

	return nil
}

func (l *loggerMap) Get(name string) *zap.Logger {
	l.lock.RLock()
	i, ok := l.instances[name]
	l.lock.RUnlock()

	if !ok {
		writer := &lumberjack.Logger{
			Filename: path.Join(config.Path, name),
			MaxSize:  int(config.MaxSize),
		}
		ws := zapcore.AddSync(writer)
		cfg := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     localTimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
		}
		logger := zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			ws,
			config.Loglevel.toZapcoreLevel(),
		))
		i = instance{
			logger: logger,
			writer: writer,
		}

		l.lock.Lock()
		if tmp, ok := l.instances[name]; !ok {
			l.instances[name] = i
		} else {
			i = tmp
		}
		l.lock.Unlock()
	}
	return i.logger
}

// ToEarlyMorningTimeDuration will 计算当前到第二日凌晨的时间
func ToEarlyMorningTimeDuration(now time.Time) time.Duration {
	hour := 24 - now.Hour() - 1
	minute := 60 - now.Minute() - 1
	second := 60 - now.Second()

	return time.Duration(hour)*time.Hour +
		time.Duration(minute)*time.Minute +
		time.Duration(second)*time.Second
}
