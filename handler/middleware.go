package handler

import (
	"net/http"

	"github.com/apex/log"
)

// Authorization token授权
type Authorization struct {
}

// MiddlewareJWTAuthorization 检验用户是否授权
func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Info("MiddlewareJWTAuthorization")
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareRequireAdminPermission 需要有管理员权限
func MiddlewareRequireAdminPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Info("MiddlewareRequireAdminPremission")
		next.ServeHTTP(rw, req)
	})
}
