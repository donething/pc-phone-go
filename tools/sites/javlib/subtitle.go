package javlib

import (
	"pc-phone-go/funcs/db"
	"strings"
)

// MatchSubtitle 匹配字幕，字幕路径的数组
//
// fanhao 番号
func MatchSubtitle(fanhao string) (string, error) {
	// 先转为大写的番号
	fh := strings.ToUpper(strings.TrimSpace(fanhao))

	var subtitle db.Subtitle
	if err := db.DB.Where("code LIKE ?", "%"+fh+"%").First(&subtitle).Error; err != nil {
		return "", err
	}

	return subtitle.Path, nil
}
