# gin web脚手架

# 一、初始化项目

```bash
#项目启用 Go Modules
mkdir web_app && cd web_app

go env -w GO111MODULE=on

go env -w GOPROXY=https://goproxy.cn,direct

#初始化项目go.mod文件
go mod init web_app

#下载项目依赖包
go mod tidy
```