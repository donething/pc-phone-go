package funcs

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
)

// FileDir 文件保存的目录，默认为用户下载目录，无法获取则为当前运行目录下的"Downloads"文件夹
func FileDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "Downloads")
}

func CheckErr(err error) {
	if err != nil {
		logger.Error.Panicln(err)
	}
}

// OpenLocal 打开或显示本地文件
//
// 参数 POST 传递 json 数据，"method"的值为"open"或"show"，"path"为目标路径：{"method": string, path: string}
//
// 返回 {"code": int, "msg": string, "data": string}
func OpenLocal(c *gin.Context) {
	var params struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	}
	err := c.Bind(&params)
	if err != nil {
		text := fmt.Sprintf("绑定显示本地文件的参数时出错：%s", err)
		logger.Error.Println(text)
		c.JSON(http.StatusOK, entity.Rest{Code: 3040, Msg: text, Data: nil})
		return
	}
	switch params.Method {
	case "open":
		err = dofile.OpenAs(params.Path)
		if err != nil {
			logger.Error.Printf("打开本地文件出错：%s\n", err)
			c.JSON(http.StatusOK, entity.Rest{Code: 3050, Msg: "打开本地文件出错", Data: err.Error()})
		}
	case "show":
		err = dofile.ShowInExplorer(params.Path)
		if err != nil {
			logger.Error.Printf("显示本地文件出错：%s\n", err)
			c.JSON(http.StatusOK, entity.Rest{Code: 3060, Msg: "显示本地文件出错", Data: err.Error()})
		}
	}
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "已显示本地文件", Data: nil})
}
