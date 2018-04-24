package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// InitLogger 初始化
// path 输出路径
// debugLevel 是否输出debug信息
// location 日志文件名所属时区
func InitLogger(path string, l Level, location *time.Location) error {
	if e := exists(path); e != nil {
		return e
	}

	directory = path
	level = l.toZapcoreLevel()

	// Fix time offset for Local
	// lt := time.FixedZone("Asia/Shanghai", 8*60*60)
	if location != nil {
		time.Local = location
	}

	lastFile := time.Now().Format(loggerByDayFormat)
	LoggerByDay = GetLogger(lastFile)
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
				LoggerByDay = GetLogger(lastFile)
			}
			time.Sleep(ToEarlyMorningTimeDuration(now))
		}
	}()

	return nil
}

// GetLogger to get logger
func GetLogger(name string) *zap.SugaredLogger {
	return loggers.Get(name)
}

// FlushAndCloseLogger flush and close logger
func FlushAndCloseLogger(name string) error {
	return loggers.Close(name)
}
