package douyin

// UserRoomStatus 返回给前端的用户信息、直播信息
type UserRoomStatus struct {
	*UserInfo
	*LiveInfo
}

// UserInfo 返回给前端的用户信息
type UserInfo struct {
	// 头像
	Avatar string `json:"avatar"`
	// 主播名/房间名
	Name string `json:"name"`
	// 是否在播
	Online bool `json:"online"`

	// 网页直播间号。**当主播不在播时，此值为空""**
	WebRid string `json:"webRid"`
	// 直播间的标题
	Title string `json:"title"`
}

// LiveInfo 返回给前端的直播信息
type LiveInfo struct {
	// 用户 ID。如 "97504125417"
	ID string `json:"id"`
	// sec id。如 "MS4wLjABAAAA_Wt4d-nxxxgEHzedsCrcK1f0kIXSlPZuI"
	SecUid string `json:"sec_uid"`
	// 头像
	Avatar string `json:"avatar"`
	// 主播名/房间名
	Name string `json:"name"`
	// 是否在播
	Online bool `json:"online"`
	// 直播间的标题
	Title string `json:"title"`
	// 流地址
	StreamUrl string `json:"streamUrl"`
}

// RoomStatus 从 Web 直播间获取直播流
//
// 可以获取到直播流地址(flv 或 .m3u8)
// 但当主播不开播时，难以获取直播间 web_rid，无法添加到chromium扩展中的关注
type RoomStatus struct {
	App struct {
		InitialState struct {
			RoomStore struct {
				RoomInfo struct {
					Anchor struct {
						IDStr       string `json:"id_str"`
						SecUID      string `json:"sec_uid"`
						Nickname    string `json:"nickname"`
						AvatarThumb struct {
							URLList []string `json:"url_list"`
						} `json:"avatar_thumb"`
						FollowInfo struct {
							FollowStatus    int    `json:"follow_status"`
							FollowStatusStr string `json:"follow_status_str"`
						} `json:"follow_info"`
					} `json:"anchor"`
					Room struct {
						IDStr string `json:"id_str"`
						// 状态为 2，表示在播
						Status int `json:"status"`
						// 状态的字符串形式，如"2"
						StatusStr string `json:"status_str"`
						// 房间标题
						Title string `json:"title"`
						// 观看人数
						UserCountStr string `json:"user_count_str"`
						// 直播间的封面图，不是用户头像。为数组，内容都一样，任选一个
						Cover struct {
							URLList []string `json:"url_list"`
						} `json:"cover"`
						// 直播流
						StreamURL struct {
							FlvPullURL struct {
								FULLHD1 string `json:"FULL_HD1"`
								HD1     string `json:"HD1"`
								SD1     string `json:"SD1"`
								SD2     string `json:"SD2"`
							} `json:"flv_pull_url"`
							// 如"HD1"
							DefaultResolution string `json:"default_resolution"`
							// 对应分辨率，ld为标清，sd为超清，hd为高清，FULL_HD1为全高清
							HlsPullURLMap struct {
								FULLHD1 string `json:"FULL_HD1"`
								HD1     string `json:"HD1"`
								SD1     string `json:"SD1"`
								SD2     string `json:"SD2"`
							} `json:"hls_pull_url_map"`
							// 默认的直播流，默认为 FULL_HD1
							HlsPullURL        string `json:"hls_pull_url"`
							StreamOrientation int    `json:"stream_orientation"`
						} `json:"stream_url"`
						AdminUserIdsStr []string `json:"admin_user_ids_str"`
						Owner           struct {
							IDStr       string `json:"id_str"`
							SecUID      string `json:"sec_uid"`
							Nickname    string `json:"nickname"`
							AvatarThumb struct {
								URLList []string `json:"url_list"`
							} `json:"avatar_thumb"`
						} `json:"owner"`
						Stats struct {
							// 总数，如"1万+"
							TotalUserStr string `json:"total_user_str"`
							// 当前，如"61"
							UserCountStr string `json:"user_count_str"`
						} `json:"stats"`
					} `json:"room"`
					// JSON中房间的ID，如"7166152177752158990"
					RoomID string `json:"roomId"`
					// Web 直播间的ID，如"492532016932"
					WebRid string `json:"web_rid"`
					// 分享的二维码
					QrcodeURL string `json:"qrcode_url"`
				} `json:"roomInfo"`
			} `json:"roomStore"`
		} `json:"initialState"`
	} `json:"app"`
}

