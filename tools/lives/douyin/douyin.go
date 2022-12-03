// Package douyin 获取抖音直播间的状态
// [douyin_dynamic_push](https://github.com/nfe-w/douyin_dynamic_push/blob/master/query_douyin.py)
// [Python爬取抖音用户相关数据(目前最方便的方法）](http://www.dagoogle.cn/n/1307.html)
package douyin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"pc-phone-go/tools/pics/pcomm"
	"regexp"
)

// GetDouyinRoomStatus 获取抖音直播间的状态
//
// secUid 用户Web主页地址栏中最后一串字符。如"MS4wLjABAAAAK9qUm1QSQAl2XhQbnuATlqe2pyW-X3gW-KykZ_Gj93o"
func GetDouyinRoomStatus(secUid string) (*RoomStatusTiny, error) {
	// 提取直播间的状态文本
	u := fmt.Sprintf("https://www.douyin.com/user/%s", secUid)
	text, err := pcomm.Client.GetText(u, nil)
	if err != nil {
		return nil, fmt.Errorf("获取抖音网页内容出错。%w", err)
	}

	// 页面会携带一段ID为"RENDER_DATA"的脚本，里面带有用户数据信息
	// 可以在页面控制台中执行`copy(decodeURIComponent(document.querySelector("#RENDER_DATA").text))`获取
	reg := regexp.MustCompile(`id="RENDER_DATA.+?>(.+?)<`)
	matches := reg.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil, fmt.Errorf("没有匹配到该主播的数据")
	}

	// 反转义非法字符
	dataText, err := url.QueryUnescape(matches[1])
	if err != nil {
		return nil, fmt.Errorf("反转义非法字符出错。%w", err)
	}

	// 解析数据
	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(dataText), &obj)
	if err != nil {
		return nil, fmt.Errorf("解析数据出错。%w", err)
	}

	// 由于键名经常变化，需要程序自动识别键名
	var homeData HomeData
	for key := range obj {
		// 识别过程：遍历对象，先将属性转为键值对，如果该属性下存在 uid 属性，则说明是目标
		// 注意 uid 需要判断 nil，不能直接判断 != ""，因为此时 uid 为 interface{}
		// 锁定目标后，还要经过序列表、反序列化，转为真正的类型 HomeData
		if userObj, ok := obj[key].(map[string]interface{}); ok &&
			userObj["uid"] != nil && userObj["uid"] != "" {
			// 转为真正的类型 HomeData
			tmp, err := json.Marshal(userObj)
			if err != nil {
				return nil, fmt.Errorf("已锁定目标数据，但序列化时出错。%w", err)
			}

			err = json.Unmarshal(tmp, &homeData)
			if err != nil {
				return nil, fmt.Errorf("已锁定目标数据，但反序列化时出错。%w", err)
			}
			break
		}
	}

	// 验证是否已正确处理好数据
	if homeData.UID == "" {
		return nil, fmt.Errorf("没有从数据从得到用户信息，需要适配")
	}

	user := homeData.User.User
	// 是否在播
	online := 0
	if user.RoomData.Status == 2 {
		online = 1
	}
	// 头像
	avatar := user.AvatarURL
	if user.Avatar300URL != "" {
		avatar = user.Avatar300URL
	}
	// URL 缺少"https"，需要补上
	avatar = "https:" + avatar

	// 网页直播间地址
	liveUrl := fmt.Sprintf("https://live.douyin.com/%s", user.RoomData.Owner.WebRid)

	var status = RoomStatusTiny{
		Avatar:    avatar,
		LiveUrl:   liveUrl,
		StreamUrl: "暂时不需要实现",
		Name:      user.Nickname,
		Online:    online,
		Title:     user.Desc,
	}

	return &status, nil
}
