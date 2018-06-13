package logger

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// Config logger config
type Config struct {
	Path     string
	Loglevel LevelString
	// MaxSize 单文件最大存储，单位MB
	MaxSize int
}

// InitLoggerWithConfig 使用config初始化logger
func InitLoggerWithConfig(cfg Config, location *time.Location) error {
	if cfg.Path == "" {
		cfg.Path = "./"
	}
	if e := exists(cfg.Path); e != nil {
		return e
	} else if cfg.MaxSize <= 0 {
		return errors.New("MaxSize must be large than zero")
	}
	config = cfg

	// Fix time offset for Local
	// lt := time.FixedZone("Asia/Shanghai", 8*60*60)
	if location != nil {
		time.Local = location
	}

	lastFile := time.Now().Format(loggerByDayFormat)
	LoggerByDay = GetSugarLogger(lastFile)
	go func() {
		for {
			now := time.Now()
			if lastFile != now.Format(loggerByDayFormat) {
				go func(name string) {
					if e := loggers.Close(name); e != nil {
						fmt.Println("writer.Close error", e.Error(), "File", name)
					}
				}(lastFile)

				lastFile = now.Format(loggerByDayFormat)
				LoggerByDay = GetSugarLogger(lastFile)
			}
			time.Sleep(ToEarlyMorningTimeDuration(now))
		}
	}()

	return nil
}

// InitLoggerWithLevel 使用String格式的level初始化logger
// path 输出路径, 默认当前路径
// logLevel 日志级别: debug,info,warn
// location 日志文件名所属时区
func InitLoggerWithLevel(path string, logLevel LevelString, location *time.Location) error {
	return InitLogger(path, logLevel.toLevel(), location)
}

// InitLogger 初始化
// path 输出路径, 默认当前路径
// logLevel 日志级别
// location 日志文件名所属时区
func InitLogger(path string, logLevel Level, location *time.Location) error {
	return InitLoggerWithConfig(Config{
		Path:     path,
		Loglevel: logLevel.toLevelString(),
		MaxSize:  1024,
	}, location)
}

// GetLogger to get logger
func GetLogger(name string) *zap.Logger {
	return loggers.Get(name)
}

// GetSugarLogger to get SugaredLogger
func GetSugarLogger(name string) *zap.SugaredLogger {
	return GetLogger(name).Sugar()
}

// FlushAndCloseLogger flush and close logger
func FlushAndCloseLogger(name string) error {
	return loggers.Close(name)
}

func exists(path string) error {
	stat, err := os.Stat(path)
	if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return errors.New("path is not exists: " + path)
	} else if stat != nil && !stat.IsDir() {
		return errors.New("path is not directory: " + path)
	} else if stat == nil {
		return errors.New("not directory: " + path)
	}
	return err
}
