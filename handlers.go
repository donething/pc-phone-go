// 手机与 PC 同步数据
// 功能包括 同步剪贴板、获取手机分到来的文本、文件
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
	"os/exec"
	"path/filepath"
	"pc-phone-go/funcs"
	"pc-phone-go/funcs/logger"
	"regexp"
	"strings"
	"time"
)

// 手机端发送的文本数据小于此值才复制到剪贴板，大于则保存到文件
const clipMAXSIZE = 512

// 首页，显示二维码
func index(c *gin.Context) {
	conn, err := net.Dial("ip:icmp", "google.com")
	funcs.CheckErr(err)
	localIP := conn.LocalAddr().String()

	qr, err := qrcode.New(fmt.Sprintf("http://%s:%d", localIP, port), qrcode.Medium)
	funcs.CheckErr(err)
	png, err := qr.PNG(256)
	funcs.CheckErr(err)
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(png)))
	_, err = c.Writer.Write(png)
	funcs.CheckErr(err)
}

// PC 端数据处理
// POST 参数：
// op: 操作的类型
//
//	为"getclip"：手机获取 PC 剪贴板数据
//	为"shutdown"：关闭 PC
//	为"shutdown_cancel"：取消关闭 PC
//	可为空：当手机发送数据到 PC 时，不需该参数
//
// data: 手机传给 PC 的数据
//
//	可为文本、文件等，可空
func handerClip(c *gin.Context) {
	// 读取请求中的参数：操作、数据
	op := c.PostForm("op")
	logger.Info.Printf("收到操作'%s'\n", op)

	// 手机端获取PC端的剪贴板数据的操作
	if op == "getclip" {
		text, err := clipboard.ReadAll()
		if err != nil {
			logger.Error.Printf("读取剪贴板出错：%v\n", err)
			c.String(http.StatusOK, "读取剪贴板出错："+err.Error())
			return
		}
		c.String(http.StatusOK, text)
		return
	}

	// 关闭 PC、取消关闭
	if op == "shutdown" || op == "shutdown_cancel" {
		var args []string
		var tips string
		if op == "shutdown" {
			args = []string{"-s", "-t", "60"}
			tips = "一分钟后将关闭 PC"
		} else if op == "shutdown_cancel" {
			args = []string{"-a"}
			tips = "取消关闭 PC"
		}

		cmd := exec.Command("shutdown", args...)
		_, err := cmd.CombinedOutput()
		if err != nil {
			logger.Error.Printf("执行(取消)关机命令时出错：%v\n", err)
			c.String(http.StatusOK, "执行(取消)关机命令时出错："+err.Error())
			return
		}
		logger.Info.Printf("已执行命令：%s\n", tips)
		c.String(http.StatusOK, "已执行命令："+tips)
		return
	}

	// 手机端发送到PC端的数据
	// 读取传输的数据
	file, header, err := c.Request.FormFile("data")
	if err != nil {
		logger.Error.Printf("PC端接收数据出错：%v\n", err)
		c.String(http.StatusOK, "PC端接收数据出错："+err.Error())
		return
	}
	// logger.Error.Printf("PC端接收数据大小：%d 字节\n", header.Size)
	// 返回给手机端以供显示的信息
	feedback := ""
	// 小文本直接复制到剪贴板（为链接时自动用浏览器打开），大文本则保存到文件中
	if header.Size <= clipMAXSIZE {
		text, err := readStreamText(file)
		if err != nil {
			logger.Error.Printf("读取文本数据出错：%v\n", err)
			c.String(http.StatusOK, "读取文本数据出错："+err.Error())
			return
		}
		// 识别网址，当去除首尾空格后为"https?://..."的字符串时，自动用浏览器打开
		urlReg := regexp.MustCompile(`^https?://\S+$`)
		if urlReg.MatchString(strings.TrimSpace(text)) {
			err = browser.OpenURL(strings.TrimSpace(text))
			feedback = "收到网址，已用浏览器打开"
		} else {
			err = clipboard.WriteAll(text)
			feedback = "收到短文本，已复制到剪贴板"
		}
	} else {
		filename := getFilename(header)
		path, _ := filepath.Abs(filepath.Join(funcs.FileDir(), dofile.ValidFileName(filename, "_")))
		logger.Info.Printf("收到大文件，将保存到：'%s'\n", op, path)
		err = c.SaveUploadedFile(header, path)
		feedback = "收到大文本或文件，已作为文件保存"
	}

	// 操作有错误
	if err != nil {
		logger.Error.Printf("执行操作'%s'时出错：%v\n", op, err)
		c.String(http.StatusOK, fmt.Sprintf("执行操作'%s'时出错：%v", op, err))
		return
	}

	// 返回结果以便显示
	c.String(http.StatusOK, fmt.Sprintf("%s", feedback))
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
func readStreamText(file multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)
	return buf.String(), err
}
