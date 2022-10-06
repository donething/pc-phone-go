package conf

import (
	"encoding/json"
	"github.com/donething/utils-go/dofile"
	"os"
	"path"
	"pc-phone-conn-go/logger"
)

type Config struct {
	// 使用代理，为空表示不使用代理
	Proxy string `json:"proxy"`

	// 青龙面板
	QLPanel struct {
		// 默认 5700
		Port         string `json:"port"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"ql_panel"`
}

const (
	// Name 配置文件的名字
	Name = "pc-phone-conn-go.json"
)

var (
	// Conf 配置的实例
	Conf Config
	// 配置文件所在的路径
	confPath string
)

func init() {
	confPath = path.Join(Name)
	exist, err := dofile.Exists(confPath)
	fatal(err)
	if exist {
		logger.Info.Printf("读取配置文件：'%s'\n", confPath)
		bs, err := dofile.Read(confPath)
		fatal(err)
		err = json.Unmarshal(bs, &Conf)
		fatal(err)

		// 指定默认配置
		if Conf.QLPanel.Port == "" {
			Conf.QLPanel.Port = "5700"
		}
	}

	bs, err := json.MarshalIndent(Conf, "", "  ")
	fatal(err)
	_, err = dofile.Write(bs, confPath, os.O_CREATE|os.O_TRUNC, 0644)
	fatal(err)
	logger.Info.Printf("已重写配置文件：'%s'\n", confPath)

	if !exist {
		logger.Info.Printf("请填写配置文件后，重新运行\n")
		os.Exit(0)
	}
}

// fatal 出错时，强制关闭程序
func fatal(err error) {
	if err != nil {
		logger.Error.Fatal(err)
	}
}