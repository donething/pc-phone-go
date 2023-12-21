package db

import (
	"github.com/donething/utils-go/dolog"
	"github.com/donething/utils-go/dotext"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"path/filepath"
	"pc-phone-go/conf"
	"pc-phone-go/funcs/logger"
	"strings"
	"time"
)

const (
	// 数据库的路径
	dbPath = "pc-phone-go.db"

	// 字幕文件的格式
	subPatten = ".srt|.ass|.ssa"
)

var gormConf = gorm.Config{
	NowFunc: func() time.Time {
		// 设置 NowFunc 为返回 UTC 时间
		return time.Now().UTC()
	},
}

var (
	DB *gorm.DB
)

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gormConf)
	dolog.CkPanic(err)

	err = DB.AutoMigrate(&Subtitle{})
	dolog.CkPanic(err)
}

// 创建字幕文件的索引到数据库
//
// 键为大写的番号（无法解析出番号时为字幕文件名），值为字幕的完整路径（因为字幕文件是按文件夹分类的，需要保存完整路径）
func indexSubtitles() {
	var payload = make([]Subtitle, 0, 1000)
	errWalk := filepath.WalkDir(conf.Conf.Javlib.SubDir, func(path string, d fs.DirEntry, err error) error {
		// 跳过目录、非字幕文件
		if d.IsDir() || !strings.Contains(subPatten, filepath.Ext(d.Name())) {
			// log.Printf("跳过目录、非字幕文件：'%s'\n", path)
			return nil
		}

		// 将番号、其字幕路径保存到数据库
		code := dotext.ResolveFanhao(d.Name())
		if code == "" {
			code = d.Name()
		}

		code = strings.ToUpper(code)
		payload = append(payload, Subtitle{Code: code, Path: path})
		// log.Printf("已记录番号'%s': '%s'\n", code, path)
		return nil
	})
	if errWalk != nil {
		logger.Error.Printf("遍历路径下的字幕文件时出错：%s\n", errWalk)
	}

	// 批量保存到桶
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&payload, len(payload)).Error; err != nil {
			return err // Rollback
		}

		return nil
	})

	if err != nil {
		logger.Error.Printf("批量写入字幕数据到数据库时出错：%s\n", err)
		return
	}

	logger.Warn.Printf("已创建字幕索引数据库\n")
}
