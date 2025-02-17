package router

import (
	"go-server-base/cmd/server/web"
	"go-server-base/i18n"
	"go-server-base/middleware"
	rou "go-server-base/router"
	"html/template"
	"net/http"

	ginI18n "github.com/gin-contrib/i18n"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func setWebStatic(rootRouter *gin.Engine) {
	rootRouter.StaticFS("/fav", http.FS(web.Favicon))
	rootRouter.GET("/assets/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(web.Assets))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	rootRouter.GET("/", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(web.IndexHtml))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.NoRoute(func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(web.IndexByte)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})
}

func Routers() *gin.Engine {

	Router = gin.Default()

	Router.Use(gzip.Gzip(gzip.DefaultCompression))
	setWebStatic(Router)
	Router.Use(i18n.GinI18nLocalize())
	Router.SetFuncMap(template.FuncMap{
		"Localize": ginI18n.GetMessage,
	})
	Router.Use(middleware.JwtAuth())
	PrivateGroup := Router.Group("/api/v1")
	RagGroup := Router.Group("/api/rag")

	ragRouter := rou.RagRouter{}
	ragRouter.InitRouter(RagGroup)

	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	return Router
}
