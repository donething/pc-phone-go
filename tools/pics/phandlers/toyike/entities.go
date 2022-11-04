package toyike

// PreResp 响应
type PreResp struct {
	ReturnType int    `json:"return_type"`
	Uploadid   string `json:"uploadid"`
	Errno      int    `json:"errno"`
}

// UpResp 上传分段的响应
type UpResp struct {
	Md5      string `json:"md5"`
	Partseq  string `json:"partseq"`
	Uploadid string `json:"uploadid"`
}

// CreateResp 创建文件的响应
type CreateResp struct {
	Errno int `json:"errno"`
	Data  struct {
		Errno          int    `json:"errno"`
		Category       int    `json:"category"`
		FromType       int    `json:"from_type"`
		FSID           int64  `json:"fs_id"`
		Isdir          int    `json:"isdir"`
		Md5            string `json:"md5"`
		Ctime          int64  `json:"ctime"`
		Mtime          int64  `json:"mtime"`
		ShootTime      int64  `json:"shoot_time"`
		Path           string `json:"path"`
		ServerFilename string `json:"server_filename"`
		Size           int64  `json:"size"`
		ServerMd5      string `json:"server_md5"`
	} `json:"data"`
}

// FilesResp 文件列表
type FilesResp struct {
	Cursor string `json:"cursor"`
	Errno  int    `json:"errno"`
	List   []struct {
		Fsid int64 `json:"fsid"`
	} `json:"list"`
}
