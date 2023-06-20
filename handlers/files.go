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
		msg := fmt.Sprintf("%s 解析多部分表单出错：%s", tagSendFiles, err)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 30000,
			Msg:  msg,
		})
		return
	}

	// 读取文件，支持读取多个文件
	files := form.File[keyFile]
	if len(files) == 0 {
		msg := fmt.Sprintf("%s 无法读取到文件，检查多部分表单中文件的 key 是否为 '%s'", tagSendFiles, keyFile)
		logger.Error.Println(msg)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 30100,
			Msg:  msg,
		})
		return
	}

	// 保存文件
	var errCount = 0
	for _, file := range files {
		filenameSaved := genFilename(file.Filename)
		path, _ := filepath.Abs(filepath.Join(funcs.FileDir(), dofile.ValidFileName(filenameSaved, "_")))

		err = c.SaveUploadedFile(file, path)
		key := fmt.Sprintf("'%s'", file.Filename)
		if err != nil {
			errCount++
			saveResult[key] = fmt.Sprintf("保存出错：%s", err)
		} else {
			saveResult[key] = fmt.Sprintf("保存到：'%s'", path)
		}
	}

	// 返回结果以便显示
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  fmt.Sprintf("共发送 %d 个文件，失败 %d 个", len(saveResult), errCount),
		Data: saveResult,
	})
}
