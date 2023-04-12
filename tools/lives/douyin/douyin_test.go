package douyin

import (
	"testing"
)

func TestGetDouyinRoomStatus(t *testing.T) {
	userInfo, err := GetUserInfo("MS4wLjABAAAA_Wt4d-n6OchUYOum24PdgEHzedsCrcK1f0kIXSlPZuI")
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("%v", userInfo)
}

func TestGetUserRoomStatus(t *testing.T) {
	status, err := GetUserRoomStatus("MS4wLjABAAAAK9qUm1QSQAl2XhQbnuATlqe2pyW-X3gW-KykZ_Gj93o")
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("%+v, %+v\n", status.LiveInfo, status.UserInfo)
}

func TestGetLiveStream(t *testing.T) {
	liveInfo, err := GetLiveInfo("897151189145")
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("%+v\n", liveInfo)
}
