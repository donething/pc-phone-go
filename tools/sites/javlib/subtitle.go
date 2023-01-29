package javlib

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"pc-phone-go/funcs/database"
	"strings"
)

// MatchSubtitle 匹配字幕，字幕路径的数组
//
// fanhao 番号
func MatchSubtitle(fanhao string) ([]string, error) {
	// 先转为大写的番号
	fh := strings.ToUpper(fanhao)
	// 通过键准确查找
	bs, err := database.DB.Get([]byte(database.PreKeySub + fh))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}

		return nil, err
	}
	// 已找到，返回
	path := string(bs)
	if path != "" {
		return []string{path}, nil
	}

	// 未找到时，遍历所有键模糊查找
	data, err := database.DB.QueryPrefix(database.PreKeySub, fanhao)
	if err != nil {
		return nil, err
	}

	var payload = make([]string, 0, len(data))
	for _, bsPath := range data {
		payload = append(payload, string(bsPath))
	}

	return payload, nil
}
