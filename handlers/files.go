package handlers

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pc-phone-go/entity"
	"pc-phone-go/funcs"
	"pc-phone-go/funcs/logger"
	"strings"
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
		filenameSaved = dofile.ValidFileName(filenameSaved, "_")
		path, _ := filepath.Abs(filepath.Join(funcs.FileDir(), filenameSaved))

		key := fmt.Sprintf("'%s'", filenameSaved)
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			errCount++
			saveResult[key] = fmt.Sprintf("保存出错：%s", err)
		}
		saveResult[key] = fmt.Sprintf("已保存到: '%s'", path)

		// 如果文件没有格式后缀，就尽量获取后增加
		if strings.Contains(filenameSaved, ".") {
			continue
		}
		// 增加格式
		err = appendExt(path)
		if err != nil {
			logger.Error.Printf("增加格式出错'%s'：%s\n", path, err)
		}
	}

	// 返回结果以便显示
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  fmt.Sprintf("共发送 %d 个文件，失败 %d 个", len(saveResult), errCount),
		Data: saveResult,
	})
}

// DownloadFile 手机获取 PC 的文件
//
// GET /api/file/download
//
// 查询字符串参数 path 需要传递完整路径
func DownloadFile(c *gin.Context) {
	path := c.Query("path")
	state, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		logger.Info.Printf("%s 文件不存在或判断出现错误 '%s'。'err'为'%v'\n", tagGetFile, err, path)
		c.JSON(http.StatusOK, entity.Rest{Code: 30200, Msg: "文件不存在或判断出现错误"})
		return
	}

	if state.IsDir() {
		logger.Info.Printf("%s 路径指向为目录 '%s'\n", tagGetFile, path)
		c.JSON(http.StatusOK, entity.Rest{Code: 30210, Msg: "路径指向为目录"})
		return
	}

	// 获取文件的名称
	// 不要用 path.Base()，path 包不跨平台
	fileName := filepath.Base(path)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	c.File(path)
}

// ListPath 返回指定路径目录下的文件信息
//
// GET /api/file/list
//
// 查询字符串参数 path 需要传递完整路径
//
// 返回 Rest<[]FileInfo>
func ListPath(c *gin.Context) {
	path := c.Query("path")
	state, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		logger.Info.Printf("%s 文件不存在或判断出现错误 '%s'。'err'为'%v'\n", taglistPath, err, path)
		c.JSON(http.StatusOK, entity.Rest{Code: 30300, Msg: "文件不存在或判断出现错误"})
		return
	}

	if !state.IsDir() {
		logger.Info.Printf("%s 路径指向为文件 '%s'\n", taglistPath, path)
		c.JSON(http.StatusOK, entity.Rest{Code: 30310, Msg: "路径指向为文件"})
		return
	}

	// 读取目录
	files, err := os.ReadDir(path)
	if err != nil {
		logger.Info.Printf("%s 读取目录'%s'出错：'%s'\n", taglistPath, path, err)
		c.JSON(http.StatusOK, entity.Rest{Code: 30320, Msg: "读取目录出错"})
		return
	}

	// 读取目录下所有子文件的信息
	payload := make([]entity.FileInfo, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			logger.Info.Printf("%s 获取文件信息出错 '%s'\n", taglistPath, err)
			c.JSON(http.StatusOK, entity.Rest{Code: 30330, Msg: "获取文件信息出错"})
			return
		}

		fileInfo := entity.FileInfo{
			Path:    filepath.Join(path, file.Name()),
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().UnixMilli(),
		}

		payload = append(payload, fileInfo)
	}

	// 返回
	logger.Info.Printf("%s 已读取目录下子文件的信息'%s'\n", taglistPath, path)
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "目录下子文件的信息", Data: payload})
}

// 增加格式
func appendExt(path string) error {
	// 读取文件前16字节，分析文件类型
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	buff := make([]byte, 16)
	_, err = io.ReadFull(f, buff)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	// 分析文件类型
	kind, err := filetype.Match(buff)
	if err != nil {
		return err
	}

	// 重命名
	return os.Rename(path, fmt.Sprintf("%s.%s", path, kind.Extension))
}
