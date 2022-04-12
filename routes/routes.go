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
