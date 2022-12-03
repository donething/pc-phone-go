// Package ql 控制青龙面板
package ql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pc-phone-go/conf"
	. "pc-phone-go/conf"
	"pc-phone-go/entity"
	"pc-phone-go/funcs/logger"
	"time"
)

func init() {
	// 指定默认配置
	if Conf.QLPanel.Port == 0 {
		Conf.QLPanel.Port = 5700
	}
}

// 接口的域名
var host = fmt.Sprintf("http://127.0.0.1:%d", conf.Conf.QLPanel.Port)

var client = dohttp.New(10*time.Second, false, false)

// SetEnv 设置环境变量
//
// 此函数**不**能用于有多个同名环境变量的场景
func SetEnv(c *gin.Context) {
	// 获取 Token
	headers, err := getTokenHeaders()
	if err != nil {
		logger.Error.Printf("获取 青龙 Token Headers 出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "获取 青龙 Token Headers 出错",
		})
		return
	}

	var setEnvReq SetEnvReq
	err = c.BindJSON(&setEnvReq)
	if err != nil {
		logger.Error.Printf("解析设置环境变量的请求内容时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "解析设置环境变量的请求内容时出错",
		})
		return
	}

	// 先判断已有的环境变量中是否存在同名变量，有就需要传递 ID 以更新，而不是添加
	var data interface{} = []SetEnvReq{setEnvReq}
	// 添加环境变量用"POST，更新用"PUT"
	var method = "POST"
	envs, err := getEnvs(setEnvReq.Name, headers)
	if len(envs) == 1 {
		method = "PUT"
		data = UpEnvReq{
			ID:        envs[0].ID,
			SetEnvReq: &setEnvReq,
		}
	}

	// 设置、更新 环境变量
	// 发送请求
	putData, err := json.Marshal(data)
	if err != nil {
		logger.Error.Printf("序列号设置环境变量的内容时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "序列号设置环境变量的内容时出错",
		})
		return
	}

	logger.Info.Printf("设置环境变量的请求：method: '%s'，body: '%s'\n", method, string(putData))
	req, err := http.NewRequest(method, fmt.Sprintf("%s/open/envs", host), bytes.NewReader(putData))
	if err != nil {
		logger.Error.Printf("创建设置环境变量的请求时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "创建设置环境变量的请求时出错",
		})
		return
	}

	// 注意请求头中的表单类型
	headers["Content-Type"] = "application/json"
	resp, err := client.Exec(req, headers)
	if err != nil {
		logger.Error.Printf("发送设置环境变量的请求时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "发送设置环境变量的请求时出错",
		})
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("读取设置环境变量的响应时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "读取设置环境变量的响应时出错",
		})
		return
	}

	// 解析响应
	var basicResp BasicResp
	err = json.Unmarshal(bs, &basicResp)
	if err != nil {
		logger.Error.Printf("解析设置环境变量的响应时出错：%s ==> '%s'\n", err, string(bs))
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6000,
			Msg:  "解析设置环境变量的响应时出错",
		})
		return
	}

	// 没有获取到正确内容
	if basicResp.Code != 200 {
		logger.Error.Printf("设置环境变量时出错：%s\n", string(bs))
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6100,
			Msg:  "设置环境变量时出错",
		})
		return
	}

	logger.Info.Printf("已设置环境变量：%s\n", setEnvReq.Name)
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  "已设置环境变量",
	})
}

// 获取指定可以的环境变量，结果可能有多个
// key 为空时，获取所有环境变量
func getEnvs(key string, headers map[string]string) ([]Env, error) {
	bs, err := client.Get(fmt.Sprintf("%s/open/envs", host), headers)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var envsResp GetEnvsResp
	err = json.Unmarshal(bs, &envsResp)
	if err != nil {
		return nil, err
	}

	// 没有获取到正确内容
	if envsResp.Code != 200 {
		return nil, fmt.Errorf("获取环境变量出错：%s", string(bs))
	}

	// 返回符合要求的环境变量
	var payload = make([]Env, 0, len(envsResp.Data))
	for _, env := range envsResp.Data {
		// 当 key 不为空时，需要需要获取指定 key 的环境变量，其它的忽略
		if key != "" && env.Name != key {
			continue
		}
		payload = append(payload, env)
	}

	return payload, nil
}