// HomeInfo 从 Web 用户主页得到用户信息
//
// 注意：当用户没有在播时，`roomId`、`roomData`信息都为空（此时当然 `web_rid` 也为空）
// 注意：`Data`对应的JSON键名会改变，如"31"、"36"，需要程序自动识别，详见代码
type HomeInfo struct {
	HomeData
}

// HomeData 用户主页的JSON数据
type HomeData struct {
	// secUid，如"MS4wLjABAAAAK9qUm1QSQAl2XhQbnuATlqe2pyW-X3gW-KykZ_Gj93o"
	UID  string `json:"uid"`
	User struct {
		User struct {
			// 用户唯一UID，最重要的ID。如"338293379838045"
			UID string `json:"uid"`
			// secUid。如"MS4wLjABAAAAK9qUm1QSQAl2XhQbnuATlqe2pyW-X3gW-KykZ_Gj93o"
			SecUID string `json:"secUid"`
			// 暂时不可用。始终为"0"
			ShortID string `json:"shortId"`
			// 用户名
			Nickname string `json:"nickname"`
			// 用户签名（描述）
			Desc string `json:"desc"`
			// 小封面
			AvatarURL string `json:"avatarUrl"`
			// 大封面，推荐。如"//p3-pc.douyinpic.com/img/aweme-avatar/tos-cn-avt-00.jpeg"
			Avatar300URL string `json:"avatar300Url"`
			// 是否已关注该主播。1表示已关注，0表示未关注
			FollowStatus int `json:"followStatus"`
			// 该主播是否已关注我
			FollowerStatus int `json:"followerStatus"`
			// 发布的视频数量
			AwemeCount int `json:"awemeCount"`
			// 该主播的关注人数（粉丝数）
			FollowingCount int `json:"followingCount"`
			// 该主播关注的人数
			FollowerCount int `json:"followerCount"`
			// 该主播的关注人数（粉丝数），移动端。通常和 FollowingCount 一致
			MplatformFollowersCount int `json:"mplatformFollowersCount"`
			// 该主播喜欢的作品数
			FavoritingCount int `json:"favoritingCount"`
			// 该主播总获赞数
			TotalFavorited int `json:"totalFavorited"`
			// 该主播的抖音号
			UniqueID string `json:"uniqueId"`
			// 直播间信息
			RoomData struct {
				// 直播间状态。2表示开播
				Status int `json:"status"`
				Owner  struct {
					// Web端的房间号。如"492432036932"，打开前添加完整域名"https://live.douyin.com/492432036932"
					WebRid string `json:"web_rid"`
				} `json:"owner"`
			} `json:"roomData"`
			// 分享时的二维码
			ShareQrcodeURL string `json:"shareQrcodeUrl"`
			// 直播间ID，可能变动，用于API获取信息。如"7166142477752159000"
			RoomID             int64 `json:"roomId"`
			IsBlocked          bool  `json:"isBlocked"`
			IsBlock            bool  `json:"isBlock"`
			FavoritePermission int   `json:"favoritePermission"`
			// 是否向其他用户开放“喜欢”列表
			ShowFavoriteList bool `json:"showFavoriteList"`
			// IP地址。如"IP属地：浙江"
			IPLocation    string `json:"ipLocation"`
			IsGovMediaVip bool   `json:"isGovMediaVip"`
			// 是否为明星？？
			IsStar bool `json:"isStar"`
		} `json:"user"`
	} `json:"user"`
	Post struct {
		StatusCode int   `json:"statusCode"`
		HasMore    int   `json:"hasMore"`
		Cursor     int64 `json:"cursor"`
		MaxCursor  int64 `json:"maxCursor"`
		MinCursor  int64 `json:"minCursor"`
		Data       []struct {
			AwemeID    string `json:"awemeId"`
			AwemeType  int    `json:"awemeType"`
			GroupID    string `json:"groupId"`
			AuthorInfo struct {
				UID                    string `json:"uid"`
				SecUID                 string `json:"secUid"`
				Nickname               string `json:"nickname"`
				RemarkName             string `json:"remarkName"`
				AvatarURI              string `json:"avatarUri"`
				FollowerCount          int    `json:"followerCount"`
				TotalFavorited         int    `json:"totalFavorited"`
				FollowStatus           int    `json:"followStatus"`
				FollowerStatus         int    `json:"followerStatus"`
				EnterpriseVerifyReason string `json:"enterpriseVerifyReason"`
				CustomVerify           string `json:"customVerify"`
				RoomData               struct {
					Status int `json:"status"`
					Owner  struct {
						WebRid string `json:"web_rid"`
					} `json:"owner"`
					LiveTypeNormal bool `json:"live_type_normal"`
					PackMeta       struct {
						Scene string `json:"scene"`
						Env   string `json:"env"`
						Dc    string `json:"dc"`
					} `json:"pack_meta"`
				} `json:"roomData"`
				AvatarThumb struct {
					Height  int      `json:"height"`
					Width   int      `json:"width"`
					URI     string   `json:"uri"`
					URLList []string `json:"urlList"`
				} `json:"avatarThumb"`
			} `json:"authorInfo"`
			AwemeControl struct {
				CanComment     bool `json:"canComment"`
				CanForward     bool `json:"canForward"`
				CanShare       bool `json:"canShare"`
				CanShowComment bool `json:"canShowComment"`
			} `json:"awemeControl"`
			DanmakuControl struct {
			} `json:"danmakuControl"`
			Desc          string        `json:"desc"`
			AuthorUserID  int64         `json:"authorUserId"`
			CreateTime    int           `json:"createTime"`
			TextExtra     []interface{} `json:"textExtra"`
			UserDigged    bool          `json:"userDigged"`
			UserCollected bool          `json:"userCollected"`
			Video         struct {
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Ratio    string `json:"ratio"`
				Duration int    `json:"duration"`
				PlayAddr []struct {
					Src string `json:"src"`
				} `json:"playAddr"`
				PlayAPI      string `json:"playApi"`
				PlayAddrH265 []struct {
					Src string `json:"src"`
				} `json:"playAddrH265"`
				PlayAPIH265 string `json:"playApiH265"`
				BitRateList []struct {
					Width    int `json:"width"`
					Height   int `json:"height"`
					PlayAddr []struct {
						Src string `json:"src"`
					} `json:"playAddr"`
					PlayAPI     string `json:"playApi"`
					IsH265      int    `json:"isH265"`
					QualityType int    `json:"qualityType"`
					BitRate     int    `json:"bitRate"`
					VideoFormat string `json:"videoFormat"`
					GearName    string `json:"gearName"`
				} `json:"bitRateList"`
				Cover        string   `json:"cover"`
				CoverURLList []string `json:"coverUrlList"`
				CoverURI     string   `json:"coverUri"`
				DynamicCover string   `json:"dynamicCover"`
				OriginCover  string   `json:"originCover"`
				Meta         struct {
					BrightRatioMean       string `json:"bright_ratio_mean"`
					BrightnessMean        string `json:"brightness_mean"`
					DiffOverexposureRatio string `json:"diff_overexposure_ratio"`
					FullscreenMaxCrop     string `json:"fullscreen_max_crop"`
					Loudness              string `json:"loudness"`
					OverexposureRatioMean string `json:"overexposure_ratio_mean"`
					Peak                  string `json:"peak"`
					Qprf                  string `json:"qprf"`
					SrScore               string `json:"sr_score"`
					StdBrightness         string `json:"std_brightness"`
					TitleInfo             string `json:"title_info"`
				} `json:"meta"`
				BigThumbs  interface{} `json:"bigThumbs"`
				VideoModel interface{} `json:"videoModel"`
			} `json:"video"`
			MixInfo struct {
				Cover  string `json:"cover"`
				Status int    `json:"status"`
			} `json:"mixInfo"`
			IsPrivate     bool `json:"isPrivate"`
			IsFriendLimit bool `json:"isFriendLimit"`
			Download      struct {
				Prevent bool   `json:"prevent"`
				URL     string `json:"url"`
			} `json:"download"`
			ImpressionData     string `json:"impressionData"`
			FakeHorizontalInfo struct {
			} `json:"fakeHorizontalInfo"`
			Tag struct {
				IsTop          bool `json:"isTop"`
				RelationLabels bool `json:"relationLabels"`
				IsStory        bool `json:"isStory"`
				ReviewStatus   int  `json:"reviewStatus"`
				InReviewing    bool `json:"inReviewing"`
			} `json:"tag"`
			Stats struct {
				CommentCount int `json:"commentCount"`
				DiggCount    int `json:"diggCount"`
				ShareCount   int `json:"shareCount"`
				PlayCount    int `json:"playCount"`
				CollectCount int `json:"collectCount"`
			} `json:"stats"`
			ShareInfo struct {
				ShareURL      string `json:"shareUrl"`
				ShareLinkDesc string `json:"shareLinkDesc"`
			} `json:"shareInfo"`
			Status struct {
				AllowShare    bool `json:"allowShare"`
				IsReviewing   bool `json:"isReviewing"`
				IsDelete      bool `json:"isDelete"`
				IsProhibited  bool `json:"isProhibited"`
				PrivateStatus int  `json:"privateStatus"`
				ReviewStatus  int  `json:"reviewStatus"`
			} `json:"status"`
			WebRawData struct {
				OftenWatchInfo struct {
				} `json:"oftenWatchInfo"`
				VideoImageInfo struct {
				} `json:"videoImageInfo"`
				RecomPhrase string      `json:"recomPhrase"`
				PendantInfo interface{} `json:"pendantInfo"`
				CTR         struct {
					RecommendScore struct {
						RelatedScore       string `json:"relatedScore"`
						HotScore           string `json:"hotScore"`
						CollectionScore    string `json:"collectionScore"`
						RelatedSearchScore string `json:"relatedSearchScore"`
					} `json:"recommendScore"`
				} `json:"CTR"`
			} `json:"webRawData"`
			Music struct {
				ID         int64  `json:"id"`
				IDStr      string `json:"idStr"`
				Mid        string `json:"mid"`
				Author     string `json:"author"`
				Title      string `json:"title"`
				CoverThumb struct {
					URLList []string `json:"urlList"`
					URI     string   `json:"uri"`
					Width   int      `json:"width"`
					Height  int      `json:"height"`
				} `json:"coverThumb"`
				CoverMedium struct {
					URLList []string `json:"urlList"`
					URI     string   `json:"uri"`
					Width   int      `json:"width"`
					Height  int      `json:"height"`
				} `json:"coverMedium"`
				PlayURL struct {
					URLList []string `json:"urlList"`
					URI     string   `json:"uri"`
					Width   int      `json:"width"`
					Height  int      `json:"height"`
					URLKey  string   `json:"urlKey"`
				} `json:"playUrl"`
				SecUID    string `json:"secUid"`
				ShareInfo struct {
				} `json:"shareInfo"`
				Extra struct {
					HasEdited int `json:"hasEdited"`
				} `json:"extra"`
				Album       string `json:"album"`
				AvatarThumb struct {
					URLList []string `json:"urlList"`
					URI     string   `json:"uri"`
					Width   int      `json:"width"`
					Height  int      `json:"height"`
				} `json:"avatarThumb"`
				OwnerNickname     string `json:"ownerNickname"`
				CollectStat       int    `json:"collectStat"`
				BindedChallengeID int    `json:"bindedChallengeId"`
				Status            int    `json:"status"`
				CanNotPlay        bool   `json:"canNotPlay"`
				MusicName         string `json:"musicName"`
				IsOriginal        bool   `json:"isOriginal"`
				Duration          int    `json:"duration"`
				UserCount         int    `json:"userCount"`
			} `json:"music"`
			Images              []interface{} `json:"images"`
			ImageInfos          string        `json:"imageInfos"`
			ImageAlbumMusicInfo struct {
				BeginTime int `json:"beginTime"`
				EndTime   int `json:"endTime"`
				Volume    int `json:"volume"`
			} `json:"imageAlbumMusicInfo"`
			ImgBitrate   []interface{} `json:"imgBitrate"`
			SuggestWords []interface{} `json:"suggestWords"`
			SeoInfo      struct {
			} `json:"seoInfo"`
			RequestTime int64 `json:"requestTime"`
			LvideoBrief struct {
			} `json:"lvideoBrief"`
		} `json:"data"`
	} `json:"post"`
}
