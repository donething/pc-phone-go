package javlib

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dotext"
	"github.com/gin-contrib/sse"
	"os"
	"path/filepath"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/funcs/mysse"
	"regexp"
	"strings"
)

const (
	// 未被屏蔽的网址，可在 javbus.com 找到
	avDetailURL = `https://www.javbus.com/%s`
	// 需要重命名的视频类型，以"."开始，以"|"分隔
	videoType = `.mp4|.mkv|.wmv|.avi|.ts`
	// 需要重命名的字幕文件类型，格式同上
	subtitleType = `.srt|.ass`
	// 根据文件名的长度，判断是否需要重命名（首次运行，尽量设为10000等大值；之后可设为60等小值）
	// needRenameLimit = 60
)

// Rename 重命名目录列表内的视频文件
//
// 该函数发送 SSEData 的格式为：Msg 为发送的信息（番号名、出错信息等），Data 为传输的数据（文件的路径）
func Rename(paths []string, eventCh chan sse.Event) {
	// 遍历目录列表
	for _, p := range paths {
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 为目录时，不需要重命名
			if info.IsDir() {
				// logger.Info.Printf("[%s]为目录\n", path)
				return nil
			}

			// 获取文件名
			name := info.Name()
			// 文件名够长则不需要重命名
			/*
				if len(info.Name()) > needRenameLimit {
					logger.Info.Printf("[%s]的文件名已经很长\n", path)
					return nil
				}
			*/

			// 获取文件的格式
			format := filepath.Ext(name)
			// 不是视频、字幕文件，不需要重命名
			if !strings.Contains(videoType, format) && !strings.Contains(subtitleType, format) {
				// logger.Info.Printf("[%s]为未知的文件类型\n", path)
				return nil
			}

			// 从文件名（去除格式后缀，以免该后缀被识别为番号）中解析出番号
			fanhao := dotext.ResolveFanhao(strings.TrimRight(name, format))
			// 如果没找到番号，就不需要重命名
			if fanhao == "" {
				logger.Error.Printf("无法解析番号：'%s'\n", path)
				mysse.SendToEventCh(eventCh, mysse.NewMsg(eventRename, false, "无法解析番号", path))
				return nil
			}
			// 联网获取番号的完整名（但不包括视频格式），如果不为""，则开始重命名
			fanhaoName, err := obtainFanhaoName(fanhao)
			if err != nil {
				logger.Error.Printf("查找番号'%s'出错：%s\n", fanhao, err)
				mysse.SendToEventCh(eventCh, mysse.NewMsg(eventRename, false,
					fmt.Sprintf("查找番号'%s'出错：'%s'", fanhao, err), path))
				return nil
			}

			// 如果番号以"-C"结尾，表示包含中文字幕，需要在文件名的番号后添加"-C"以体现
			if strings.Contains(strings.ToUpper(name), "-C") {
				fanhaoName = strings.Replace(fanhaoName, fanhao, fanhao+"-C", 1)
			}
			// 如果番号以"-4k"结尾，表示为4k分辨率，需要在文件名的番号后添加"-K"以体现
			if strings.Contains(strings.ToUpper(name), "-4K") {
				fanhaoName = strings.Replace(fanhaoName, fanhao, fanhao+"-4K", 1)
			}

			// 准备重命名
			fullname := fanhaoName + format
			newName := filepath.Join(filepath.Dir(path), dofile.ValidFileName(fullname, " "))
			err = os.Rename(path, newName)
			// 每完成一个发送通知
			// logger.Info.Printf("完成重命名为'%s'\n", newName)
			if err != nil {
				logger.Error.Printf("重命名文件'%s'出错：%s\n", path, err)
				mysse.SendToEventCh(eventCh, mysse.NewMsg(eventRename, false,
					fmt.Sprintf("重命名文件'%s'出错：'%s'", fanhao, err), path))
			} else {
				mysse.SendToEventCh(eventCh, mysse.NewMsg(eventRename, true, fanhao, newName))
			}
			// time.Sleep(1 * time.Second)
			return nil
		})

		if err != nil {
			logger.Error.Printf("重命名时遍历路径出错：%s\n", err)
			mysse.SendToEventCh(eventCh, mysse.NewMsg(eventRename, false,
				fmt.Sprintf("遍历路径出错：'%s'", err), p))
			continue
		}
	}
}

// 获取番号对应的全名
func obtainFanhaoName(fanhao string) (string, error) {
	// 联网获取包含番号的文本
	url := fmt.Sprintf(avDetailURL, fanhao)
	text, err := client.GetText(url, nil)
	if err != nil {
		return "", err
	}

	// 提取番号名
	p := regexp.MustCompile(`<h3>(.+)</h3>`)
	m := p.FindStringSubmatch(text)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1]), nil
	}

	return "", fmt.Errorf("not found")
}
