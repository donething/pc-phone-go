// Package pcomm 公共数据，不引用本工程下自己编写的库
package pcomm

const (
	// TagWeibo 微博
	TagWeibo = "weibo"
)

const (
	// HandlerToLocal HandlerToTG HandlerToYike 在 pworker 中对数据的处理方法
	HandlerToLocal = "ToLocal" // 保存到本地
	HandlerToTG    = "ToTG"    // 发送到 Telegram
	HandlerToYike  = "ToYike"  // 发送到一刻相册
)

const (
	// DBFail 数据库中图片下载相关数据的键前缀，下载失败的数据
	DBFail = "picsfail_"
	// DBSkip 数据库中图片下载相关数据的键前缀，下载跳过的数据
	DBSkip = "picsskip_"
)
