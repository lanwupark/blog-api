package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

var (
	validate *validator.Validate
	router   *config.Router
)

func init() {
	validate = validator.New()
	router = config.GetDefaultRouter()
}

// MiddlewareRequireAuthorization 必须要授权中间件
func MiddlewareRequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Debugf("excute MiddlewareRequireAuthorization method | uri:%s\n", req.RequestURI)
		token, ok := req.Header["Authorization"]
		if !ok {
			// 没有请求头
			rw.WriteHeader(http.StatusUnauthorized)
			util.ToJSON(data.NewFailedResponse("Authorization Header Not Found", http.StatusUnauthorized), rw)
			return
		}
		user, err := util.ParseToken(token[0])
		if err != nil {
			log.WithError(err).Error("parse token failed")
			// token 解析失败
			rw.WriteHeader(http.StatusUnauthorized)
			util.ToJSON(data.NewFailedResponse("unauthorization", http.StatusUnauthorized), rw)
			return
		}
		// 创建context 将user结构体传给之后需要用的handler
		ctx := context.WithValue(req.Context(), UserHandler{}, user)
		// 赋值新的request
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareOptionalAuthorization 授权信息是可选的 这样也不会有错
func MiddlewareOptionalAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Debugf("excute MiddlewareOptionalAuthorization method | uri:%s\n", req.RequestURI)
		token, ok := req.Header["Authorization"]
		if ok {
			user, err := util.ParseToken(token[0])
			if err == nil {
				// 创建context 将user结构体传给之后需要用的handler
				ctx := context.WithValue(req.Context(), UserHandler{}, user)
				// 赋值新的request
				req = req.WithContext(ctx)
			}
		}
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareRequireAdminPermission 需要有管理员权限
func MiddlewareRequireAdminPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Debugf("excute MiddlewareRequireAdminPermission method | uri:%s\n", req.RequestURI)
		token, ok := req.Header["Authorization"]
		if !ok {
			// 没有请求头
			rw.WriteHeader(http.StatusUnauthorized)
			util.ToJSON(data.NewFailedResponse("Authorization Header Not Found", http.StatusUnauthorized), rw)
			return
		}
		user, err := util.ParseToken(token[0])
		if err != nil {
			log.WithError(err).Error("parse token failed")
			// token 解析失败
			rw.WriteHeader(http.StatusUnauthorized)
			util.ToJSON(data.NewFailedResponse("unauthorization", http.StatusUnauthorized), rw)
			return
		}
		if !user.IsAdmin {
			// 权限不足
			rw.WriteHeader(http.StatusUnauthorized)
			util.ToJSON(data.NewFailedResponse("unauthorization", http.StatusUnauthorized), rw)
			return
		}
		// 创建context 将user结构体传给之后需要用的handler
		ctx := context.WithValue(req.Context(), UserHandler{}, user)
		// 赋值新的request
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareUserValidation 校验User中间件
func MiddlewareUserValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var user data.User
		validateThen(next, rw, req, nil, &user)
	})
}

// MiddlewareFeedbackValidation 反馈中间件
func MiddlewareFeedbackValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var feedback data.FeedbackRequest
		validateThen(next, rw, req, FeedbackContextKey{}, &feedback)
	})
}

// MiddlewareAddArticleValidation 校验分类中间件
func MiddlewareAddArticleValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var article data.AddArticleRequest
		validateThen(next, rw, req, ArticleHandler{}, &article)
	})
}

// MiddlewareCheckArticleIDValidation 检测文章ID中间件
func MiddlewareCheckArticleIDValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		key := "article_id"
		idStr, ok := mux.Vars(req)[key]
		if !ok {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s can't be null", key))
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s format error", key))
			return
		}
		// 传输
		ctx := context.WithValue(req.Context(), ArticleIDContextKey{}, uint64(id))
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareCheckAlbumIDValidation 检测相册ID中间件
func MiddlewareCheckAlbumIDValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		key := "album_id"
		idStr, ok := mux.Vars(req)[key]
		if !ok {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s can't be null", key))
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s format error", key))
			return
		}
		// 传输
		ctx := context.WithValue(req.Context(), AlbumIDContextKey{}, uint64(id))
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareCheckUserIDValidation 检查用户id中间件
func MiddlewareCheckUserIDValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		key := "user_id"
		idStr, ok := mux.Vars(req)[key]
		if !ok {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s can't be null", key))
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			RespondBadRequest(rw, fmt.Sprintf("uri error: %s format error", key))
			return
		}
		// 传输
		ctx := context.WithValue(req.Context(), UserIDContextKey{}, uint(id))
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareEditArticleValidation 校验编辑文章中间件
func MiddlewareEditArticleValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var article data.AddArticleRequest
		validateThen(next, rw, req, ArticleHandler{}, &article)
	})
}

// MiddlewareAddCommentValidation 校验添加评论中间件
func MiddlewareAddCommentValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var comment data.AddCommentRequest
		validateThen(next, rw, req, CommentContextKey{}, &comment)
	})
}

// MiddlewareLikeArticleValidation 检验
func MiddlewareLikeArticleValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var likeRequest data.LikeArticleRequest
		validateThen(next, rw, req, LikeArticleContextKey{}, &likeRequest)
	})
}

