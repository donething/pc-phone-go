package mysse

import "encoding/json"

// Data SSE 消息的 data 部分
type Data struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// 转换为 JSON 文本
func (s *Data) String() string {
	bs, _ := json.Marshal(*s)
	return string(bs)
}
