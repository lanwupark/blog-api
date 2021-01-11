package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
)

type AdminHanlder struct{}

func NewAdminHandler() *AdminHanlder {
	return &AdminHanlder{}
}

// ArticleQuery 文章查询
func (AdminHanlder) ArticleQuery(rw http.ResponseWriter, req *http.Request) {

}

// ArticleUpdate 文章更新
func (AdminHanlder) ArticleUpdate(rw http.ResponseWriter, req *http.Request) {

}

// PhotoQuery 照片查询
func (AdminHanlder) PhotoQuery(rw http.ResponseWriter, req *http.Request) {

}

// PhotoUpdate 照片更新
func (AdminHanlder) PhotoUpdate(rw http.ResponseWriter, req *http.Request) {

}

// CommentQuery 评论查询
func (AdminHanlder) CommentQuery(rw http.ResponseWriter, req *http.Request) {

}

// CommentUpdate 评论更新
func (AdminHanlder) CommentUpdate(rw http.ResponseWriter, req *http.Request) {

}

// UserQuery 用户查询
func (AdminHanlder) UserQuery(rw http.ResponseWriter, req *http.Request) {

}

// UserUpdate 用户更新
func (AdminHanlder) UserUpdate(rw http.ResponseWriter, req *http.Request) {

}

// GetRoutes 获取路由配置
func (admin *AdminHanlder) GetRoutes() []*config.Route {
	articleQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/article",
		Handler:         admin.ArticleQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	articleUpdate := &config.Route{
		Method:          http.MethodPost,
		Path:            "/admin/article",
		Handler:         admin.ArticleUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	photoQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/photo",
		Handler:         admin.PhotoQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	photoUpdate := &config.Route{
		Method:          http.MethodPost,
		Path:            "/admin/photo",
		Handler:         admin.PhotoUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	commentQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/comment",
		Handler:         admin.CommentQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	commentUpdate := &config.Route{
		Method:          http.MethodPost,
		Path:            "/admin/comment",
		Handler:         admin.CommentUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	userQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/user",
		Handler:         admin.UserQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	userUpdate := &config.Route{
		Method:          http.MethodPost,
		Path:            "/admin/user",
		Handler:         admin.UserUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	return []*config.Route{articleQuery, articleUpdate, photoQuery, photoUpdate, commentQuery, commentUpdate, userQuery, userUpdate}
}
