// Package weibo 下载、发送微博指定用户的所有图集到 PicTG
package weibo

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dofile"
	"math/rand"
	"pc-phone-conn-go/funcs/logger"
	"pc-phone-conn-go/tools/pics/pcomm"
	"time"
)

// PostPage 一次请求某页 API 时的返回信息
type PostPage struct {
	Data struct {
		List []struct {
			Idstr   string   `json:"idstr"`
			PicIds  []string `json:"pic_ids"`
			TextRaw string   `json:"text_raw"`
			Created int64    `json:"created_at"`
		} `json:"list"`
	} `json:"data"`
}

const (
	// API
	mymblogAPI  = "https://weibo.com/ajax/statuses/mymblog?uid=%s&page=%d&feature=1"
	downloadAPI = "https://weibo.com/ajax/common/download?pid=%s"
)

// 从文件解析图集信息
func unMarshalPosts(path string) map[string]pcomm.Album {
	bs, err := dofile.Read(path)
	if err != nil {
		logger.Error.Println("读取图集数据的文件时出错：", err)
		return nil
	}
	var posts map[string]pcomm.Album
	err = json.Unmarshal(bs, &posts)
	if err != nil {
		logger.Error.Println("解析图集数据时出错：", err)
		return nil
	}

	return posts
}

// 保存图集信息到文件
func marshalPosts(posts map[string]pcomm.Album, path string) {
	bs, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		logger.Error.Println("文本化图集数据时出错：", err)
		return
	}
	_, err = dofile.Write(bs, path, dofile.OTrunc, 0644)
	if err != nil {
		logger.Error.Println("将图集数据保存到文件时出错：", err)
		return
	}
}

// 获取微博指定用户所有帖子的图集
func getAllAlbums(uid string, idDone string, headers *map[string]string) map[string]pcomm.Album {
	// 用于保存所有图集，将返回
	posts := make(map[string]pcomm.Album)

	// 读取所有帖子
	page := 1
	var postPage PostPage
	for {
		// 读取 API，解析
		bs, err := pcomm.Client.Get(fmt.Sprintf(mymblogAPI, uid, page), *headers)
		if err != nil {
			logger.Error.Println("联网获取图集数据时出错：", err)
			return nil
		}
		err = json.Unmarshal(bs, &postPage)
		if err != nil {
			logger.Error.Println("解析获取到底图集数据时出错：", err)
			return nil
		}

		// 返回内容的帖子数量为 0 时，表示遍历完成，退出循环
		if len(postPage.Data.List) == 0 {
			return posts
		}

		// 遍历帖子
		for _, post := range postPage.Data.List {
			// 当读取的帖子的 idstr 和已保存的进度记录相同时，说明已完成任务，直接返回数据
			if post.Idstr == idDone {
				return posts
			}

			// 读取、保存该贴的图集
			task := pcomm.Album{
				Plat:    pcomm.TagWeibo,
				Caption: post.TextRaw,
				ID:      post.Idstr,
				UID:     uid,
				Created: post.Created,
				URLs:    make([]string, len(post.PicIds)),
				URLsM:   nil,
				Header:  *headers,
			}
			for _, pid := range post.PicIds {
				task.URLs = append(task.URLs, fmt.Sprintf(downloadAPI, pid))
			}

			// 添加到所有图集中，以返回
			posts[post.Idstr] = task
		}
		logger.Info.Printf("[%s][%s] 已添加第 %d 页的图集\n", page)
		page++

		// 等待不固定的时间，以防被禁止访问
		r := rand.Intn(5)
		time.Sleep(time.Duration(r) * time.Second)
	}
}
