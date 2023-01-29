package mysse

import (
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"strings"
	"time"
)

const (
	// 消息通道的缓存容量
	eventsCount = 100
	exitTimeout = "Exit: timeout"
)

// Send 使用 SSE 向客户端传递消息
//
// 消息的格式：{"success": bool, "msg": string, "data": string}
func Send(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 提取 hash
	hash := c.Query("hash")
	if strings.TrimSpace(hash) == "" {
		logger.Warn.Printf("SSE 连接缺少参数 hash，无法传输消息")
		c.Status(204)
		return
	}

	// 保存该请求的消息通道
	eventCh := make(chan sse.Event, eventsCount)
	MuChs.Lock()
	EventsChs[hash] = &EventsChsData{EventCh: eventCh, Tick: time.Now()}
	MuChs.Unlock()

	// 定时检测客户端是否关闭了连接，若是则删除其消息通道
	// 客户端每 1 分钟发送一次心跳；服务端每 3 分钟检测心跳
	// 如果 2 分钟内都没有收到心跳，表示客户端已关闭连接，服务端也退出该请求函数、移除该消息通道
	go func() {
		for range time.Tick(3 * time.Minute) {
			MuChs.Lock()
			last := EventsChs[hash].Tick
			MuChs.Unlock()
			if last.Add(2 * time.Minute).Before(time.Now()) {
				MuChs.Lock()
				delete(EventsChs, hash)
				MuChs.Unlock()
				logger.Warn.Printf("删除了已超时的消息通道，hash:'%s'\n", hash)
				// 发送结束该请求的消息
				SendToEventCh(eventCh, NewMsg("message", true, exitTimeout, nil))
				return
			}
		}
	}()

	// 发送 Hello 回应
	logger.Info.Printf("有新的 SSE 连接，hash:'%s'\n", hash)
	c.Render(-1, sse.Event{Data: "已连接到服务端"})
	c.Writer.Flush()

	// 从消息通道获取消息，发送到客户端
	for event := range eventCh {
		// logger.Info.Printf("收到需发送的消息：%#v\n", event.Data)
		if event.Data.(Data).Msg == exitTimeout {
			logger.Info.Printf("已收到超时退出请求的消息，hash:'%s'\n", hash)
			break
		}
		c.Render(-1, event)
		c.Writer.Flush()
	}
	logger.Info.Printf("已结束 SSE 请求，hash:'%s'\n", hash)
}

// Tick 接收客户端的心跳
//
// 返回值 {"code": int, "msg": string, "data": interface{}}
func Tick(c *gin.Context) {
	// 提取视频的所有分段链接
	var params Params
	err := c.Bind(&params)
	if err != nil {
		text := fmt.Sprintf("绑定心跳的参数时出错：%s", err)
		logger.Error.Println(text)
		c.JSON(http.StatusOK, entity.Rest{Code: 4000, Msg: text, Data: err.Error()})
		return
	}
	MuChs.Lock()
	EventsChs[params.Hash].Tick = time.Now()
	MuChs.Unlock()
	logger.Info.Printf("服务端已接收到心跳，hash:'%s'\n", params.Hash)
	c.JSON(http.StatusOK, entity.Rest{Code: 0,
		Msg: fmt.Sprintf("服务端已接收到心跳，hash:'%s'", params.Hash), Data: nil})
}

// GetAll 获取所有的消息通道
//
// 响应 格式：{"code": int, "msg": string, "data": string}
func GetAll(c *gin.Context) {
	MuChs.Lock()
	var keys = make([]string, len(EventsChs))
	i := 0
	for key := range EventsChs {
		keys[i] = key
		i++
	}
	MuChs.Unlock()
	c.JSON(200, entity.Rest{Code: 0, Msg: "所有的消息通道", Data: keys})
}
