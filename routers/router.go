package routers

import (
	"caspar/gin-blog/middleware/jwt"
	"caspar/gin-blog/pkg/setting"
	"caspar/gin-blog/routers/api"
	v1 "caspar/gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"
)

/*func Default() *Engine {
	debugPrintWARNINGDefault()
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}*/

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)
	/*r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})*/
	r.GET("auth", api.GetAuth)

	apiv1 := r.Group("api/v1")
	apiv1.Use(jwt.JWT())
	{
		//tags
		apiv1.GET("tags", v1.GetTags)
		apiv1.POST("tags", v1.AddTags)
		apiv1.PUT("tags/:id", v1.EditTag)
		apiv1.DELETE("tags/:id", v1.DeleteTag)
		//argicle
		apiv1.GET("article", v1.GetArticles)
		apiv1.GET("article/:id", v1.GetArticle)
		apiv1.POST("article", v1.AddArticle)
		apiv1.PUT("article/:id", v1.EditArticle)
		apiv1.DELETE("article/:id", v1.DeleteArticle)

	}
	return r
}
