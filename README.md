# pc-phone-go

# 可在手机端向 PC 端发送文件、文本，互设剪贴板

# 此为 PC 服务端

## 编译

```shell
# 生成`GUI`程序，点击运行不打开终端
go build -ldflags="-H windowsgui"
```

# 使用说明

本服务默认监听地址为：http://IP:8800/clip (需将 IP 改为 PC 端内网 IP)
通过 POST 方式向此地址传递数据，表单键值有两个：

1. type 字符串 操作/数据的类型，可选的值为
    1. "getclip"：获取 PC 端剪贴板（此时 content 参数可为空）
    2. "URL"：向 PC 端发送链接，将自动打开此链接
    3. "文本"：向 PC 端发送文本，当字节数小于 512 B 时，复制到剪贴板，否则作为文件保存到用户下载目录
    4. 当没有匹配到上述情况时，将作为文件保存到 用户下载目录
2. content 文件 需保存到 PC 端的文件（保存在用户下载文满中）

## javlib

* `/api/openfile` 打开/显示本地的文件（夹）。`POST`传递 JSON 数据：`{method: string, path: string}`。`method`可选值为`open`
  、`show`，`path`
  为文件（夹）的路径
* `/api/fanhaos` 查询本地是否有番号对应的视频。`POST`传递 JSON 数据：`[string, string]`。返回番号以其文件的路径的键值对。
* `/api/subtitle` 查询本地是否有番号对应的字幕。`POST`传递 JSON 数据：`{fanhao: string}`。返回字幕的路径。

### Linux 上还需要安装依赖包

* `apt install xclip`
* `apt install libgtk-3-dev libappindicator3-dev`

### 依赖项目

* [atotto/clipboard](https://github.com/atotto/clipboard)
* [getlantern/systray](https://github.com/getlantern/systray)