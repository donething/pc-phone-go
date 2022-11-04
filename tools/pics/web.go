package pics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-conn-go/conf"
	"pc-phone-conn-go/entity"
	"pc-phone-conn-go/funcs/database"
	"pc-phone-conn-go/funcs/logger"
	push "pc-phone-conn-go/funcs/notify"
	"pc-phone-conn-go/tools/pics/pcomm"
	"pc-phone-conn-go/tools/pics/phandlers/toyike"
	"pc-phone-conn-go/tools/pics/pworker"
	"strings"
	"time"
)

// 请求头
var wbHeaders = map[string]string{
	"Referer": "https://weibo.com",
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/91.0.4472.106 Safari/537.36",
}

// Donwload 下载图集
//
// POST /api/pics/dl
//
// POST JSON 数据，类型为 []Album
func Donwload(c *gin.Context) {
	var albums []pcomm.Album
	err := c.BindJSON(&albums)
	if err != nil {
		logger.Error.Printf("下载发送图片出错：绑定 JSON 数据出错：%s\n", err)
		push.WXPushCard("VPS 下载发送图片出错",
			fmt.Sprintf("绑定下载信息参数出错：%s", err), "", "")
		c.JSON(http.StatusOK, entity.Rest{Code: 5000, Msg: fmt.Sprintf("绑定 JSON 数据出错：%s", err)})
		return
	}

	go func() {
		err = start(albums)
		if err != nil {
			logger.Error.Printf("%s", err)
			push.WXPushCard("VPS 下载发送图片出错", err.Error(), "", "")
		}
	}()

	logger.Info.Printf("已提交 %d 个图集的下载任务", len(albums))
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: fmt.Sprintf("已提交 %d 个图集的下载任务", len(albums))})
	return
}

// Retry 重试下载失败的图集
//
// POST /api/pics/dl/retry
func Retry(c *gin.Context) {
	albums, err := database.GetAll(pcomm.DBFail)
	if err != nil {
		c.JSON(http.StatusOK, entity.Rest{Code: 5010, Msg: fmt.Sprintf("获取下载失败的图集时出错：%s", err)})
		return
	}
	if len(albums) == 0 {
		c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: fmt.Sprintf("没有下载失败的图集任务，不需重试")})
		return
	}

	go func() {
		err = start(albums)
		if err != nil {
			logger.Error.Printf("%s", err)
			push.WXPushCard("VPS 下载发送图片出错", err.Error(), "", "")
		}
	}()

	logger.Info.Printf("已提交 %d 个图集的重试下载任务", len(albums))
	c.JSON(http.StatusOK, entity.Rest{Code: 0,
		Msg: fmt.Sprintf("已提交 %d 个图集的重试下载任务", len(albums))})
	return
}

// FailList 已下载失败的图集列表
//
// GET /api/pics/dl/faillist
func FailList(c *gin.Context) {
	albums, err := database.GetAll(pcomm.DBFail)
	if err != nil {
		logger.Error.Printf("读取下载失败的图集时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Code: 5020, Msg: "读取下载失败的图集时出错：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "下载失败的图集", Data: albums})
}

// SkipList 已跳过下载的图集列表
//
// GET /api/pics/dl/skiplist
func SkipList(c *gin.Context) {
	albums, err := database.GetAll(pcomm.DBSkip)
	if err != nil {
		logger.Error.Printf("读取已跳过下载的图集时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Code: 5030, Msg: "读取已跳过下载的图集时出错：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "已跳过下载的图集", Data: albums})
}

// Status 图集下载的状态
//
// GET /api/pics/dl/status
func Status(c *gin.Context) {
	statusList := make(map[string]pcomm.PStatus)
	pworker.MapPWorker.Range(func(key, value interface{}) bool {
		if pw, ok := value.(*pworker.PWorkers); ok {
			// 不包含通道已关闭的任务（此时说明已完成任务，会通过发送通知消息）
			if pw.IsClosed {
				pworker.MapPWorker.Delete(key)
				return true
			}

			// 状态
			status := new(pcomm.PStatus)
			status.Total = pw.Total
			// 已成功发送数
			pw.MuDone.Lock()
			status.Done = pw.DoneCount
			pw.MuDone.Unlock()
			// 跳过数
			pw.MuSkip.Lock()
			status.Skip = pw.SkipCount
			pw.MuSkip.Unlock()
			// 失败数
			pw.MuFail.Lock()
			status.Fail = pw.FailCount
			pw.MuFail.Unlock()
			statusList[fmt.Sprintf("%s", key)] = *status
		}
		return true
	})

	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "图集下载的进度", Data: statusList})
	return
}