// StartCommCrons 执行常规定时任务
func StartCommCrons(c *gin.Context) {
	num, err := StartCommCronsCall()
	if err != nil {
		logger.Error.Printf("执行定时任务时出错：%s\n", err)
		c.JSON(http.StatusOK, entity.Rest{
			Code: 6100,
			Msg:  fmt.Sprintf("执行定时任务时出错：%s", err),
		})
		return
	}

	logger.Info.Printf("已发送执行定时任务的请求，共计 %d 个任务\n", num)
	c.JSON(http.StatusOK, entity.Rest{
		Code: 0,
		Msg:  fmt.Sprintf("已发送执行定时任务的请求，共计 %d 个任务\n", num),
	})
}

// StartCommCronsCall 执行定时任务
//
// 返回已提交执行的任务数量
func StartCommCronsCall() (int, error) {
	// 获取 Token
	headers, err := getTokenHeaders()
	if err != nil {
		return 0, fmt.Errorf("获取青龙 Token Headers 出错：%w", err)
	}

	// 获取任务
	crons, err := getAllCrons(headers)
	if err != nil {
		return 0, fmt.Errorf("获取所有定时任务出错：%w", err)
	}

	// 过滤任务，排除置顶、禁用的任务
	var ids = make([]int, 0, len(crons))
	for _, cron := range crons {
		// 排除置顶、禁用的任务
		if cron.IsDisabled == 1 || cron.IsPinned == 1 {
			continue
		}
		ids = append(ids, cron.ID)
	}

	// 发送请求，执行定时任务
	putData, err := json.Marshal(ids)
	if err != nil {
		return 0, fmt.Errorf("执行定时任务时出错：%w", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/open/crons/run", host),
		bytes.NewReader(putData))
	if err != nil {
		return 0, fmt.Errorf("创建网络请求出错：%w", err)
	}

	// 注意请求头中的表单类型
	headers["Content-Type"] = "application/json"
	resp, err := client.Exec(req, headers)
	if err != nil {
		return 0, fmt.Errorf("网络错误：%w", err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应内容出错：%w", err)
	}

	// 解析响应
	var basicResp BasicResp
	err = json.Unmarshal(bs, &basicResp)
	if err != nil {
		return 0, fmt.Errorf("解析JSON文本出错：%w ==> '%s'", err, string(bs))
	}

	// 没有获取到正确内容
	if basicResp.Code != 200 {
		return 0, fmt.Errorf("状态响应码异常：'%s'", string(bs))
	}

	return len(ids), nil
}

// 获取所有定时任务
func getAllCrons(headers map[string]string) ([]Cron, error) {
	bs, err := client.Get(fmt.Sprintf("%s/open/crons", host), headers)
	if err != nil {
		return nil, fmt.Errorf("网络错误：%w", err)
	}

	// 解析响应
	var cronsResp GetCronsResp
	err = json.Unmarshal(bs, &cronsResp)
	if err != nil {
		return nil, fmt.Errorf("解析JSON文本出错：%w ==> '%s'", err, string(bs))
	}

	// 没有获取到正确内容
	if cronsResp.Code != 200 {
		return nil, fmt.Errorf("状态响应码异常：'%s'", string(bs))
	}

	// 返回符合要求的定时任务
	var payload = make([]Cron, 0, cronsResp.Data.Total)
	for _, cron := range cronsResp.Data.Data {
		payload = append(payload, cron)
	}

	return payload, nil
}

// 获取 token
func getTokenHeaders() (map[string]string, error) {
	// 发送请求
	bs, err := client.Get(fmt.Sprintf("%s/open/auth/token?client_id=%s&client_secret=%s",
		host, conf.Conf.QLPanel.ClientID, conf.Conf.QLPanel.ClientSecret), nil)
	if err != nil {
		return nil, fmt.Errorf("网络错误：%w", err)
	}

	// 解析响应
	var tokenResp TokenResp
	err = json.Unmarshal(bs, &tokenResp)
	if err != nil {
		return nil, fmt.Errorf("解析JSON文本出错：%w ==> '%s'", err, string(bs))
	}

	// 没有获取到正确内容
	if tokenResp.Code != 200 {
		return nil, fmt.Errorf("状态响应码异常：'%s'", string(bs))
	}

	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", tokenResp.Data.Token),
		"Accept": "application/json"}
	return headers, nil
}
