package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

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
		user, err := util.ParseToken(token[0])
		if err != nil {
			// token 解析失败
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(data.NewFailedResponse("unauthorization", http.StatusUnauthorized), rw)
			return
		}
		// 创建context 将user结构体传给之后需要用的handler
		ctx := context.WithValue(req.Context(), UserHandler{}, user)
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

// MiddlewareUserValidation 校验User中间件
func MiddlewareUserValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var user data.User
		if deserializeStruct(rw, req, &user) && validateStruct(rw, req, &user) {
			next.ServeHTTP(rw, req)
		}
	})
}

// MiddlewareCategoryValidation 校验分类中间件
func MiddlewareCategoryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var category data.Category
		if deserializeStruct(rw, req, &category) && validateStruct(rw, req, &category) {
			// 创建context 将user结构体传给之后需要用的handler
			ctx := context.WithValue(req.Context(), CategoryHandler{}, &category)
			// 赋值新的request
			req = req.WithContext(ctx)
			next.ServeHTTP(rw, req)
		}
	})
}

// 校验参数是否正确 不正确的话向response写入校验信息
func validateStruct(rw http.ResponseWriter, req *http.Request, s interface{}) bool {
	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(s)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			json, _ := data.ToJSONString(err)
			data.ToJSON(data.NewFailedResponse(json, http.StatusBadRequest), rw)
			return false
		}
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf(
				"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
				err.Namespace(),
				err.Field(),
				err.Tag(),
			)
			errors = append(errors, errMsg)
		}
		msg, _ := data.ToJSONString(errors)
		data.ToJSON(data.NewFailedResponse(msg, http.StatusBadRequest), rw)
		// from here you can create your own error messages in whatever language you wish
		return false
	}

	return true
}

// true: 反序列化成功 false:反序列化失败 并且回写response
func deserializeStruct(rw http.ResponseWriter, req *http.Request, s interface{}) bool {
	err := data.FromJSON(s, req.Body)
	// 反序列化
	if err != nil {
		msg := fmt.Sprintf("deserialize %T struct error:%s", s, err)
		log.Warnf(msg)
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(data.NewFailedResponse(msg, http.StatusBadRequest), rw)
		return false
	}
	return true
}
