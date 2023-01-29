package database

import (
	"encoding/json"
	"github.com/donething/utils-go/dodb/dobadger"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dotext"
	"io/fs"
	"path/filepath"
	"pc-phone-go/conf"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/tools/pics/pcomm"
	"strings"
)

const (
	// 数据库的路径
	dbDir = "mydb"

	// PreKeySub 数据库中存放字幕信息的键前缀
	PreKeySub = "sub_"

	// 字幕文件的格式
	subPatten = ".srt|.ass|.ssa"
)

var (
	DB *dobadger.DoBadger
)

func init() {
	// 先判断数据库是否存在时，以便执行数据索引等操作
	// 因为打开就会创建数据库的文件，所以需要放在 Exists() 之后，在判断其值之前
	exist, err := dofile.Exists(dbDir)
	if err != nil {
		logger.Error.Printf("判断数据库是否存在时出错：%s\n", err)
		return
	}

	// 打开数据库
	DB = dobadger.Open(dbDir, nil)

	if !exist {
		indexSubtitles()
	}
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

// 创建字幕文件的索引到数据库
//
// 键为大写的番号（无法解析出番号时为字幕文件名），值为字幕的完整路径（因为字幕文件是按文件夹分类的，需要保存完整路径）
func indexSubtitles() {
	var payload = make(map[string][]byte)
	errWalk := filepath.WalkDir(conf.Conf.Javlib.SubDir, func(path string, d fs.DirEntry, err error) error {
		// 跳过目录、非字幕文件
		if d.IsDir() || !strings.Contains(subPatten, filepath.Ext(d.Name())) {
			// log.Printf("跳过目录、非字幕文件：'%s'\n", path)
			return nil
		}

		// 将番号、其字幕路径保存到数据库
		key := dotext.ResolveFanhao(d.Name())
		if key == "" {
			key = d.Name()
		}

		key = strings.ToUpper(key)
		payload[PreKeySub+key] = []byte(path)
		// log.Printf("已记录番号'%s': '%s'\n", key, path)
		return nil
	})
	if errWalk != nil {
		logger.Error.Printf("遍历路径下的字幕文件时出错：%s\n", errWalk)
	}

	// 批量保存到桶
	err := DB.BatchSet(payload)
	if err != nil {
		logger.Error.Printf("批量写入字幕数据到数据库时出错：%s\n", err)
		return
	}

	logger.Info.Printf("已创建字幕索引数据库\n")
}
