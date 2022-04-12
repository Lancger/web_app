package main

import (
	"fmt"
	"wep_app/settings"
)

// Go Web开发通用的脚手架模板

func main() {
	// 1、加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	// // 2、初始化日志
	// if err := logger.Init(); err != nil {
	// 	fmt.Printf("init logger failed, err:%v\n", err)
	// 	return
	// }
	// // 3、初始化MySQL连接
	// if err := mysql.Init(); err != nil {
	// 	fmt.Printf("init mysql failed, err:%v\n", err)
	// 	return
	// }
	// // 4、初始化Redis连接
	// if err := redis.Init(); err != nil {
	// 	fmt.Printf("init redis failed, err:%v\n", err)
	// 	return
	// }
	// 5、注册路由

	// 6、启动服务优雅关机
}
