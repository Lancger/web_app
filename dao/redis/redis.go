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

// rdb小写实例，没有对外暴露，这里可以通过封装一个对外的方法用于关闭释放redis连接
func Close() {
	_ = rdb.Close()
}
