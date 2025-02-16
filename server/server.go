package server

import (
	"github.com/gin-gonic/gin"
	"go-server-base/init/log"
	"go-server-base/init/router"
	"go-server-base/init/viper"
)

func Start() {

	viper.Init()
	log.Init()
	//db.Init()
	//migration.Init()
	gin.SetMode("debug")

	rootRouter := router.Routers()

	err := rootRouter.Run()
	if err != nil {
		return
	}
}
