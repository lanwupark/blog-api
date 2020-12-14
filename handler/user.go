package handler

import (
	"net/http"

	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
)

var (
	userdao dao.UserDao
)

func init() {

}

// UserHandler 处理用户表的接口
type UserHandler struct {
}

// NewUserHandler 新建
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetRoutes 获取该handler下所有路由
func (u *UserHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: u.GetUsers,
	}
	return []*config.Route{route}
}

// GetUsers 获取用户
func (UserHandler) GetUsers(rw http.ResponseWriter, req *http.Request) {
	users := userdao.SelectAll()
	data.ToJSON(users, rw)
}
