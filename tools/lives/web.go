package lives

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"pc-phone-go/tools/lives/douyin"
)

// GetDouyinUserInfo 获取抖音直播间状态
//
// Get /api/lives/douyin/live?sec_uid=test-uid
func GetDouyinUserInfo(c *gin.Context) {
	secUid := c.Query("sec_uid")
	if secUid == "" {
		c.JSON(http.StatusOK, entity.Rest{Code: 1000, Msg: "缺少参数'sec_uid'"})
		return
	}

	userInfo, err := douyin.GetUserInfo(secUid)
	if err != nil {
		logger.Error.Printf("获取抖音用户信息出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{Code: 1010, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "抖音用户信息", Data: userInfo})
}
