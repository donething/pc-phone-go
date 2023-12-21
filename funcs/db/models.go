package db

import "gorm.io/gorm"

// Subtitle 对应表 Subtitle 的一行行
type Subtitle struct {
	Code string `gorm:"unique"` // 番号
	Path string // 路径

	gorm.Model // 可以自动创建 ID、添加、删除、更新的时间
}

// TableName 返回表名。可通过实例 movie.TableName() 引用
func (Subtitle) TableName() string {
	return "m_subtitles"
}
