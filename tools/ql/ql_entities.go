package ql

import "time"

// BasicResp 基础响应内容
type BasicResp struct {
	Code int `json:"code"` // 200 表示正确
}

// TokenResp 获取 Token 时的响应内容
type TokenResp struct {
	*BasicResp
	Data struct {
		Token      string `json:"token"`      // 获取的 Token
		TokenType  string `json:"token_type"` // Token 类型，此时为"Bearer"
		Expiration int    `json:"expiration"` // 过期时间，30天
	} `json:"data"`
}

// Env 环境变量的类型
type Env struct {
	ID        int         `json:"id"`             // 更新环境变量时，需要指定
	Value     string      `json:"value"`          // 环境变量值
	Name      string      `json:"name,omitempty"` // 环境变量名
	Remarks   interface{} `json:"remarks"`        // 环境变量备注
	Status    int         `json:"status"`         // 0 表示启用，1 表示禁用
	Timestamp string      `json:"timestamp"`
	Position  float64     `json:"position"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// GetEnvsResp 获取环境变量的响应内容
type GetEnvsResp struct {
	*BasicResp
	Data []Env `json:"data"`
}

// SetEnvReq 设置环境变量的请求内容
type SetEnvReq struct {
	Value   string `json:"value"`
	Name    string `json:"name"`
	Remarks string `json:"remarks"`
}

// UpEnvReq 更新环境变量的请求内容
type UpEnvReq struct {
	// 更新环境变量需要指定 ID，否则会创建同名环境变量
	ID int `json:"id"`
	*SetEnvReq
}

// Cron 定时任务
// 可排除执行置顶或禁用的任务
type Cron struct {
	ID                int         `json:"id"`       // 任务 ID，可用于发送执行任务的请求时作为参数
	Name              string      `json:"name"`     // 任务名
	Command           string      `json:"command"`  // 执行命令，如"task helloworld.py"，可忽略
	Schedule          string      `json:"schedule"` // 定时执行的时间，如"0 0 * * *"
	Timestamp         string      `json:"timestamp"`
	Saved             bool        `json:"saved"`
	Status            int         `json:"status"`
	IsSystem          int         `json:"isSystem"`
	Pid               interface{} `json:"pid"`
	IsDisabled        int         `json:"isDisabled"` // 任务是否被禁用，0表示否；1表示是，可用于过滤
	IsPinned          int         `json:"isPinned"`   // 任务是否被置顶，0表示否；1表示是，可用于过滤
	LogPath           string      `json:"log_path"`
	Labels            []string    `json:"labels"`
	LastRunningTime   int         `json:"last_running_time"`
	LastExecutionTime int         `json:"last_execution_time"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

// GetCronsResp 获取所有定时任务的响应
type GetCronsResp struct {
	*BasicResp
	Data struct {
		Data  []Cron `json:"data"`
		Total int    `json:"total"`
	} `json:"data"`
}
