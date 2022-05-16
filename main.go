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
	"runtime"
)

const (
	// 端口和路径：http://host:8899/topc
	port = "8899"
	path = "topc"
)

func init() {
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

	// 处理请求
	router.GET("/"+path, func(c *gin.Context) {
		c.String(http.StatusOK, "请使用 POST 方式访问")
	})
	router.POST("/"+path, pcHander)

	logger.Info.Println("开始本地服务：http://127.0.0.1:" + port)
	logger.Info.Println("本地服务地址：http://127.0.0.1:" + port + "/" + path)
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
		case <-mQuit.ClickedCh:
			// 退出程序
			logger.Info.Println("退出程序")
			systray.Quit()
			os.Exit(0)
		}
	}
}
