package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
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
	usersRoute := &config.Route{
		Method:          http.MethodGet,
		Path:            "/users",
		Handler:         u.GetUsers,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	userRoute := &config.Route{
		Method:          http.MethodGet,
		Path:            "/user",
		Handler:         u.GetUser,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization},
	}
	return []*config.Route{usersRoute, userRoute}
}

// GetUsers 获取用户
func (UserHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	users := userdao.SelectAll()
	resp := data.NewResultListResponse(users)
	data.ToJSON(resp, rw)
}

// GetUser 获取用户
func (UserHandler) GetUser(rw http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(UserHandler{}).(*data.User)
	data.ToJSON(user, rw)
}
