package router

import (
	v1 "go-server-base/app/api/v1"

	"github.com/gin-gonic/gin"
)

type RagRouter struct{}

func (s *RagRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		baseRouter.GET("*action", baseApi.Rag)
		baseRouter.POST("*action", baseApi.Rag)
		baseRouter.PUT("*action", baseApi.Rag)
		baseRouter.DELETE("*action", baseApi.Rag)
	}
}
