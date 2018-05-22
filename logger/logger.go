package logger

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// InitLogger 初始化
// path 输出路径, 默认当前路径
// debugLevel 是否输出debug信息
// location 日志文件名所属时区
func InitLogger(path string, logLevel Level, location *time.Location) error {
	if path == "" {
		path = "./"
	}
	if e := exists(path); e != nil {
		return e
	}

	directory = path
	level = logLevel.toZapcoreLevel()

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
