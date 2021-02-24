package main

import (
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dolog"
	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"pc-phone-conn-go/icons"
)

const (
	// 端口和路径：http://192.168.1.52:8899/topc
	port = "8899"
	path = "topc"
)

var (
	logFile *os.File
)

func init() {
	var err error
	logFile, err = dolog.LogToFile(dolog.LogName, os.O_CREATE|os.O_APPEND, dolog.LogFormat)
	CheckErr(err)
	go systray.Run(onReady, nil)
}
func main() {
	defer logFile.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 显示 PC 端地址的二维码
	router.GET("/", index)

	// 处理请求
	router.GET("/"+path, func(c *gin.Context) {
		c.String(http.StatusOK, "请使用 POST 方式访问")
	})
	router.POST("/"+path, pcHander)

	log.Println("开始本地服务：http://127.0.0.1:" + port)
	log.Println("本地服务地址：http://127.0.0.1:" + port + "/" + path)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 显示systray托盘
func onReady() {
	systray.SetIcon(icons.Tray)
	systray.SetTitle("手机和 PC 传递数据")
	systray.SetTooltip("手机和 PC 传递数据")

	mOpenLog := systray.AddMenuItem("打开日志", "打开日志文件")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出程序")

	for {
		select {
		case <-mOpenLog.ClickedCh:
			err := dofile.OpenAs(dolog.LogName)
			if err != nil {
				log.Printf("打开日志文件(%s)出错：%s\n", dolog.LogName, err)
			}
		case <-mQuit.ClickedCh:
			// 退出程序
			log.Println("退出程序")
			systray.Quit()
			os.Exit(0)
		}
	}
}
