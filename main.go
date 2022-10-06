package main

import (
	"github.com/donething/utils-go/dofile"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"pc-phone-conn-go/icons"
	"pc-phone-conn-go/logger"
	"pc-phone-conn-go/tools/ql"
	"pc-phone-conn-go/tools/qx"
	"runtime"
	"time"
)

const (
	// 服务端口
	port = "8800"
)

func init() {
	go func() {
		// 在本应用运行后需等一段时间，等 Docker 启动目标容器后才执行脚本，用于电脑刚开机时
		t := time.NewTimer(3 * time.Minute)
		<-t.C

		num, err := ql.StartCommCronsCall()
		if err != nil {
			logger.Error.Printf("执行定时任务时出错：%s\n", err)
			return
		}
		logger.Info.Printf("已发送执行定时任务的请求，共计 %d 个任务\n", num)
	}()

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

	// 剪贴板
	router.POST("/api/clip", handerClip)

	// qx
	router.GET("/api/qx/parse_surge", qx.ParseSurge)

	// ql
	router.POST("/api/ql/set_env", ql.SetEnv)
	router.POST("/api/ql/start_comm_crons", ql.StartCommCrons)

	logger.Info.Println("开始本地服务：http://127.0.0.1:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func CheckErr(err error) {
	if err != nil {
		logger.Error.Panicln(err)
	}
}

// 显示systray托盘
func onReady() {
	systray.SetIcon(icons.Tray)
	systray.SetTitle("手机与 PC 传递数据")
	systray.SetTooltip("手机与 PC 传递数据")

	mMatchQL := systray.AddMenuItem("运行青龙脚本", "运行青龙脚本")
	mOpenLog := systray.AddMenuItem("打开日志", "打开日志文件")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	for {
		select {
		case <-mOpenLog.ClickedCh:
			err := dofile.OpenAs(logger.LogName)
			if err != nil {
				logger.Error.Printf("打开日志文件(%s)出错：%s\n", logger.LogName, err)
			}

			// 运行青龙脚本
		case <-mMatchQL.ClickedCh:
			num, err := ql.StartCommCronsCall()
			if err != nil {
				logger.Error.Printf("执行定时任务时出错：%s\n", err)
				return
			}
			logger.Info.Printf("已发送执行定时任务的请求，共计 %d 个任务\n", num)

		case <-mQuit.ClickedCh:
			// 退出程序
			logger.Info.Println("退出程序")
			systray.Quit()
			os.Exit(0)
		}
	}
}
