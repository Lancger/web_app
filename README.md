# gin web脚手架

## 一、初始化项目

```go
#项目启用 Go Modules
mkdir web_app && cd web_app

go env -w GO111MODULE=on

go env -w GOPROXY=https://goproxy.cn,direct

#初始化项目go.mod文件
go mod init web_app

#下载项目依赖包
go mod tidy
```

## 二、Viper

```go
	viper.SetConfigFile("config.yaml") // 指定配置文件路径
	// viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	// viper.SetConfigType("yaml")   // 指定配置文件类型(专用于从远程获取配置信息时指定配置文件类型的)
	viper.AddConfigPath(".")   // 指定查找配置文件的路径（这里使用相对路径）
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	// 配置文件热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
```

https://www.liwenzhou.com/posts/Go/viper_tutorial/

## 三、zap日志库

```go

```