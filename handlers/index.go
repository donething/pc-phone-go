package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"net"
	"pc-phone-go/comm"
	"pc-phone-go/funcs"
)

// Index 首页
//
// 显示二维码
func Index(c *gin.Context) {
	conn, err := net.Dial("ip:icmp", "google.com")
	funcs.CheckErr(err)
	localIP := conn.LocalAddr().String()

	qr, err := qrcode.New(fmt.Sprintf("http://%s:%d", localIP, comm.Port), qrcode.Medium)
	funcs.CheckErr(err)
	png, err := qr.PNG(256)
	funcs.CheckErr(err)
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(png)))
	_, err = c.Writer.Write(png)
	funcs.CheckErr(err)
}
