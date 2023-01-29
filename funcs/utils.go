package funcs

import (
	"os"
	"path/filepath"
	"pc-phone-go/funcs/logger"
)

// FileDir 文件保存的目录，默认为用户下载目录，无法获取则为当前运行目录下的"Downloads"文件夹
func FileDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "Downloads")
}

func CheckErr(err error) {
	if err != nil {
		logger.Error.Panicln(err)
	}
}
