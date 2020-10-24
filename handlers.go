package main

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

// 首页
func index(c *gin.Context) {
	conn, err := net.Dial("ip:icmp", "google.com")
	CheckErr(err)
	localIP := conn.LocalAddr().String()

	qr, err := qrcode.New(fmt.Sprintf("http://%s:%s/%s", localIP, port, path),
		qrcode.Medium)
	CheckErr(err)
	png, err := qr.PNG(256)
	CheckErr(err)
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(png)))
	_, err = c.Writer.Write(png)
	CheckErr(err)
}

// PC 端数据处理
func pcHander(c *gin.Context) {
	// 读取数据类型、数据内容
	ctype := c.PostForm("type")
	log.Printf("收到 '%s' 类型的数据\n", ctype)

	file, header, err := c.Request.FormFile("content")
	if err != nil && ctype != "getclip" {
		log.Printf("提取数据出错：%v\n", err)
		c.String(http.StatusOK, "提取数据出错，传递的数据可能为空："+err.Error())
		return
	}

	// 根据数据类型分别处理
	switch ctype {
	case "getclip":
		text, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("读取剪贴板出错：%v\n", err)
			c.String(http.StatusOK, "读取剪贴板出错："+err.Error())
			return
		}
		c.String(http.StatusOK, text)
	case "URL", "文本":
		// 先读取纯文本类型的数据
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			log.Printf("读取纯文本内容出错：%v\n", err)
			c.String(http.StatusOK, "读取纯文本内容出错："+err.Error())
			return
		}
		text := buf.String()

		if ctype == "URL" {
			// 链接
			err = browser.OpenURL(text)
		} else if ctype == "文本" {
			if header.Size <= 512 {
				err = clipboard.WriteAll(text)
			} else {
				path, _ := filepath.Abs(filepath.Join(FileDir(), header.Filename))
				log.Printf("收到 '%s' 类型的数据，保存到 %s\n", ctype, path)
				err = c.SaveUploadedFile(header, path)
			}
		}
	default:
		path, _ := filepath.Abs(filepath.Join(FileDir(), header.Filename))
		log.Printf("收到 '%s' 类型的数据，保存到 %s\n", ctype, path)
		err = c.SaveUploadedFile(header, path)
	}

	// 操作有错误
	if err != nil {
		log.Printf("执行 '%s' 类型的操作时出错：%v\n", ctype, err)
		c.String(http.StatusOK, fmt.Sprintf("执行 '%s' 类型的操作时出错：%v", ctype, err))
		return
	}

	// 正常完成
	c.String(http.StatusOK, fmt.Sprintf("执行 '%s' 类型的操作完成", ctype))
}
