package qx

import (
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pc-phone-go/conf"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"strings"
)

// ParseSurge 解析 surge 分流规则
func ParseSurge(c *gin.Context) {
	var client = dohttp.New(false, false)
	if conf.Conf.Comm.Proxy != "" {
		err := client.SetProxy(conf.Conf.Comm.Proxy)
		if err != nil {
			msg := fmt.Sprintf("解析surge分流规则出错，设置网络代理时出错：%s", err)
			logger.Error.Println(msg)
			c.JSON(http.StatusOK, entity.Rest{Code: 5000, Msg: msg})
			return
		}
	}

	// 获取分流配置文件内容
	var url = c.Query("url")
	resp, err := http.Get(url)
	if err != nil {
		msg := fmt.Sprintf("解析surge分流规则出错，请求出错：%s", err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 5100, Msg: msg})
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("解析surge分流规则出错，无法读取响应内容：'%s'", url)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{Code: 5200, Msg: msg})
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
