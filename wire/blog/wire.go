//go:build wireinject
// +build wireinject

package main

import (
	"golang-example/wire/blog/ioc"
	"golang-example/wire/blog/web"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeApp() *gin.Engine {
	wire.Build(
		// web.NewPostHandler,
		// service.NewPostService,
		web.PostSet,
		ioc.NewGinEngineAndRegisterRoute,
	)
	return &gin.Engine{}
}
