package lives

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pc-phone-go/entity"
	"pc-phone-go/tools/lives/douyin"
)

// GetDouyinRoom 获取抖音直播间状态
//
// Get /api/lives/douyin/live?sec_uid=test-uid
func GetDouyinRoom(c *gin.Context) {
	secUid := c.Query("sec_uid")
	if secUid == "" {
		c.JSON(http.StatusOK, entity.Rest{Code: 1000, Msg: "没有提取到请求参数'web_rid'"})
		return
	}

	status, err := douyin.GetDouyinRoomStatus(secUid)
	if err != nil {
		c.JSON(http.StatusOK, entity.Rest{Code: 2000, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity.Rest{Code: 0, Msg: "抖音直播间状态", Data: status})
}
