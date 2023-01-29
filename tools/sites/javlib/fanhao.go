package javlib

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// 已找到文件，抛出错误，以退出filepath.walk()
	hasFindError = errors.New("已找到文件：")
)

// 全局变量，保存上次剪贴板中的内容
// var clip = ""

// 监视剪贴板，查找本地文件
/*
func see() {
	// 读取剪贴板
	text, err := clipboard.ReadAll()
	// 此处可能误报，需要排除
	if err != nil && !strings.Contains(err.Error(), `The operation completed successfully`) {
		logger.Warn.Printf("读取剪贴板错误：%s（该错误可能为误报）\n", err)
	}

	// 解析番号
	fahhao := ResolveFanhao(text)

	// 如果剪贴板没有变动，或没有包含番号，则不需继续处理
	if text == clip || fahhao == "" {
		return
	}

	// 剪贴板变动过，需要更新剪贴板
	clip = text

	// 查找本地文件中是否已存在对应的番号文件
	logger.Info.Printf("发现剪贴板中的番号，开始查找对应的本地文件：")
	path, _ := seekFile(root, fahhao)
	dealFanhao(path, fahhao)
}
*/

// SeekFile 根据番号查找本地目录中，是否已存在番号对应的文件
// 找到，则返回文件路径；否则返回空字符串：""
func SeekFile(paths []string, keywords []string) (map[string]string, error) {
	var results = make(map[string]string)
	for _, p := range paths {
		err := filepath.Walk(p, func(path string, info os.FileInfo, errWalk error) error {
			if errWalk != nil {
				return errWalk
			}
			// 不需要检查目录和字幕文件夹
			if strings.Index(path, "字幕") >= 0 {
				// logger.Info.Printf("跳过字幕目录：%s\n", path)
				return filepath.SkipDir
			}
			if info.IsDir() {
				return nil
			}
			// 检查文件名
			for _, keyword := range keywords {
				if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(keyword)) {
					results[keyword] = strings.Replace(path, `\`, `/`, -1)
					// 找到文件，在退出filepath.walk()时，需要抛出错误，才能停止查找
					if len(results) == len(keywords) {
						return hasFindError
					}
				}
			}
			return nil
		})

		// 找到文件，返回目标路径
		if err == hasFindError {
			err = nil
			return results, nil
		}

		if err != nil {
			return results, err
		}
	}
	return results, nil
}
