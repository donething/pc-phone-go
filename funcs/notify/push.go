package notify

import (
	"github.com/donething/utils-go/dowx"
	. "pc-phone-go/conf"
	"pc-phone-go/funcs/logger"
)

var (
	// WXQiYe 微信推送
	WXQiYe *dowx.QiYe
)

// WXPushCard 推送微信卡片消息
func WXPushCard(title string, description string, url string, btnText string) {
	if !initPush() {
		return
	}

	err := WXQiYe.PushCard(Conf.Comm.WXPush.Agentid, title, description, "", url, btnText)
	if err != nil {
		logger.Error.Printf("微信推送消息出错：%s\n", err)
		return
	}
}

// WXPushText 推送微信文本消息
func WXPushText(content string) {
	if !initPush() {
		return
	}
	err := WXQiYe.PushText(Conf.Comm.WXPush.Agentid, content, "")
	if err != nil {
		logger.Error.Printf("微信推送消息出错：%s\n", err)
		return
	}
}

// 初始化
func initPush() bool {
	// 初始化微信推送
	if WXQiYe == nil {
		if Conf.Comm.WXPush.Appid == "" || Conf.Comm.WXPush.Secret == "" {
			logger.Warn.Printf("微信推送的 Token 为空，无法推送消息\n")
			return false
		}
		WXQiYe = dowx.NewQiYe(Conf.Comm.WXPush.Appid, Conf.Comm.WXPush.Secret)
	}

	return true
}
