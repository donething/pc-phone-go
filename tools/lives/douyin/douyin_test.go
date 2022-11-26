package douyin

import (
	"testing"
)

func TestGetDouyinRoomStatus(t *testing.T) {
	status, err := GetDouyinRoomStatus("MS4wLjABAAAAK9qUm1QSQAl2XhQbnuATlqe2pyW-X3gW-KykZ_Gj93o")
	if err != nil {
		t.Log(err)
		return
	}

	t.Logf("%v", status)
}
