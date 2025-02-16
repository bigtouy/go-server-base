package router

import (
	v1 "go-server-base/app/api/v1"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("/")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		baseRouter.Any("*action", baseApi.Rag)
	}
}
