package pworker

import (
	"encoding/json"
	"errors"
	"fmt"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/funcs/notify"
	"pc-phone-go/tools/pics/pcomm"
	"pc-phone-go/tools/pics/phandlers/tolocal"
	"pc-phone-go/tools/pics/phandlers/totg"
	"pc-phone-go/tools/pics/phandlers/toyike"
	"strings"
	"sync"
)

// PWorkers 可多实例的工作对象
type PWorkers struct {
	// 当前任务的图集数量
	Total int

	// TasksCh 任务通道
	TasksCh chan pcomm.Album
	WGTask  sync.WaitGroup
	// 通道是否已被关闭
	IsClosed bool

	// 已成功发送任务的数量
	DoneCount int
	MuDone    sync.RWMutex

	// 已跳过图集（广告、低质量）的数量
	SkipCount int
	MuSkip    sync.RWMutex

	// 发送失败图集的数量
	FailCount int
	MuFail    sync.RWMutex
}

// MapPWorker 所有的PWorker
var MapPWorker sync.Map

// Init 初始化工作协程
func (w *PWorkers) Init(workerCount int) {
	// 创建一个有缓冲的通道来管理工作
	w.TasksCh = make(chan pcomm.Album, workerCount)
	w.IsClosed = false

	// 启动 goroutine 来完成工作
	for id := 1; id <= workerCount; id++ {
		go w.work(id)
	}
	logger.Info.Println("[PWorker] 工作协程已准备就绪")
}

// 工作
func (w *PWorkers) work(id int) {
	defer w.WGTask.Done()

	for {
		// 等待分配工作
		task, ok := <-w.TasksCh
		if !ok {
			// 这意味着通道已经空了，并且已被关闭
			logger.Info.Printf("[Worker%02d] 通道已关闭，完成任务\n", id)
			w.IsClosed = true
			return
		}

		var err error
		switch Conf.Pics.Handler {
		// 发送到一刻相册
		case pcomm.HandlerToYike:
			err = toyike.Send(task)
		case pcomm.HandlerToLocal:
			err = tolocal.Save(task)
		case pcomm.HandlerToTG:
			err = totg.Send(task)
		default:
			logger.Error.Printf("[Worker] 未知的 Handler：'%s'\n", id, Conf.Pics.Handler)
			notify.WXPushCard("发送图片出错",
				fmt.Sprintf("未知的 Handler：'%s'", Conf.Pics.Handler), "", "")
			return
		}

		// 该图集在下载失败时的保存到数据桶中的键
		dbkey := fmt.Sprintf("%s_%s_%s", task.Plat, task.UID, task.ID)
		if err != nil {
			// 不是错误，而是跳过
			if errors.Is(err, toyike.ErrImgTooSmall) {
				logger.Warn.Printf("[Worker][%s] 跳过下载小图图集'%s'\n", task.Plat, task.ID)
				w.MuSkip.Lock()
				w.SkipCount++
				w.MuSkip.Unlock()

				// 保存到数据库
				taskBS, _ := json.Marshal(task)
				errSet := DB.Set([]byte(pcomm.DBSkip+dbkey), taskBS)
				if errSet != nil {
					logger.Error.Printf("[Worker][%s] 保存跳过下载的图集'%s'到数据库时出错：%s\n",
						task.Plat, task.ID, errSet)
					notify.WXPushCard("发送图片出错",
						fmt.Sprintf("保存跳过下载的图集到数据库出错：%s", errSet), "", "")
					continue
				}

				delLog(task, dbkey)
			} else {
				// 发生错误
				logger.Error.Printf("%s\n", err)
				w.MuFail.Lock()
				w.FailCount++
				w.MuFail.Unlock()
				// 保存到数据库
				taskBS, _ := json.Marshal(task)
				errSet := DB.Set([]byte(pcomm.DBFail+dbkey), taskBS)
				if errSet != nil {
					logger.Error.Printf("[Worker][%s] 保存下载失败的图集'%s'到数据库时出错：%s\n",
						task.Plat, task.ID, errSet)
					notify.WXPushCard("发送图片出错",
						fmt.Sprintf("保存下载失败的图集到数据库出错：%s", errSet), "", "")
				}
			}

			continue
		}

		logger.Info.Printf("[Worker][%s] 已完成下载、发送图集'%s'\n", task.Plat, task.ID)
		w.MuDone.Lock()
		w.DoneCount++
		w.MuDone.Unlock()

		delLog(task, dbkey)
	}
}

// 当重试成功后，从记录中删除
func delLog(task pcomm.Album, dbkey string) {
	// 为重试任务时，成功后需要删除该任务记录
	if strings.TrimSpace(task.RetryFrom) != "" {
		err := DB.Del([]byte(task.RetryFrom + dbkey))
		if err != nil {
			logger.Error.Printf("[Worker][%s] 删除数据库中图集'%s'的下载失败记录时出错：%s\n",
				task.Plat, task.ID, err)
			notify.WXPushCard("发送图片出错",
				fmt.Sprintf("删除数据库中图集下载失败的记录时出错：%s", err), "", "")
		}
	}
}
