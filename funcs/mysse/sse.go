package mysse

import (
	"fmt"
	"github.com/gin-contrib/sse"
	"pc-phone-go/funcs/logger"
	"sync"
	"time"
)

// Params request 请求的基础参数
type Params struct {
	// 消息通道的 hash
	Hash string `json:"hash"`
}
type EventsChsData struct {
	EventCh chan sse.Event
	Tick    time.Time
}

var (
	// EventsChs 向客户端实时传递消息
	//
	// 每个请求各用一个通知 channel，其键为 main.go 中路由"/api/mysse?hash=xxx"的参数"xxx"
	//
	// 每个请求需要唯一，如用当前时间的毫秒值作为 hash 等
	EventsChs = make(map[string]*EventsChsData)
	MuChs     = sync.Mutex{}
)

// NewMsg 返回 SSE 消息实体
func NewMsg(event string, success bool, msg string, data interface{}) sse.Event {
	return sse.Event{
		Event: event,
		Data:  Data{Success: success, Msg: msg, Data: data},
		Id:    fmt.Sprintf("%s_%d", event, time.Now().UnixNano()),
		Retry: 5,
	}
}

// SendToEventCh 发送消息到 SSE
func SendToEventCh(ch chan sse.Event, data sse.Event) {
	// 就是要判断通道是否为 nil，以免阻塞协程
	if ch != nil {
		logger.Info.Printf("发送消息：%v\n", data.Data)
		ch <- data
	}
}
