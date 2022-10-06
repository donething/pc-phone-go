package qx

import (
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pc-phone-conn-go/conf"
	"pc-phone-conn-go/entity"
	"pc-phone-conn-go/logger"
	"strings"
	"time"
)

// ParseSurge 解析 surge 分流规则
func ParseSurge(c *gin.Context) {
	var client = dohttp.New(120*time.Second, false, false)
	if conf.Conf.Proxy != "" {
		err := client.SetProxy(conf.Conf.Proxy, nil)
		if err != nil {
			logger.Error.Printf("解析surge分流规则出错，设置网络代理时出错：%s\n", err)
			c.JSON(http.StatusOK, entity.Rest{Errcode: 5000,
				Msg: fmt.Sprintf("设置网络代理时出错")})
			return
		}
	}

	// 获取分流配置文件内容
	var url = c.Query("url")
	resp, err := http.Get(url)
	if err != nil {
		logger.Error.Printf("解析surge分流规则出错，请求出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Errcode: 5100,
			Msg: fmt.Sprintf("请求出错")})
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("解析surge分流规则出错，无法读取响应内容：'%s'\n", url)
		c.JSON(http.StatusOK, entity.Rest{Errcode: 5200,
			Msg: fmt.Sprintf("解析surge分流规则出错，无法读取响应内容：'%s'", url)})
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
