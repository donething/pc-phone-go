package handlers

import (
	"github.com/atotto/clipboard"
	"pc-phone-go/funcs/logger"
	"strings"
	"testing"
)

func TestGetClip(t *testing.T) {
	text, err := clipboard.ReadAll()

	// 排除误报的错误
	if err != nil && !strings.Contains(err.Error(), "The operation completed successfully") {
		// 剪贴板为空
		if strings.Contains(err.Error(), "Element not found") {
			logger.Warn.Printf("%s PC 剪贴板为空：%s\n", tagGetClip, err)
			return
		}

		// 其它为真正的错误
		logger.Error.Printf("%s 读取 PC 的剪贴板出错：%s\n", tagGetClip, err)
		return
	}

	logger.Info.Printf("剪贴板：%s\n", text)
}
