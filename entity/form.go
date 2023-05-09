package entity

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pc-phone-go/funcs/logger"
)

// OpForm 操作的表单
type OpForm[T any] struct {
	// 操作类型。如 "shutdown"
	Op string `json:"op"`

	// 传递的数据
	Data T `json:"data"`
}

// ParseForm 解析 JSON 表单参数
func ParseForm[T any](c *gin.Context, tag string) (*OpForm[T], *Rest) {
	var form OpForm[T]
	err := c.Bind(&form)
	if err != nil {
		logger.Error.Printf("%s 解析 JSON 表单出错：%s\n", tag, err)
		return nil, &Rest{
			Code: 10000,
			Msg:  fmt.Sprintf("%s 解析 JSON 表单出错", tag),
			Data: err.Error(),
		}
	}

	logger.Info.Printf("%s 收到操作：%+v\n", tag, form)

	return &form, nil
}
