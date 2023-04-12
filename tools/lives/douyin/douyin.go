// Package douyin 获取抖音直播间的状态
// [douyin_dynamic_push](https://github.com/nfe-w/douyin_dynamic_push/blob/master/query_douyin.py)
// [Python爬取抖音用户相关数据(目前最方便的方法）](http://www.dagoogle.cn/n/1307.html)
package douyin

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"net/url"
	"regexp"
)

var (
	client = dohttp.New(false, false)
	// 目前除了要浏览器代理，还需要提供 cookie，否则获取到的是滑动验证页面
	headers = map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
			"Chrome/109.0.0.0 Safari/537.36",
		"cookie": "douyin.com; __ac_nonce=063dcba970055a998205e; " +
			"__ac_signature=_02B4Z6wo00f01bSyleQAAIDA17hVjICxGym0kpFAAA7xf4; " +
			"ttwid=1|hbZh2jof8RzmIXsu_BGVWQS1bfxRIZoGEqPIYl4Ps8s|1675410071|" +
			"d846634924650788e7d99d099f25d5cf91baae821266cad889624e7ae604ad76; " +
			"home_can_add_dy_2_desktop=\"0\"; passport_csrf_token=2986562d012d3ba98667185e47bd660b; " +
			"passport_csrf_token_default=2986562d012d3ba98667185e47bd660b; " +
			"s_v_web_id=verify_ldo7wf2n_20LGcyTj_0FHk_4qQY_BOil_GeL7Km1fmacF; " +
			"msToken=swZMKXmEvtMrXyFFtgXAGXJE8p5FopEje7AcVn3qHojIATy-" +
			"GvxftK5PcbjeX05K6lyh55OomGxw0uH6hV_jkewjBvN-h-stkIEA5HwGJT1ZDgVkfYW_f5v_7Kb5gQ==; " +
			"ttcid=04b43e7318154b3fb129f9f1b9410ee065; strategyABtestKey=\"1675410110.969\"; " +
			"msToken=n7oN7QSD3cif-WFwCwktWyG3lbAPgiwLvRR5LyhTDs9WcvJE45ii2avHdsMEzIICwvD9Q2oa" +
			"nh2hntg4N3wCaCCVqQs3S63uxVlc3y7pZjrGThaQHZuEm_N_k32YWQ==; tt_scid=iTcB8Hh361Vcyaz3P" +
			"gCzRaYnW2-mJlGBieJqBNN7PkjSbazXj4TfcSQRMBAwRNsT31ce",
	}
)

// GetUserRoomStatus 获取抖音用户信息和直播流
//
// secUid 用户 Web 主页地址栏中最后一串字符。如"MS4wLjABAAAAK9qUxxxqe2pyW-X3gW-KykZ_Gj93o"
func GetUserRoomStatus(secUid string) (*UserRoomStatus, error) {
	var roomStatus UserRoomStatus

	// 获取用户信息
	userInfo, err := GetUserInfo(secUid)
	if err != nil {
		return nil, err
	}
	roomStatus.UserInfo = userInfo

	// 不在播，无法获取直播流
	if !userInfo.Online {
		return &roomStatus, nil
	}

	// 获取直播流
	liveInfo, err := GetLiveInfo(userInfo.WebRid)
	if err != nil {
		return nil, err
	}
	roomStatus.LiveInfo = liveInfo

	return &roomStatus, nil
}

// GetUserInfo 获取抖音用户信息
//
// secUid 用户 Web 主页地址栏中最后一串字符。如"MS4wLjABAAAAK9qUxxxqe2pyW-X3gW-KykZ_Gj93o"
func GetUserInfo(secUid string) (*UserInfo, error) {
	// 提取主页的用户信息
	u := fmt.Sprintf("https://www.douyin.com/user/%s", secUid)
	homeInfo, err := parseRenderData[map[string]interface{}](u)
	if err != nil {
		return nil, fmt.Errorf("获取主页信息出错(%s)：%w", secUid, err)
	}

	// 由于键名经常变化，需要程序自动识别键名
	var homeData HomeData
	for key := range *homeInfo {
		// 识别过程：遍历对象，先将属性转为键值对，如果该属性下存在 uid 属性，则说明是目标
		// 注意 uid 需要判断 nil，不能直接判断 != ""，因为此时 uid 为 interface{}
		// 锁定目标后，还要经过序列表、反序列化，转为真正的类型 HomeData
		if userObj, ok := (*homeInfo)[key].(map[string]interface{}); ok &&
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
		return nil, fmt.Errorf("没有从数据中得到用户信息，需要适配")
	}

	user := homeData.User.User

	// 头像。URL 缺少"https"，需要补上
	avatar := "https:" + user.AvatarURL
	if user.Avatar300URL != "" {
		avatar = "https:" + user.Avatar300URL
	}

	var status = UserInfo{
		Avatar: avatar,
		WebRid: user.RoomData.Owner.WebRid,
		Name:   user.Nickname,

		Online: user.RoomData.Status == 2,
		Title:  user.Desc,
	}

	return &status, nil
}

// GetLiveInfo 获取抖音直播流地址
//
// roomid 直播间号
func GetLiveInfo(roomid string) (*LiveInfo, error) {
	// 提取直播间的直播信息
	u := fmt.Sprintf("https://live.douyin.com/%s", roomid)
	roomStatus, err := parseRenderData[RoomStatus](u)
	if err != nil {
		return nil, fmt.Errorf("获取直播间信息出错(%s)：%w", roomid, err)
	}

	// 是否开播，关系到页面中是否存在数据
	if roomStatus.App.InitialState.RoomStore.RoomInfo.Anchor.Nickname == "" {
		return nil, fmt.Errorf("不存在直播间(%s)", roomid)
	}

	roomInfo := roomStatus.App.InitialState.RoomStore.RoomInfo
	liveInfo := LiveInfo{
		ID:        roomInfo.Anchor.IDStr,
		SecUid:    roomInfo.Anchor.SecUID,
		Avatar:    roomInfo.Anchor.AvatarThumb.URLList[0],
		Name:      roomInfo.Anchor.Nickname,
		Online:    roomStatus.App.InitialState.RoomStore.RoomInfo.Room.Status == 2,
		Title:     roomInfo.Room.Title,
		StreamUrl: roomInfo.Room.StreamURL.FlvPullURL.FULLHD1,
	}

	return &liveInfo, nil
}

// 提取网页中的 RENDER_DATA
func parseRenderData[T any](dyUrl string) (*T, error) {
	// 获取抖音网页文本
	text, err := client.GetText(dyUrl, headers)
	if err != nil {
		return nil, fmt.Errorf("获取网页内容出错(%s)：%w", dyUrl, err)
	}

	// 页面会携带一段ID为"RENDER_DATA"的脚本，里面带有用户数据信息
	// 可以在页面控制台中执行`copy(decodeURIComponent(document.querySelector("#RENDER_DATA").text))`获取
	reg := regexp.MustCompile(`(?m)id="RENDER_DATA".+?>(.+?)<`)
	matches := reg.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil, fmt.Errorf("没有匹配到'RENDER_DATA'数据")
	}

	// 反转义非法字符
	dataText, err := url.QueryUnescape(matches[1])
	if err != nil {
		return nil, fmt.Errorf("反转义非法字符出错：%w", err)
	}

	// 解析数据
	var obj = new(T)
	err = json.Unmarshal([]byte(dataText), obj)

	if err != nil {
		return nil, fmt.Errorf("解析数据出错：%w", err)
	}

	return obj, nil
}
