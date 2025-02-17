package server

import (
	"encoding/gob"
	"go-server-base/init/cache"
	"go-server-base/init/log"
	"go-server-base/init/router"
	"go-server-base/init/session"
	"go-server-base/init/session/psession"
	"go-server-base/init/validator"
	"go-server-base/init/viper"

	"github.com/gin-gonic/gin"
)

func Start() {

	viper.Init()
	log.Init()
	//db.Init()
	//migration.Init()
	validator.Init()
	gob.Register(psession.SessionUser{})
	cache.Init()
	session.Init()
	gin.SetMode("debug")

	rootRouter := router.Routers()

	err := rootRouter.Run()
	if err != nil {
		return
	}
}
