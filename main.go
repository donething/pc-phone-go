package main

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"pc-phone-go/funcs"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/funcs/mysse"
	"pc-phone-go/icons"
	"pc-phone-go/tools/lives"
	"pc-phone-go/tools/ql"
	"pc-phone-go/tools/qx"
	"pc-phone-go/tools/sites/javlib"
	"runtime"
	"time"
)

const (
	// 服务端口
	port = 8800
)

func init() {
	// go func() {
	// 	// 在本应用运行后需等一段时间，等 Docker 启动目标容器后才执行脚本，用于电脑刚开机时
	// 	t := time.NewTimer(3 * time.Minute)
	// 	<-t.C
	//
	// 	_, err := ql.StartCommCronsCall()
	// 	if err != nil {
	// 		logger.Error.Printf("执行定时任务时出错：%s\n", err)
	// 		return
	// 	}
	// }()

	// 显示托盘
	go func() {
		runtime.LockOSThread()
		systray.Run(onReady, nil)
		runtime.UnlockOSThread()
	}()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 显示 PC 端地址的二维码
	router.GET("/", index)

	// 打开或显示本地文件
	router.POST("/api/openfile", funcs.OpenLocal)

	// 剪贴板
	router.POST("/api/clip", handerClip)

	// qx
	router.GET("/api/qx/parse_surge", qx.ParseSurge)

	// ql
	router.POST("/api/ql/set_env", ql.SetEnv)
	router.POST("/api/ql/start_comm_crons", ql.StartCommCrons)

	// lives
	router.GET("/api/lives/douyin/live", lives.GetDouyinRoom)

	// 使用 SSE 向客户端传递消息
	// 客户端请求需传递查询字符串参数 hash："?hash=xxx"，其值可看下面的说明
	// 尽量每个请求使用不同的消息通道，以避免当打开了多个连接时，消息传到其它请求上
	// 同时为了在其它包的代码中向该消息通道传消息，需要指定"hash"后存储到 map 中
	// 连接 SSE 传递消息，不能使用 POST 方法，需使用 GET，所以 hash 参数通过查询字串串传递
	// 参考：[How to send message to all connections in pool](https://stackoverflow.com/a/55208320)
	// [How to send message to all connections in pool](https://stackoverflow.com/questions/55207853)
	router.GET("/api/sse", mysse.Send)
	router.POST("/api/sse/tick", mysse.Tick)
	// 返回所有的消息通道
	router.GET("/api/sse/all", mysse.GetAll)

	// javlib
	router.POST("/api/fanhao/exist", javlib.ExistFanhaoFile)
	router.POST("/api/fanhao/subtitle", javlib.ExistSubtitle)
	// 重命名路径下的文件
	router.POST("/api/fanhao/rename", javlib.RenameDir)

	logger.Info.Printf("开始本地服务：http://127.0.0.1:%d\n", port)
	logger.Error.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// 显示systray托盘
func onReady() {
	systray.SetIcon(icons.Tray)
	systray.SetTitle("手机与 PC 传递数据")
	systray.SetTooltip("手机与 PC 传递数据")

	// mMatchQL := systray.AddMenuItem("运行青龙脚本", "运行青龙脚本")
	mOpenLog := systray.AddMenuItem("打开日志", "打开日志文件")
	systray.AddSeparator()
	mShutdown := systray.AddMenuItemCheckbox("自动关机", "完成所有任务后将自动关闭计算机", false)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	// 每分钟检查是否有任务，没有将自动关机
	taskTicker := time.NewTicker(1 * time.Minute)
	taskTicker.Stop()
	go func() {
		for range taskTicker.C {
			// 判断任务

			// 所有任务经过判断，已完成，关机
			if err := exec.Command("cmd", "/C", "shutdown", "/s", "/t", "60").Run(); err != nil {
				fmt.Println("执行关机出错：", err)
				return
			}
		}
	}()

	for {
		select {
		case <-mOpenLog.ClickedCh:
			err := dofile.OpenAs(logger.LogName)
			if err != nil {
				logger.Error.Printf("打开日志文件(%s)出错：%s\n", logger.LogName, err)
			}

			// 运行青龙脚本
		// case <-mMatchQL.ClickedCh:
		// 	_, err := ql.StartCommCronsCall()
		// 	if err != nil {
		// 		logger.Error.Printf("执行定时任务时出错：%s\n", err)
		// 		return
		// 	}
		case <-mShutdown.ClickedCh:
			if mShutdown.Checked() {
				taskTicker.Reset(1 * time.Minute)
			} else {
				taskTicker.Stop()
			}

		case <-mQuit.ClickedCh:
			// 退出程序
			logger.Warn.Println("退出程序")
			systray.Quit()
			os.Exit(0)
		}
	}
}
