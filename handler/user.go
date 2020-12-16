package handler

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
)

var (
	userdao     = dao.GetUserDaoInstance()
	userOnce    sync.Once
	userHanlder *userHandler
)

// userHandler 处理用户表的接口
type userHandler struct {
}

// GetUserHandlerInstance 新建
func GetUserHandlerInstance() *userHandler {
	userOnce.Do(func() {
		userHanlder = &userHandler{}
	})
	return userHanlder
}

// GetRoutes 获取该handler下所有路由
func (u *userHandler) GetRoutes() []*config.Route {
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
func (userHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	users := userdao.SelectAll()
	resp := data.NewResultListResponse(users)
	data.ToJSON(resp, rw)
}

// GetUser 获取用户
func (userHandler) GetUser(rw http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(userStructKey{}).(*data.User)
	data.ToJSON(user, rw)
}
