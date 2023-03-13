package javlib

import (
	"pc-phone-go/funcs/db"
	"strings"
)

// MatchSubtitle 匹配字幕，字幕路径的数组
//
// fanhao 番号
func MatchSubtitle(fanhao string) ([]string, error) {
	// 先转为大写的番号
	fh := strings.ToUpper(strings.TrimSpace(fanhao))

	// 通过键准确查找
	bs, err := db.DB.Get([]byte(fh), db.BkSubtitle)
	if err != nil {
		return nil, err
	}

	// 已找到时，直接返回
	if bs != nil {
		return []string{string(bs)}, nil
	}

	// 未找到时，模糊查找
	data, err := db.DB.Query(fh, db.BkSubtitle)
	if err != nil {
		return nil, err
	}

	var payload = make([]string, 0, len(data))
	for _, bsPath := range data {
		payload = append(payload, string(bsPath))
	}

	return payload, nil
}
