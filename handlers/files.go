package handlers

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"pc-phone-go/entity"
	"pc-phone-go/funcs"
	"pc-phone-go/funcs/logger"
)

// SendFiles 手机发送文件到 PC
//
// POST /api/files/send
//
// 请求头的`Content-Type`为`multipart/form-data`
//
// 支持发送多个文件，注意表单中文件的键都必须设为`file`
//
// 需要传递每个文件的文件名
func SendFiles(c *gin.Context) {
	const keyFile = "file"

	// 记录保存文件的失败记录，返回信息
	saveResult := make(map[string]string)

	form, err := c.MultipartForm()
	if err != nil {
		logger.Error.Printf("%s 解析多部分表单出错：%s\n", tagSendFiles, err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 30000,
			Msg:  fmt.Sprintf("%s 解析多部分表单出错", tagSendFiles),
			Data: err.Error(),
		})
		return
	}

	// 读取文件，支持读取多个文件
	files := form.File[keyFile]
	if len(files) == 0 {
		logger.Error.Printf("%s 无法读取到文件，检查多部分表单中文件的 key 是否为 '%s'\n", tagSendFiles, keyFile)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 30100,
			Msg:  fmt.Sprintf("%s 无法读取到文件，检查多部分表单中文件的 key 是否为 '%s'", tagSendFiles, keyFile),
			Data: nil,
		})
		return
	}

	// 保存文件
	for _, file := range files {
		filename := genFilename(file.Filename)
		path, _ := filepath.Abs(filepath.Join(funcs.FileDir(), dofile.ValidFileName(filename, "_")))

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			saveResult[file.Filename] = fmt.Sprintf("保存出错：%s", err)
		} else {
			saveResult[file.Filename] = fmt.Sprintf("保存成功：%s", path)
		}
	}

	// 返回结果以便显示
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  "保存文件的结果",
		Data: saveResult,
	})
}
