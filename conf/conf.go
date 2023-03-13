package conf

import (
	"github.com/donething/utils-go/doconf"
	"github.com/donething/utils-go/dolog"
	"os"
	"pc-phone-go/funcs/logger"
)

// Name 配置文件的名字
const confPath = "./pc-phone-go.json"

// Conf 配置的实例
var Conf Config

func init() {
	exist, err := doconf.Init(confPath, &Conf)
	dolog.CkPanic(err)

	if !exist {
		logger.Warn.Printf("请填写配置文件后，重新运行\n")
		os.Exit(0)
	}
}
