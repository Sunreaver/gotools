package log

import (
	L "log"
	"time"
)

// Println 输出带YYYY-mm-dd HH:MM:SS前缀的log信息
func Println(v ...interface{}) {
	L.Print(time.Now().Format("2006-01-02 15:04:05"))
	for _, item := range v {
		L.Print(item)
	}
	L.Print("\r\n")
}
