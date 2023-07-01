package handlers

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"net/http"
	"path"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"regexp"
	"strings"
	"time"
)

// 标签
const (
	tagGetClip   = "[GetClip]"
	tagSendText  = "[SendText]"
	tagSendFiles = "[SendFiles]"
)

// GetClip 手机读取 PC 的剪贴板
//
// 当剪贴板中的文本是有效路径时，手机端将获取到该文件
//
// GET /api/clip/get
func GetClip(c *gin.Context) {
	// 读取剪贴板
	text, err := clipboard.ReadAll()

	// 排除误报的错误
	if err != nil && !strings.Contains(err.Error(), "The operation completed successfully") {
		// 剪贴板为空
		if strings.Contains(err.Error(), "Element not found") {
			logger.Warn.Printf("%s PC 剪贴板为空：%s\n", tagGetClip, err)
			c.JSON(http.StatusOK, entity.Rest{
				Code: 20000,
				Msg:  fmt.Sprintf("%s PC 剪贴板为空：%s", tagGetClip, err),
			})
			return
		}

		// 其它为真正的错误
		logger.Error.Printf("%s 读取 PC 的剪贴板出错：%s\n", tagGetClip, err)
		c.JSON(http.StatusOK, entity.Rest{Code: 20010, Msg: fmt.Sprintf("%s ：%s", tagGetClip, err)})
		return
	}

	// 分析剪贴板的文本是否为文件路径
	// 判断出错或文件不存在，当文本发送
	exist, err := dofile.Exists(text)
	if err != nil || !exist {
		logger.Info.Printf("%s 作为文本发送：%s\n", tagGetClip, text)
		c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: text})
		return
	}

	// 当文件发送
	filePath := text
	logger.Info.Printf("%s 作为文件发送：%s\n", tagGetClip, filePath)

	// 获取文件的名称
	fileName := path.Base(filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	c.File(filePath)
}

// SendText 手机发送数据到 PC 的剪贴板
//
// POST /api/clip/send
//
// 表单：entity.OpForm
//
// 参数 Op 为空
//
// 参数 Data 为发送的文本数据
func SendText(c *gin.Context) {
	// 提取表单
	form, errRest := entity.ParseForm[string](c, tagSendText)
	if errRest != nil {
		c.JSON(http.StatusOK, errRest)
		return
	}

	// 发送的文本数据
	text := form.Data
	if strings.TrimSpace(text) == "" {
		logger.Error.Printf("%s 文本数据为空\n", tagSendText)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 20100,
			Msg:  fmt.Sprintf("%s 文本数据为空", tagSendText),
		})
		return
	}

	// 返回给手机端以供显示的信息
	var feedback string
	var err error

	// 识别网址，当去除首尾空格后为"https?://..."的字符串时，自动用浏览器打开
	urlReg := regexp.MustCompile(`^https?://\S+$`)
	if urlReg.MatchString(strings.TrimSpace(text)) {
		err = browser.OpenURL(strings.TrimSpace(text))
		feedback = "PC 收到网址，已用浏览器打开"
	} else {
		err = clipboard.WriteAll(text)
		feedback = "PC 收到文本，已复制到剪贴板"
	}

	// 操作有错误
	if err != nil {
		logger.Error.Printf("%s 执行操作时出错：%s\n", tagSendText, err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 20200,
			Msg:  fmt.Sprintf("%s 执行操作时出错：%s", tagSendText, err),
		})
		return
	}

	// 返回结果以便显示
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  feedback,
	})
}

// 根据当前时间、请求头中的文件名 生成保存文件的文件名
func genFilename(name string) string {
	filename := fmt.Sprintf("PP-%d", time.Now().Unix())
	if strings.TrimSpace(name) != "" {
		filename += "_" + name
	}

	return filename
}
