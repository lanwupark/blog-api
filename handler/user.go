package handler

import (
	"database/sql"
	"errors"
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

// UpdateFriendStatusRequestKey 更新好友请求上下文
type UpdateFriendStatusRequestKey struct{}

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

// UpdateFriendStatus 更新好友状态
func (UserHandler) UpdateFriendStatus(rw http.ResponseWriter, req *http.Request) {
	fromUser := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	updateReq := req.Context().Value(UpdateFriendStatusRequestKey{}).(*data.UpdateFriendStatusRequest)
	if err := userservice.UpdateFriendStatus(fromUser.UserID, updateReq); err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, errors.New("friend user not found"))
			return
		}
		if err == service.ErrNotMyself {
			rw.WriteHeader(http.StatusNotAcceptable)
			resp := data.NewFailedResponse(err.Error(), http.StatusNotAcceptable)
			util.ToJSON(resp, rw)
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetFriendList 获取好友集合
func (UserHandler) GetFriendList(rw http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	resp, err := userservice.GetFriendList(user.UserID)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	data.NewResultListResponse(resp)
	util.ToJSON(resp, rw)
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
		Path:            "/user/{user_id:[0-9]+}",
		Handler:         u.GetUserInfo,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCheckUserIDValidation},
	}
	updateStatusReq := &config.Route{
		Method:          http.MethodPost,
		Path:            "/user/friends",
		Handler:         u.UpdateFriendStatus,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareUpdateFriendRequestValidtion},
	}
	updateStatusReq2 := &config.Route{
		Method:          http.MethodPut,
		Path:            "/user/friends",
		Handler:         u.UpdateFriendStatus,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareUpdateFriendRequestValidtion},
	}
	getFriendList := &config.Route{
		Method:          http.MethodGet,
		Path:            "/user/friends",
		Handler:         u.GetFriendList,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization},
	}
	return []*config.Route{userRoute, getUserInfoRoute, updateStatusReq, updateStatusReq2, getFriendList}
}
