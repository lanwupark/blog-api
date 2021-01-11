package handler

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/lanwupark/blog-api/util"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	articleservice = service.NewArticleSrrvice()
)

// ArticleHandler 文章请求处理器
type ArticleHandler struct{}

// ArticleMaintainQueryKey 上下文
type ArticleMaintainQueryKey struct{}

// ArticleIDContextKey 上下文
type ArticleIDContextKey struct{}

// UserIDContextKey ...
type UserIDContextKey struct{}

// LikeArticleContextKey 上下文
type LikeArticleContextKey struct{}

// CommentContextKey ...
type CommentContextKey struct{}

// NewArticleHandler 新建
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// AddArticle 添加一篇文章
func (ArticleHandler) AddArticle(rw http.ResponseWriter, req *http.Request) {
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
func (ArticleHandler) AddComment(rw http.ResponseWriter, req *http.Request) {
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
func (ArticleHandler) EditArticle(rw http.ResponseWriter, req *http.Request) {
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
func (ArticleHandler) LikeArticle(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	likeArticleRequest := req.Context().Value(LikeArticleContextKey{}).(*data.LikeArticleRequest)
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	// go in
	if err := articleservice.LikeArticle(id, user.UserID, likeArticleRequest.LikeType); err != nil {
		if err == mongo.ErrNoDocuments {
			RespondNotFound(rw, err)
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// CancelLikeArticle 取消喜欢文章
func (ArticleHandler) CancelLikeArticle(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	// 从路由里面获取
	likeType := data.LikeType(mux.Vars(req)["like_type"])
	// go in
	if err := articleservice.CancelLikeArticle(id, user.UserID, likeType); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetArticleDetail 获取文章详情
func (ArticleHandler) GetArticleDetail(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	// in
	detail, err := articleservice.GetArticleDetail(id)
	if err != nil {
		// 404
		if err == mongo.ErrNoDocuments {
			RespondNotFound(rw, err)
			return
		}
		// 500
		RespondInternalServerError(rw, err)
		return
	}
	res := data.NewResultResponse(detail)
	util.ToJSON(res, rw)
}

// DeleteArticleOrComment 删除某条评论或文章 及其子评论
func (ArticleHandler) DeleteArticleOrComment(rw http.ResponseWriter, req *http.Request) {
	id := req.Context().Value(ArticleIDContextKey{}).(uint64)
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	if err := articleservice.DeleteArticleOrComment(id, user.UserID); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetFavoriteList 获取收藏夹
func (ArticleHandler) GetFavoriteList(rw http.ResponseWriter, req *http.Request) {
	userid := req.Context().Value(UserIDContextKey{}).(uint)
	res, err := articleservice.GetFavoriteList(userid)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewResultListResponse(res)
	util.ToJSON(resp, rw)
}

// GetUsualCategories 获取常用的分类排行
func (ArticleHandler) GetUsualCategories(rw http.ResponseWriter, req *http.Request) {
	res, err := articleservice.GetUsualCategories()
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewResultListResponse(res)
	util.ToJSON(resp, rw)
}

// ArticleMaintainQuery 文章大概查询
func (ArticleHandler) ArticleMaintainQuery(rw http.ResponseWriter, req *http.Request) {
	queryRequest := req.Context().Value(ArticleMaintainQueryKey{}).(*data.ArticleMaintainQuery)
	articleMaintainResp, pageInfo, err := articleservice.ArticleMaintainQuery(queryRequest)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewPageInfoResultListResponse(articleMaintainResp, pageInfo)
	util.ToJSON(resp, rw)
}

//
// -----------------------------------------------------------------------------------
//

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
	deleteArticleOrComment := &config.Route{
		Method:          http.MethodDelete,
		Path:            "/article/comment/{article_id:[0-9]+}", //这里都是雪花算法id
		Handler:         ah.DeleteArticleOrComment,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation},
	}
	likeArticle := &config.Route{
		Method:          http.MethodPost,
		Path:            "/article/like/{article_id:[0-9]+}",
		Handler:         ah.LikeArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation, MiddlewareLikeArticleValidation},
	}
	canelLikeArticle := &config.Route{
		Method:          http.MethodDelete,
		Path:            "/article/like/{article_id:[0-9]+}/{like_type:[SF]{1}}", //like type S或者F
		Handler:         ah.CancelLikeArticle,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckArticleIDValidation},
	}
	getArticle := &config.Route{
		Method:          http.MethodGet,
		Path:            "/article/{article_id:[0-9]+}", //正则判断
		Handler:         ah.GetArticleDetail,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCheckArticleIDValidation},
	}
	getFavoriteList := &config.Route{
		Method:          http.MethodGet,
		Path:            "/article/favorite/{user_id:[0-9]+}", //正则判断
		Handler:         ah.GetFavoriteList,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCheckUserIDValidation},
	}
	getCategories := &config.Route{
		Method:          http.MethodGet,
		Path:            "/article/categories",
		Handler:         ah.GetUsualCategories,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	articleMaintainQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/article/query",
		Handler:         ah.ArticleMaintainQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareArticleMaintainQueryValidation},
	}
	return []*config.Route{articleMaintainQuery, addArticle, editArticle, addComment, deleteArticleOrComment, likeArticle, canelLikeArticle, getArticle, getFavoriteList, getCategories}
}
