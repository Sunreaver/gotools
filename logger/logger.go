package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"github.com/sunreaver/gotools/system"
	"github.com/uber-go/zap"
)

// Logger is zap json logger wrapper.
var Logger zap.Logger

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("DEBUG", "0")
	viper.SetDefault("LOGFILE", "0")

	if runtime.GOOS == "linux" {
		viper.SetDefault("LOGPATH", "/mnt/log")
	} else {
		viper.SetDefault("LOGPATH", system.CurPath())
	}
	opts := getOptions()
	Logger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("time")), opts...)

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(86400-(time.Now().Unix()%86400)))
			opts := getOptions()
			Logger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("time")), opts...)
		}
	}()
}

func getFile() *os.File {
	path := viper.GetString("LOGPATH")
	fn := time.Now().Format("2006-01-02") + ".log"
	fp := filepath.Join(path, fn)
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func getOptions() []zap.Option {
	var opts []zap.Option
	if viper.GetString("DEBUG") != "0" {
		opts = append(opts, zap.DebugLevel)
	} else {
		opts = append(opts, zap.InfoLevel)
	}
	if viper.GetString("LOGFILE") != "0" {
		f := getFile()
		opts = append(opts, zap.Output(f))
	}
	return opts
}
