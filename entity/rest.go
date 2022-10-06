package entity

// Rest REST 响应体
type Rest struct {
	Errcode int         `json:"errcode"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
