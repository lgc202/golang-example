package service

import (
	"context"
)

type IPostService interface {
	GetPostById(ctx context.Context, id string) string
}

var _ IPostService = (*PostService)(nil)

type PostService struct {
}

func (s *PostService) GetPostById(ctx context.Context, id string) string {
	return "欢迎访问博客"
}

func NewPostService() IPostService {
	return &PostService{}
}
