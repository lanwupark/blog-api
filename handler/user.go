package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/lanwupark/blog-api/util"
)

var (
	userservice = service.NewUserService()
)

// UserHandler 处理用户表的接口
type UserHandler struct {
}

// NewUserHandler 新建
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUserSelf 获取用户自身信息
func (UserHandler) GetUserSelf(rw http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	util.ToJSON(user, rw)
}

// GetUserInfo 获取用户信息
func (UserHandler) GetUserInfo(rw http.ResponseWriter, req *http.Request) {
	userid := req.Context().Value(UserIDContextKey{}).(uint)
	userInfo, err := userservice.GetUserInfo(userid)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, fmt.Errorf("user_id:%d not found", userid))
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	res := data.NewResultResponse(userInfo)
	util.ToJSON(res, rw)
}

// GetRoutes 获取该handler下所有路由
func (u *UserHandler) GetRoutes() []*config.Route {
	userRoute := &config.Route{
		Method:          http.MethodGet,
		Path:            "/user",
		Handler:         u.GetUserSelf,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization},
	}
	getUserInfoRoute := &config.Route{
		Method:          http.MethodGet,
		Path:            "/user/{user_id}",
		Handler:         u.GetUserInfo,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCheckUserIDValidation},
	}
	return []*config.Route{userRoute, getUserInfoRoute}
}
