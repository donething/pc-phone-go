package database

import (
	"encoding/json"
	"github.com/donething/utils-go/dodb/dobadger"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/tools/pics/pcomm"
)

const (
	dbDir = "mydb"
)

var (
	DB *dobadger.DoBadger
)

func init() {
	// 打开数据库
	DB = dobadger.Open(dbDir, nil)
}

// GetAll 读取保存的指定桶内的图集信息
func GetAll(prefixStr string) ([]pcomm.Album, error) {
	data, err := DB.QueryPrefix(prefixStr, "")
	if err != nil {
		return nil, err
	}

	albums := make([]pcomm.Album, 0, 0)
	for _, bs := range data {
		var album pcomm.Album
		errUnmarshal := json.Unmarshal(bs, &album)
		if errUnmarshal != nil {
			logger.Error.Printf("解析 JSON 数据出错：'%s' ==> '%s'\n", string(bs), err)
			continue
		}
		album.RetryFrom = prefixStr
		albums = append(albums, album)
	}

	return albums, nil
}