// Count 需要重试下载的图集数
//
// GET /api/pics/dl/count
func Count(c *gin.Context) {
	// 总失败数
	failAlbums, err := database.GetAll(pcomm.DBFail)
	if err != nil {
		logger.Error.Printf("读取需要重试下载的图集数时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Code: 5040, Msg: "读取需要重试下载的图集数时出错：" + err.Error()})
		return
	}
	// 总跳过数
	skipAlbums, err := database.GetAll(pcomm.DBSkip)
	if err != nil {
		logger.Error.Printf("读取需要跳过下载的图集数时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Code: 5050, Msg: "读取跳过重试下载的图集数时出错：" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "需要重试下载的图集数",
		Data: pcomm.PCount{Fail: len(failAlbums), Skip: len(skipAlbums)}})
}

// DelYikeAll 删除一刻相册中所有的图片
//
// Post /api/pics/del/yikeall
func DelYikeAll(c *gin.Context) {
	for {
		err := toyike.DelAll()
		if err != nil {
			if strings.Index(err.Error(), "50503") >= 0 {
				logger.Info.Printf("删除一刻相册中的图片出现一点错误，重试删除\n")
				continue
			} else {
				logger.Info.Printf("%s\n", err.Error())
				push.WXPushCard("删除一刻相册中所有的图片失败",
					fmt.Sprintf("删除一刻相册中所有的图片失败：%s", err), "", "")
				c.JSON(http.StatusOK, entity.Rest{Code: 5060, Msg: err.Error()})
				return
			}
		}
		break
	}

	logger.Info.Println("已尝试删除一刻相册的所有图片，若还有遗漏，可再次运行本程序")
	push.WXPushCard("删除一刻相册所有图片完成",
		"已尝试删除一刻相册的所有图片，若还有遗漏，可再次运行本程序", "", "")
	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "已尝试删除一刻相册的所有图片，若还有遗漏，可再次运行本程序"})
}

// 下载
func start(albums []pcomm.Album) error {
	logger.Info.Printf("开始下载 %d 个图集\n", len(albums))
	// 初始化下载、发送任务
	workerCount := conf.Conf.Pics.WorkerCount
	if workerCount == 0 {
		workerCount = 10
	}
	pworkers := pworker.PWorkers{}
	pworker.MapPWorker.Store(fmt.Sprintf("%d", time.Now().UnixMicro()), &pworkers)
	pworkers.Total = len(albums)
	pworkers.Init(workerCount)
	pworkers.WGTask.Add(workerCount)

	// 逆序发送任务，完成后关闭任务通道
	for i := len(albums) - 1; i >= 0; i-- {
		album := albums[i]
		// 跳过没有图集链接的任务
		if len(album.URLs) == 0 || len(album.URLsM) == 0 {
			continue
		}

		switch album.Plat {
		case pcomm.TagWeibo:
			album.Header = wbHeaders
		default:

			return fmt.Errorf("下载发送图片出错：未适配的平台：'%s'", album.Plat)
		}

		pworkers.TasksCh <- album
	}

	// 关闭任务通道
	close(pworkers.TasksCh)

	// 等待任务完成
	pworkers.WGTask.Wait()

	// 统计跳过
	var skipCount = 0
	pworkers.MuSkip.Lock()
	skipCount = pworkers.SkipCount
	pworkers.MuSkip.Unlock()
	// 统计失败
	var failCount = 0
	pworkers.MuFail.Lock()
	failCount = pworkers.FailCount
	pworkers.MuFail.Unlock()

	// 发送消息
	skipMsg := ""
	if skipCount > 0 {
		skipMsg = fmt.Sprintf("，已跳过 %d 个图集", skipCount)
	}
	failMsg := ""
	if failCount > 0 {
		failMsg = fmt.Sprintf("，下载失败 %d 个图集", failCount)
	}

	logger.Info.Printf("图集下载任务已完成：共有 %d 个图集%s%s", len(albums), skipMsg, failMsg)
	push.WXPushCard("VPS 图集下载任务已完成", fmt.Sprintf("图集下载任务已完成：共有 %d 个图集%s%s",
		len(albums), skipMsg, failMsg), "", "")
	return nil
}
