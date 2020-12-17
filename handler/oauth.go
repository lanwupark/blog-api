package handler

import (
	"net/http"
	"sync"

	"github.com/lanwupark/blog-api/config"
)

var (
	oauthHandler *OAuthHandler
	oauthOnce    sync.Once
)

// OAuthHandler 获取github第三方授权
type OAuthHandler struct{}

// NewOAuthHandler 获取OAuthHandler单例对象
func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

// LoginOAuth 获取登录令牌
func (OAuthHandler) LoginOAuth(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello World"))
}

// GetRoutes 实现接口
func (o *OAuthHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:  http.MethodGet,
		Path:    "/oauth/redirect",
		Handler: o.LoginOAuth,
	}
	return []*config.Route{route}
}
