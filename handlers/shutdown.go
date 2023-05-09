package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
)

// 标签
const (
	tagShutdown = "[Shutdown]"
)

// 操作的常量值
const (
	opShutdown       = "shutdown"
	opCancelShutdown = "cancel"
)

// Shutdown 关闭 PC 及取消
//
// POST /api/shutdown
//
// 表单：entity.OpForm
//
// 参数 Op 为 op* 常量中的值
//
// 参数 Data 为等待的秒数（为 0 时默认为 60）。仅当操作为关机时传递
func Shutdown(c *gin.Context) {
	// 提取表单
	form, rest := entity.ParseForm[int](c, tagShutdown)
	if rest != nil {
		c.JSON(http.StatusOK, rest)
		return
	}

	// 默认 60 秒关机
	if form.Data == 0 {
		form.Data = 60
	}

	// 根据操作参数，执行命令
	var args []string
	var tips string
	if form.Op == opShutdown {
		args = []string{"-s", "-t", fmt.Sprintf("%d", form.Data)}
		tips = fmt.Sprintf("%d 秒后将关闭 PC", form.Data)
	} else if form.Op == opCancelShutdown {
		args = []string{"-a"}
		tips = "取消关闭 PC"
	} else {
		logger.Error.Printf("%s 未知的操作：'%s'\n", tagShutdown, form.Op)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 10100,
			Msg:  fmt.Sprintf("%s 未知的操作：'%s'\n", tagShutdown, form.Op),
			Data: nil,
		})
		return
	}

	cmd := exec.Command("shutdown", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Printf("%s 执行/取消 关机命令时出错：%s\n", tagShutdown, err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 10200,
			Msg:  fmt.Sprintf("%s 执行/取消 关机命令时出错", tagShutdown),
			Data: err.Error(),
		})
		return
	}

	// 正确执行
	logger.Info.Printf("%s 已执行命令：%s\n", tagShutdown, tips)
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  fmt.Sprintf("%s 已执行命令：%s", tagShutdown, tips),
		Data: nil,
	})
}
