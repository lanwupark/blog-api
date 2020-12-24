package handler

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

// ArticleHandler 文章请求处理器
type ArticleHandler struct{}

// NewArticleHandler 新建
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// AddArticle 添加一篇文章
func (ah *ArticleHandler) AddArticle(rw http.ResponseWriter, req *http.Request) {
	articleRequest := req.Context().Value(ArticleHandler{}).(*data.AddArticleRequest)
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	log.Infof("article:%+v\n", articleRequest)
	log.Infof("user:%+v\n", user)
}

// GetRoutes 实现接口
func (ah *ArticleHandler) GetRoutes() []*config.Route {
	addArticle := &config.Route{
		Method:          http.MethodPost,
		Path:            "/article",
		Handler:         ah.AddArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareAddArticleValidation},
	}
	return []*config.Route{addArticle}
}
