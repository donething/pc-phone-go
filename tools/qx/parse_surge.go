package qx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pc-phone-conn-go/entity"
	"pc-phone-conn-go/logger"
	"strings"
)

// ParseSurge 解析 surge 分流规则
func ParseSurge(c *gin.Context) {
	// 获取分流配置文件内容
	var url = c.Query("url")
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("解析surge分流规则失败，URL有误：'%s'", url))
		c.JSON(http.StatusOK, entity.Rest{Errcode: 5000,
			Msg: fmt.Sprintf("解析surge分流规则失败，URL有误：'%s'", url)})
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(fmt.Sprintf("解析surge分流规则失败，无法获取内容：'%s'", url))
		c.JSON(http.StatusOK, entity.Rest{Errcode: 5000,
			Msg: fmt.Sprintf("解析surge分流规则失败，无法获取内容：'%s'", url)})
		return
	}

	// 转换
	lines := strings.Split(string(bs), "\n")
	payload := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 以"#"开头的注释，直接追加到数组，其它为分流规则，在末尾添加规则后再追加
		newLine := line
		if strings.Index(line, "#") != 0 {
			newLine = fmt.Sprintf("%s,reject", line)
		}
		payload = append(payload, newLine)
	}

	// 成功返回
	logger.Info.Println("解析surge分流规则成功，将返回到客户端")
	c.String(http.StatusOK, strings.Join(payload, "\n"))
}
