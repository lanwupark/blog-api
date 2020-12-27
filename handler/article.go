package handler

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/lanwupark/blog-api/util"
)

var (
	articleservice = service.NewArticleSrrvice()
)

// ArticleHandler 文章请求处理器
type ArticleHandler struct{}

// ArticleIDContextKey 上下文
type ArticleIDContextKey struct{}

// LikeArticleContextKey 上下文
type LikeArticleContextKey struct{}

// CommentContextKey ...
type CommentContextKey struct{}

// NewArticleHandler 新建
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// AddArticle 添加一篇文章
func (ah *ArticleHandler) AddArticle(rw http.ResponseWriter, req *http.Request) {
	// 在中间件里获取
	articleRequest := req.Context().Value(ArticleHandler{}).(*data.AddArticleRequest)
	// 在中间件里获取
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	var article data.Article
	// 复制 title content
	err := data.DuplicateStructField(articleRequest, &article)
	if err != nil {
		log.Error(err.Error())
		RespondInternalServerError(rw, err)
		return
	}
	categories := articleRequest.Categories
	// 设置UserID
	article.UserID = user.UserID
	id, err := articleservice.AddArticle(&article, categories)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewResultResponse(id)
	util.ToJSON(resp, rw)
}

// AddComment 添加评论
func (ah *ArticleHandler) AddComment(rw http.ResponseWriter, req *http.Request) {
	log.Info("add comment")
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	comment := req.Context().Value(CommentContextKey{}).(*data.AddCommentRequest)
	// 在中间件里获取
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	// 设置user ID
	comment.UserID = user.UserID
	// 添加评论
	commentID, err := articleservice.AddComment(id, comment)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewResultResponse(commentID)
	util.ToJSON(resp, rw)
}

// EditArticle 编辑文章
func (ah *ArticleHandler) EditArticle(rw http.ResponseWriter, req *http.Request) {
	log.Info("edit article")
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	// 在中间件里获取
	articleRequest := req.Context().Value(ArticleHandler{}).(*data.AddArticleRequest)
	// 在中间件里获取
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	var article data.Article
	// 复制 title content
	err := data.DuplicateStructField(articleRequest, &article)
	if err != nil {
		log.Error(err.Error())
		RespondInternalServerError(rw, err)
		return
	}
	categories := articleRequest.Categories
	// 设置UserID
	article.UserID = user.UserID
	// 设置articleID
	article.ArticleID = id
	err = articleservice.EditArticle(&article, categories)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// LikeArticle 喜欢文章
func (ah *ArticleHandler) LikeArticle(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	likeArticleRequest := req.Context().Value(LikeArticleContextKey{}).(*data.LikeArticleRequest)
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	// go in
	if err := articleservice.LikeArticle(id, user.UserID, likeArticleRequest.LikeType); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetRoutes 实现接口
func (ah *ArticleHandler) GetRoutes() []*config.Route {
	addArticle := &config.Route{
		Method:          http.MethodPost,
		Path:            "/article",
		Handler:         ah.AddArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareAddArticleValidation},
	}
	editArticle := &config.Route{
		Method:          http.MethodPut,
		Path:            "/article/{article_id:[0-9]+}", //正则判断
		Handler:         ah.EditArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation, MiddlewareEditArticleValidation},
	}
	addComment := &config.Route{
		Method:          http.MethodPost,
		Path:            "/article/comment/{article_id:[0-9]+}",
		Handler:         ah.AddComment,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation, MiddlewareAddCommentValidation},
	}
	likeArticle := &config.Route{
		Method:          http.MethodPost,
		Path:            "/article/like/{article_id:[0-9]+}",
		Handler:         ah.LikeArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation, MiddlewareLikeArticleValidation},
	}
	return []*config.Route{addArticle, editArticle, addComment, likeArticle}
}
