package main

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/icons"
	"pc-phone-go/tools/lives"
	"pc-phone-go/tools/pics"
	"pc-phone-go/tools/pics/pworker"
	"pc-phone-go/tools/ql"
	"pc-phone-go/tools/qx"
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

	// 剪贴板
	router.POST("/api/clip", handerClip)

	// 图片下载发送
	router.POST("/api/pics/dl", pics.Donwload)
	router.GET("/api/pics/dl/status", pics.Status)
	router.POST("/api/pics/dl/retry", pics.Retry)
	router.GET("/api/pics/dl/count", pics.Count)
	router.GET("/api/pics/dl/faillist", pics.FailList)
	router.GET("/api/pics/dl/skiplist", pics.SkipList)
	router.POST("/api/pics/del/yikeall", pics.DelYikeAll)

	// qx
	router.GET("/api/qx/parse_surge", qx.ParseSurge)

	// ql
	router.POST("/api/ql/set_env", ql.SetEnv)
	router.POST("/api/ql/start_comm_crons", ql.StartCommCrons)

	// lives
	router.GET("/api/lives/douyin/live", lives.GetDouyinRoom)

	logger.Info.Printf("开始本地服务：http://127.0.0.1:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// 显示systray托盘
func onReady() {
	systray.SetIcon(icons.Tray)
	systray.SetTitle("手机与 PC 传递数据")
	systray.SetTooltip("手机与 PC 传递数据")

	mMatchQL := systray.AddMenuItem("运行青龙脚本", "运行青龙脚本")
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
			// 判断图集任务
			hasPicsTask := false
			pworker.MapPWorker.Range(func(_, _ interface{}) bool {
				hasPicsTask = true
				return false
			})
			if hasPicsTask {
				break
			}

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
		case <-mMatchQL.ClickedCh:
			_, err := ql.StartCommCronsCall()
			if err != nil {
				logger.Error.Printf("执行定时任务时出错：%s\n", err)
				return
			}
		case <-mShutdown.ClickedCh:
			if mShutdown.Checked() {
				taskTicker.Reset(1 * time.Minute)
			} else {
				taskTicker.Stop()
			}

		case <-mQuit.ClickedCh:
			// 退出程序
			logger.Info.Println("退出程序")
			systray.Quit()
			os.Exit(0)
		}
	}
}
