# pc-phone-go

## 可在手机端向 PC 端发送文件、文本，互设剪贴板

## 此为 PC 服务端

## 使用说明

本服务默认监听地址为：http://IP:8800/clip (需将 IP 改为 PC 端内网 IP)
通过 POST 方式向此地址传递数据，表单键值有两个：

1. type 字符串 操作/数据的类型，可选的值为
    1. "getclip"：获取 PC 端剪贴板（此时 content 参数可为空）
    2. "URL"：向 PC 端发送链接，将自动打开此链接
    3. "文本"：向 PC 端发送文本，当字节数小于 512 B 时，复制到剪贴板，否则作为文件保存到用户下载目录
    4. 当没有匹配到上述情况时，将作为文件保存到 用户下载目录
2. content 文件 需保存到 PC 端的文件（保存在用户下载文满中）
