package conf

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

	// Javlib
	Javlib struct {
		// 视频目录的数组
		FanDirs []string `json:"fan_dirs"`
		// 字幕文件目录
		SubDir string `json:"sub_dir"`
	} `json:"javlib"`
}
