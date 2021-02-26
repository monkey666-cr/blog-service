package routers

import (
	"blog-service/global"
	"blog-service/internal/middleware"
	v1 "blog-service/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.Tracing())
	r.Use(middleware.Translations())

	demo := v1.NewDemo()
	article := v1.NewArticle()
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/ping", demo.Ping)

		// 获取指定文章
		apiv1.GET("/article/:id", article.Get)
	}
	return r
}
