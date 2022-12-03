package tolocal

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"os"
	"path/filepath"
	. "pc-phone-go/conf"
	"pc-phone-go/tools/pics/pcomm"
	"strings"
)

// Save 保存到本地
func Save(album pcomm.Album) error {
	// 保存到本地文件时，先创建目录
	destDir := filepath.Join(Conf.Pics.LocalRoot, album.Plat, album.UID, album.ID)
	err := os.MkdirAll(destDir, 0644)
	if err != nil {
		return fmt.Errorf("创建目录'%s'出错：%s\n", destDir, err)
	}

	for _, url := range album.URLs {
		// 下载链接，获取二进制数组数据
		bs, err := pcomm.Client.Get(url, album.Header)
		if err != nil {
			return fmt.Errorf("下载文件'%s'出错：%s\n", url, err)
		}

		// 写入文件
		filename := filepath.Base(url)
		if i := strings.Index(filename, "?"); i >= 0 {
			filename = filename[:i]
		}
		_, err = dofile.Write(bs, filepath.Join(destDir, filename), dofile.OTrunc, 0644)
		if err != nil {
			return fmt.Errorf("将数据保存到文件'%s'时出错：%s\n", filename, err)
		}
	}
	return nil
}
