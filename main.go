package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "oauthServer/docs"
	"oauthServer/handler"
	"oauthServer/pkg/ginServer"
)

// @title Swagger API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

//////@license.name Apache 2.0
/////@license.url http://www.apache.org/licenses/LICENSE-2.0.html

///// @host petstore.swagger.io
///// @BasePath /
func main() {
	ginServer.Default()
	g := gin.Default()
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // The url pointing to API definition
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// wolf-rbac插件 适配apisix网关的接口
	wolf := g.Group("/wolf/rbac")
	{
		//核心接口，apisix核心鉴权能力
		wolf.GET("/access_check", handler.AccessCheckHandler)
		//登陆调用 可以不实现,如果不实现则无法使用插件对应的login
		wolf.POST("/login.rest", func(c *gin.Context) {
			ginServer.HandleTokenRequest(c)
		})
		//获取用户信息 可以不实现,如果不实现则无法使用插件对应的user_info
		wolf.GET("/user_info", func(c *gin.Context) {

		})
	}

	// 标准oauth2 端点可以不实现
	auth := g.Group("/oauth2")
	{
		//token 生成端点,注意：token生成端点创建的token不会销毁之前的token
		auth.GET("/token_endpoint", ginServer.HandleTokenRequest)
		auth.POST("/token_endpoint", ginServer.HandleTokenRequest)
		//token authorization认证端点
		auth.GET("/authorization_endpoint", ginServer.HandleAuthorizeRequest)
		auth.POST("/authorization_endpoint", ginServer.HandleAuthorizeRequest)
		//校验token
		auth.GET("/check_token", ginServer.ValidationBearerToken)
		auth.POST("/check_token", ginServer.ValidationBearerToken)
	}
	g.Run(":8000")
}