// MiddlewareAddAlbumValidation 检验
func MiddlewareAddAlbumValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var addAlbumRequest data.AddAlbumRequest
		validateThen(next, rw, req, AddAlbumContextKey{}, &addAlbumRequest)
	})
}

// MiddlewareEditAlbumValidtion ...
func MiddlewareEditAlbumValidtion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var editAlbumRequest data.EditAlbumRequest
		validateThen(next, rw, req, EditAlbumContextKey{}, &editAlbumRequest)
	})
}

// MiddlewareUpdateFriendRequestValidtion ...
func MiddlewareUpdateFriendRequestValidtion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var ppdateFriendStatusRequest data.UpdateFriendStatusRequest
		validateThen(next, rw, req, UpdateFriendStatusRequestKey{}, &ppdateFriendStatusRequest)
	})
}

// MiddlewareArticleMaintainQueryValidation 从query parameter里序列化数据
func MiddlewareArticleMaintainQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var articleMaintain data.ArticleMaintainQuery
		err := schema.NewDecoder().Decode(&articleMaintain, req.URL.Query())
		if err != nil {
			RespondBadRequest(rw, err.Error())
			return
		}
		// 传输
		ctx := context.WithValue(req.Context(), ArticleMaintainQueryKey{}, &articleMaintain)
		req = req.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// MiddlewareAdminArticleMaintainQueryValidation 文章查询中间件
func MiddlewareAdminArticleMaintainQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var adminArticleQuery data.AdminArticleQuery
		decodeQueryStringAndDoNextHandler(next, rw, req, AdminArticleQueryKey{}, &adminArticleQuery)
	})
}

// MiddlewareAdminPhotoQueryValidation 照片查询中间件
func MiddlewareAdminPhotoQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var adminPhotoQuery data.AdminPhotoQuery
		decodeQueryStringAndDoNextHandler(next, rw, req, AdminPhotoQueryKey{}, &adminPhotoQuery)
	})
}

// MiddlewareAdminCommentQueryValidation 评论查询中间件
func MiddlewareAdminCommentQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var adminCommentQuery data.AdminCommentQuery
		decodeQueryStringAndDoNextHandler(next, rw, req, AdminCommentQueryKey{}, &adminCommentQuery)
	})
}

// MiddlewareAdminUserQueryValidation 用户查询中间件
func MiddlewareAdminUserQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var adminUserQuery data.AdminUserQuery
		decodeQueryStringAndDoNextHandler(next, rw, req, AdminUserQueryKey{}, &adminUserQuery)
	})
}

// MiddlewareAdminFeedbackQueryValidation 反馈查询中间件
func MiddlewareAdminFeedbackQueryValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var pageInfo data.PageInfo
		decodeQueryStringAndDoNextHandler(next, rw, req, AdminFeedbackQueryKey{}, &pageInfo)
	})
}

// checkID 检查路径上的id
func checkID(rw http.ResponseWriter, req *http.Request, key string) (uint64, bool) {
	idStr, ok := mux.Vars(req)[key]
	if !ok {
		RespondBadRequest(rw, fmt.Sprintf("uri error: %s can't be null", key))
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(rw, fmt.Sprintf("uri error: %s format error", key))
		return 0, false
	}
	return uint64(id), true
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
			json, _ := util.ToJSONString(err)
			util.ToJSON(data.NewFailedResponse(json, http.StatusBadRequest), rw)
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
		msg, _ := util.ToJSONString(errors)
		util.ToJSON(data.NewFailedResponse(msg, http.StatusBadRequest), rw)
		// from here you can create your own error messages in whatever language you wish
		return false
	}

	return true
}

// DeserializeStruct 反序列化结构 true: 反序列化成功 false:反序列化失败 并且回写response
func DeserializeStruct(rw http.ResponseWriter, req *http.Request, s interface{}) bool {
	err := util.FromJSON(s, req.Body)
	// 反序列化
	if err != nil {
		msg := fmt.Sprintf("deserialize %T struct error:%s", s, err)
		log.Warnf(msg)
		rw.WriteHeader(http.StatusBadRequest)
		util.ToJSON(data.NewFailedResponse(msg, http.StatusBadRequest), rw)
		return false
	}
	return true
}

// 重复代码抽取 key为空时 不传值
func validateThen(next http.Handler, rw http.ResponseWriter, req *http.Request, key interface{}, s interface{}) {
	if DeserializeStruct(rw, req, s) && validateStruct(rw, req, s) {
		if key != nil {
			// 创建context 将user结构体传给之后需要用的handler
			ctx := context.WithValue(req.Context(), key, s)
			// 赋值新的request
			req = req.WithContext(ctx)
		}
		next.ServeHTTP(rw, req)
	}
}

// decodeQueryStringAndDoNextHandler 解密然后处理下一个handler
func decodeQueryStringAndDoNextHandler(next http.Handler, rw http.ResponseWriter, req *http.Request, key interface{}, s interface{}) {
	err := schema.NewDecoder().Decode(s, req.URL.Query())
	if err != nil {
		RespondBadRequest(rw, err.Error())
		return
	}
	// 传输
	ctx := context.WithValue(req.Context(), key, s)
	req = req.WithContext(ctx)
	next.ServeHTTP(rw, req)
}
