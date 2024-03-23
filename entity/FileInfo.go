package entity

// FileInfo 文件信息
type FileInfo struct {
	Path    string `json:"path"` // 文件的完整路径
	Name    string `json:"name"` // 文件名
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"modTime"` // 时间戳（毫秒）
}
