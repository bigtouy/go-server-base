package router

import (
	rou "go-server-base/router"

	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func Routers() *gin.Engine {

	Router = gin.Default()
	PrivateGroup := Router.Group("/api/v1")
	RagGroup := Router.Group("/api/rag")

	ragRouter := rou.RagRouter{}
	ragRouter.InitRouter(RagGroup)

	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	return Router
}
