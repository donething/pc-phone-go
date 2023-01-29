package conf

import (
	"encoding/json"
	"github.com/donething/utils-go/dofile"
	"os"
	"path"
	"pc-phone-go/funcs/logger"
)

type Config struct {
	Comm struct {
		// 使用代理，为空表示不使用代理
		Proxy string `json:"proxy"`

		// 微信推送
		WXPush struct {
			Appid   string `json:"appid"`   // 组织 ID
			Secret  string `json:"secret"`  // 秘钥
			Agentid int    `json:"agentid"` // 应用（频道） ID
		}
	} `json:"comm"`

	// 青龙面板
	QLPanel struct {
		// 默认 5700
		Port         int    `json:"port"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"ql_panel"`

	Pics struct {
		// 工作池的容量，默认 10
		WorkerCount int `json:"worker_count"`
		// 对文件数据的处理，可从常量中选择 Handler***
		Handler string `json:"handler"`
		// 当 Handler 的值为 HandlerToLocal 时，保存文件到的本地目录
		LocalRoot string `json:"local_root"`
		// 推送
		// 一刻相册
		Yike struct {
			// 在一刻页面的网络调试工具中选择一个 fetch/XHR 请求，点击 Payload 标签可看到
			Bdstoken string `json:"bdstoken"`
			Cookie   string `json:"cookie"`
		} `json:"toyike"`
		// Telegram 推送消息
		TG struct {
			PicSaveToken  string `json:"pic_save_token"`
			PicSaveChatID string `json:"pic_save_chat_id"`
		} `json:"tg"`
	} `json:"pics"`

	// Javlib
	Javlib struct {
		// 视频目录的数组
		FanDirs []string `json:"fan_dirs"`
		// 字幕文件目录
		SubDir string `json:"sub_dir"`
	} `json:"javlib"`
}

const (
	// Name 配置文件的名字
	Name = "pc-phone-go.json"
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

		bs, err = json.MarshalIndent(Conf, "", "  ")
		fatal(err)
		_, err = dofile.Write(bs, confPath+".bak", os.O_CREATE|os.O_TRUNC, 0644)
		fatal(err)
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
