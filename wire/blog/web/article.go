package web

import (
	"golang-example/wire/blog/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type PostHandler struct {
	serv service.IPostService
}

var PostSet = wire.NewSet(NewPostHandler, service.NewPostService)

func (h *PostHandler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/post/:id", h.GetPostById)
}

func (h *PostHandler) GetPostById(ctx *gin.Context) {
	content := h.serv.GetPostById(ctx, ctx.Param("id"))
	ctx.String(http.StatusOK, content)
}

func NewPostHandler(serv service.IPostService) *PostHandler {
	return &PostHandler{serv: serv}
}
