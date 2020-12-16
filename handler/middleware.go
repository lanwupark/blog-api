package handler

import (
	"context"
	"net/http"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
)

type userStructKey struct{}

// MiddlewareRequireAuthorization 必须要授权中间件
func MiddlewareRequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Debugf("excute MiddlewareRequireAuthorization method | uri:%s\n", req.RequestURI)
		token, ok := req.Header["Authorization"]
		if !ok {
			// 没有请求头
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(data.NewFailedResponse("Authorization Header Not Found", http.StatusUnauthorized), rw)
			return
		}
		user, err := ParseToken(token[0])
		if err != nil {
			// token 解析失败
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(data.NewFailedResponse("unauthorization", http.StatusUnauthorized), rw)
			return
		}
		// 创建context 将user结构体传给之后需要用的handler
		ctx := context.WithValue(req.Context(), userStructKey{}, user)
		// 赋值新的request
		req = req.WithContext(ctx)
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
