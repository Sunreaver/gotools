package goTools

import (
	"os"
	"os/exec"
	"path/filepath"
)

//获取当前运行目录
func CurPath() (path string) {
	file, _ := exec.LookPath(os.Args[0])
	pt, _ := filepath.Abs(file)

	return filepath.Dir(pt)
}

//获取系统分隔符
func SystemSep() (path string) {

	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	return path
}

//判断目录是否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

//判断文件是否存在
func IsFileExists(filePath string) bool {
	fi, err := os.Open(filePath)
	if err != nil {
		return false
	}
	fi.Close()
	return true
}
