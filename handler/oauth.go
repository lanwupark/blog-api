package handler

import (
	"net/http"
	"sync"

	"github.com/lanwupark/blog-api/config"
)

var (
	oauthHandler *oAuthHandler
	oauthOnce    sync.Once
)

// OAuthHandler 获取github第三方授权
type oAuthHandler struct{}

// GetOAuthHandlerInstance 获取OAuthHandler单例对象
func GetOAuthHandlerInstance() *oAuthHandler {
	oauthOnce.Do(func() {
		oauthHandler = &oAuthHandler{}
	})
	return oauthHandler
}

// LoginOAuth 获取登录令牌
func (oAuthHandler) LoginOAuth(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Hello World"))
}

func (o *oAuthHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:  http.MethodGet,
		Path:    "/oauth/redirect",
		Handler: o.LoginOAuth,
	}
	return []*config.Route{route}
}
