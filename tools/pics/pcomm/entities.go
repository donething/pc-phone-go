package pcomm

// Album 向 pworker 发送的图集对象
type Album struct {
	// 基础，必需
	Plat    string   `json:"plat"`             // 所在的平台，如“微博”
	Caption string   `json:"caption"`          // 标题
	Created int64    `json:"created"`          // 创建时间
	ID      string   `json:"id"`               // 图集的 ID
	UID     string   `json:"uid"`              // 图集所属用户的 ID
	URLs    []string `json:"urls"`             // 发送的数据，如 URL 的数组（若为图片则为最大分辨率）
	URLsM   []string `json:"urls_m,omitempty"` // 若为图片，则为中等分辨率

	// 后续设置，可空
	// 表示为任务来源，以便在此处成功后删除记录。其值为失败、跳过的数据前缀
	// 空""为初次下载，pcomm.DBFail 为失败，pcomm.DBSkip 为跳过
	// 类型不能直接用 []byte，在持久化保存到数据库中时会按某种方式转为字符串，和 string() 产生的方式不一样
	RetryFrom string
	Header    map[string]string `json:"header,omitempty"` // 下载文件的请求头，可空
}

// PStatus 任务进度
type PStatus struct {
	Total int `json:"total"`
	Done  int `json:"done"`
	Skip  int `json:"skip"`
	Fail  int `json:"fail"`
}

// PCount 需要重试下载的图集数
type PCount struct {
	Fail int `json:"fail"`
	Skip int `json:"skip"`
}
