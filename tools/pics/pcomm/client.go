package pcomm

import (
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/funcs/notify"
	"time"
)

var (
	// Client 执行 HTTP 请求的客户端
	Client = dohttp.New(30*time.Second, false, false)
)

func init() {
	// 如果配置中指定了代理，需要设置
	proxy := Conf.Comm.Proxy
	if proxy != "" {
		err := Client.SetProxy(proxy)
		if err != nil {
			logger.Error.Println("下载发送图片出错：设置 HTTP 代理时出错：", err)
			notify.WXPushCard("发送图片出错",
				fmt.Sprintf("设置 HTTP 代理时出错：%s", err.Error()), "", "")
			return
		}
	}
}
