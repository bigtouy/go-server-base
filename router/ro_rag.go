package router

import (
	v1 "go-server-base/app/api/v1"
	"go-server-base/middleware"

	"github.com/gin-gonic/gin"
)

type RagRouter struct{}

func (s *RagRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("")
	baseRouter.Use(middleware.JwtAuth()).Use(middleware.SessionAuth()).Use(middleware.PasswordExpired())

	baseApi := v1.ApiGroupApp.BaseApi
	{
		baseRouter.GET("*action", baseApi.Rag)
		baseRouter.POST("*action", baseApi.Rag)
		baseRouter.PUT("*action", baseApi.Rag)
		baseRouter.DELETE("*action", baseApi.Rag)
	}
}
