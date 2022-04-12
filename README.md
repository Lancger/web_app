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

## 二、引入viper

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

## 三、引入zap日志库

https://www.liwenzhou.com/posts/Go/use_zap_in_gin/  使用zap接收gin框架默认的日志并配置日志归档

## 四、引入Mysql

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

## 五、引入redis

```go
package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// // 初始化连接
// func initClient() (err error) {
// 	rdb = redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	_, err = rdb.Ping().Result()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// 初始化连接
func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
```

https://www.liwenzhou.com/posts/Go/go_redis/  Go语言操作Redis

## 六、注册路由

- 编写路由

```go
package routes

import (
	"net/http"
	"wep_app/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// func Setup() *gin.Engine {
// 	r := gin.New()
// 	r.Use(logger.GinLogger(), logger.GinRecovery(true))

// 	r.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "OK")
// 	})
// 	return r
// }

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})
	return r
}
```

- 引入路由

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wep_app/dao/mysql"
	"wep_app/dao/redis"
	"wep_app/logger"
	"wep_app/routes"
	"wep_app/settings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Go Web开发通用的脚手架模板

func main() {
	// 1、加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2、初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 3、初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}

	// 4、初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}

	// 5、注册路由
	// r := routes.Setup()
	r := routes.Setup(gin.DebugMode)

	// 6、启动服务优雅关机
	srv := &http.Server{
		// Addr:    ":8080",
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log.Fatalf("listen: %s\n", err)
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	// log.Println("Shutdown Server ...")
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		// log.Fatal("Server Shutdown: ", err)
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	// log.Println("Server exiting")
	zap.L().Info("Server exiting")
}
```

https://www.liwenzhou.com/posts/Go/graceful_shutdown/  优雅地关机或重启

## 七、运行测试

```go
➜  web_app git:(main) ✗ go run main.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /version                  --> wep_app/routes.Setup.func1 (3 handlers)

➜  ~ curl "http://127.0.0.1:8080/version"
v0.1.4
```

## 八、优化，连接关闭

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wep_app/dao/mysql"
	"wep_app/dao/redis"
	"wep_app/logger"
	"wep_app/routes"
	"wep_app/settings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Go Web开发通用的脚手架模板

func main() {
	// 1、加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2、初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// 缓存区的日志追加到日志中
	defer zap.L().Sync()

	// 3、初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	// mysql连接释放
	defer mysql.Close()

	// 4、初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	// redis连接释放
	defer redis.Close()

	// 5、注册路由
	// r := routes.Setup()
	r := routes.Setup(gin.DebugMode)

	// 6、启动服务优雅关机
	srv := &http.Server{
		// Addr:    ":8080",
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log.Fatalf("listen: %s\n", err)
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	// log.Println("Shutdown Server ...")
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		// log.Fatal("Server Shutdown: ", err)
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	// log.Println("Server exiting")
	zap.L().Info("Server exiting")
}
```