package ioc

import (
	"golang-example/wire/blog/web"

	"github.com/gin-gonic/gin"
)

func NewGinEngineAndRegisterRoute(postHandler *web.PostHandler) *gin.Engine {
	engine := gin.Default()
	postHandler.RegisterRoutes(engine)
	return engine
}
