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

func init() {

}

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
	route := &config.Route{
		Method:          http.MethodGet,
		Path:            "/users",
		Handler:         u.GetUsers,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareJWTAuthorization},
	}
	return []*config.Route{route}
}

// GetUsers 获取用户
func (userHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	users := userdao.SelectAll()
	data.ToJSON(users, rw)
}
