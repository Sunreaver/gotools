package logger

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type instance struct {
	logger *zap.SugaredLogger
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
	directory string
	level     zapcore.LevelEnabler
	// Logger zap.Logger实例
	Logger *zap.SugaredLogger
)

// InitLogger 初始化
// path 输出路径
// debugLevel 是否输出debug信息
// location 日志文件名所属时区/默认东八区
func InitLogger(path string, debugLevel bool, location *time.Location) {
	directory = path
	if debugLevel {
		level = zap.DebugLevel
	} else {
		level = zap.InfoLevel
	}
	//Fix time offset for Local
	// lt := time.FixedZone("Asia/Shanghai", 8*60*60)
	if location == nil {
		location = time.FixedZone("Asia/Shanghai", 8*60*60)
	}
	time.Local = location

	Logger = GetLogger(time.Now().Format("2006-01-02"))

	go func() {
		lastFile := time.Now().Format("2006-01-02")
		for {
			time.Sleep(time.Minute)
			if lastFile != time.Now().Format("2006-01-02") {
				lastFile = time.Now().Format("2006-01-02")
				Logger = GetLogger(lastFile)
			}
		}
	}()
}

// BeiJingTimeFormatter encodes the entry time as an RFC3339-formatted string under
// the provided key.
// func BeiJingTimeFormatter(key string) zap.TimeFormatter {
// 	return zap.TimeFormatter(func(t time.Time) zap.Field {
// 		return zap.String(key, t.Local().Format(time.RFC3339))
// 	})
// }

func localTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.Local()
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *loggerMap) Get(name string) *zap.SugaredLogger {
	i, ok := l.instances[name]
	if !ok {
		l.lock.Lock()
		i, ok = l.instances[name]
		if !ok {
			writer := &lumberjack.Logger{
				Filename: path.Join(directory, name+".log"),
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
				level,
			))
			i = instance{
				logger: logger.Sugar(),
				writer: writer,
			}
			l.instances[name] = i
		}
		defer l.lock.Unlock()
	}
	return i.logger
}

// RotateLog to causes Logger to close the existing log file
// and immediately create a new one.
func RotateLog() {
	loggers.lock.RLock()
	for _, i := range loggers.instances {
		i.writer.Rotate()
	}
	loggers.lock.RUnlock()
}

func exists(path string) error {
	stat, err := os.Stat(path)
	if err == nil {
		return errors.Wrap(err, "directory")
	}
	if os.IsNotExist(err) {
		return errors.New("path is not exists: " + path)
	}
	if !stat.IsDir() {
		return errors.New("path is not directory: " + path)
	}
	return err
}

// GetLogger to get zap.SugaredLogger
func GetLogger(name string) *zap.SugaredLogger {
	return loggers.Get(name)
}
