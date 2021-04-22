package main

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"github.com/skip2/go-qrcode"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"path/filepath"
	"pc-phone-conn-go/logger"
	"strings"
	"time"
)

// 手机端发送的文本数据小于此值才复制到剪贴板，大于则保存到文件
const clipMAXSIZE = 512

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
// POST 参数：
// type: 操作的类型
//   可为"getclip"：手机端获取PC端的剪贴板数据
//   可为"URL"、"文本"：手机端发送数据到PC端
// content: 手机端传输的数据
//   可为文本、文件
func pcHander(c *gin.Context) {
	// 读取数据类型、数据内容
	ctype := c.PostForm("type")
	logger.Info.Printf("收到 '%s' 类型的请求或数据", ctype)

	// 手机端获取PC端的剪贴板数据的操作
	if ctype == "getclip" {
		text, err := clipboard.ReadAll()
		if err != nil {
			logger.Error.Printf("读取剪贴板出错：%v\n", err)
			c.String(http.StatusOK, "读取剪贴板出错："+err.Error())
			return
		}
		c.String(http.StatusOK, text)
		return
	}

	// 手机端发送数据到PC端的操作
	// 读取传输的数据
	file, header, err := c.Request.FormFile("content")
	if err != nil {
		logger.Error.Printf("PC端接收数据出错：%v\n", err)
		c.String(http.StatusOK, "PC端接收数据出错："+err.Error())
		return
	}
	//logger.Error.Printf("PC端接收数据大小：%d 字节\n", header.Size)

	// 返回给手机端显示的额外的信息（可空）
	extraInfo := ""
	// 根据数据类型分别处理
	if ctype == "URL" {
		url, err := readText(file)
		if err != nil {
			logger.Error.Printf("读取URL数据出错：%v\n", err)
			c.String(http.StatusOK, "读取URL数据出错："+err.Error())
			return
		}
		err = browser.OpenURL(url)
	} else if ctype == "文本" && header.Size <= clipMAXSIZE {
		text, err := readText(file)
		if err != nil {
			logger.Error.Printf("读取文本数据出错：%v\n", err)
			c.String(http.StatusOK, "读取文本数据出错："+err.Error())
			return
		}
		err = clipboard.WriteAll(text)
	} else {
		filename := getFilename(header)
		path, _ := filepath.Abs(filepath.Join(FileDir(), dofile.ValidFileName(filename, "_")))
		logger.Info.Printf("收到 '%s' 类型的数据，保存到 '%s'\n", ctype, path)
		err = c.SaveUploadedFile(header, path)
	}

	// 操作有错误
	if err != nil {
		logger.Error.Printf("执行 '%s' 类型的操作时出错：%v\n", ctype, err)
		c.String(http.StatusOK, fmt.Sprintf("执行 '%s' 类型的操作时出错：%v", ctype, err))
		return
	}

	// 正常完成
	if ctype == "文本" && header.Size > clipMAXSIZE {
		extraInfo = "，已作为文件保存"
	}
	c.String(http.StatusOK, fmt.Sprintf("执行 '%s' 类型的操作完成%s", ctype, extraInfo))
}

// 根据当前时间、请求头中的文件名 生成保存文件的文件名
func getFilename(header *multipart.FileHeader) string {
	filename := fmt.Sprintf("PPC-%d", time.Now().Unix())
	if header != nil && strings.TrimSpace(header.Filename) != "" {
		filename += "_" + header.Filename
	}
	return filename
}

// 读取流的文本数据
func readText(file multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)
	return buf.String(), err
}
