package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

var (
	userdao = dao.NewUserDao()
)

// UserHandler 处理用户表的接口
type UserHandler struct {
}

// NewUserHandler 新建
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetRoutes 获取该handler下所有路由
func (u *UserHandler) GetRoutes() []*config.Route {
	userRoute := &config.Route{
		Method:          http.MethodGet,
		Path:            "/user",
		Handler:         u.GetUserSelf,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization},
	}
	return []*config.Route{userRoute}
}

// GetUserSelf 获取用户自身信息
func (UserHandler) GetUserSelf(rw http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	util.ToJSON(user, rw)
}
