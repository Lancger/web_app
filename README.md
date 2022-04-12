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

https://www.liwenzhou.com/posts/Go/viper_tutorial/  Go语言配置管理神器——Viper中文教程

## 三、zap日志库

https://www.liwenzhou.com/posts/Go/use_zap_in_gin/  使用zap接收gin框架默认的日志并配置日志归档

## 四、Mysql

```go
package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

// func initDB() (err error) {
// 	dsn := "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
// 	// 也可以使用MustConnect连接不成功就panic
// 	db, err = sqlx.Connect("mysql", dsn)
// 	if err != nil {
// 		fmt.Printf("connect DB failed, err:%v\n", err)
// 		return
// 	}
// 	db.SetMaxOpenConns(20)
// 	db.SetMaxIdleConns(10)
// 	return
// }

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}
```