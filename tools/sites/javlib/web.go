package javlib

import (
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-go/conf"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/funcs/mysse"
)

var client = dohttp.New(false, false)

// ParamsRename request 重命名的参数
type ParamsRename struct {
	*mysse.Params
	Paths []string `json:"paths"`
}

const (
	// 重命名完一个文件后发送事件给前端更新界面
	eventRename = "rename_update"
)

func init() {
	// 是否设置网络代理
	if conf.Conf.Comm.Proxy != "" {
		err := client.SetProxy(conf.Conf.Comm.Proxy)
		if err != nil {
			logger.Error.Printf("设置代理'%s'出错：%s\n", conf.Conf.Comm.Proxy, err)
			return
		}
	}
}

// ExistFanhaoFile 查询番号列表对应的文件列表
//
// 参数：POST 传递 json 数据，番号列表：[string, string, string]
//
// 返回：data：番号的路径，为空""表示没有相应的文件：
//
// {"code": int, "msg": string, "data": {"fanhao1": string, "fanhao2": string}}
func ExistFanhaoFile(c *gin.Context) {
	var fanhaos []string
	err := c.Bind(&fanhaos)
	if err != nil {
		msg := fmt.Sprintf("绑定字幕的番号参数出错：%s", err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3000, Msg: msg, Data: false})
		return
	}
	paths, err := SeekFile(conf.Conf.Javlib.FanDirs, fanhaos)
	if err != nil {
		msg := fmt.Sprintf("查找番号的本地文件时出错：%s\n", err)
		logger.Error.Printf(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3010, Msg: msg, Data: false})
		return
	}
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "查找番号的结果", Data: paths})
}

// ExistSubtitle 查找番号对应的字幕
//
// 参数：POST 传递 json 数据，番号：{"fanhao": string}
//
// 返回 其中若有字幕，则 data 为字幕的路径：
//
// {"code": int, "msg": string, "data": {[string]:string}}
func ExistSubtitle(c *gin.Context) {
	var fanhao struct {
		Fanhao string `json:"fanhao"`
	}
	err := c.Bind(&fanhao)
	if err != nil {
		msg := fmt.Sprintf("绑定字幕的番号参数出错：%s", err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3020, Msg: msg, Data: false})
		return
	}

	if fanhao.Fanhao == "" {
		c.JSON(http.StatusOK, entity.Rest{Code: 3021, Msg: "关键字为空", Data: false})
		return
	}
	path, err := MatchSubtitle(fanhao.Fanhao)
	if err != nil {
		msg := fmt.Sprintf("查找番号(%s)本地的字幕文件时出错：%s", fanhao, err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3030, Msg: msg, Data: false})
		return
	}

	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "字幕的路径", Data: path})
}

// RenameDir 重命名文件夹内的文件
//
// 参数 POST 传递 json 数据，"path"为文件夹的路径列表：[string]
//
// 返回消息 当成功时 msg 为番号名， data 为新路径；失败时 msg 为错误说的描述， data 为错误的详情（可空）
//
// {"success": bool, "msg": string, "data": interface{}}
//
// 返回 通用信息，只表示操作完成，不指明成功还是失败
func RenameDir(c *gin.Context) {
	var params ParamsRename
	err := c.Bind(&params)
	if err != nil {
		msg := fmt.Sprintf("重命名提取路径参数时出错：%s", err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3050, Msg: msg, Data: nil})
		return
	}
	// 重命名并实时反馈消息
	// 获取消息通道
	mysse.MuChs.Lock()
	eventCh := mysse.EventsChs[params.Hash].EventCh
	mysse.MuChs.Unlock()
	if eventCh == nil {
		msg := fmt.Sprintf("重命名的 channel 为空，终止重命名")
		logger.Info.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 3060, Msg: msg, Data: nil})
		return
	}
	// 重命名
	Rename(params.Paths, eventCh)
	// 慎重关闭，会导致该 SSE 请求完成然后退出，客户端疯狂重连
	// 等待重命名操作完成，关闭消息通道
	// close(eventCh)
	msg := "已完成重命名操作"
	logger.Info.Println(msg)
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: msg, Data: nil})
}
