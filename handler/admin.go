package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/lanwupark/blog-api/util"
)

var (
	adminservice = service.NewAdminService()
)

// AdminHanlder ...
type AdminHanlder struct{}

// AdminArticleQueryKey ...
type AdminArticleQueryKey struct{}

// AdminPhotoQueryKey ...
type AdminPhotoQueryKey struct{}

// AdminCommentQueryKey ...
type AdminCommentQueryKey struct{}

// AdminUserQueryKey ...
type AdminUserQueryKey struct{}

// AdminFeedbackQueryKey ...
type AdminFeedbackQueryKey struct{}

// NewAdminHandler new admin handler
func NewAdminHandler() *AdminHanlder {
	return &AdminHanlder{}
}

// ArticleQuery 文章查询
func (AdminHanlder) ArticleQuery(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin article query")
	adminArticleQuery := req.Context().Value(AdminArticleQueryKey{}).(*data.AdminArticleQuery)
	res, err := adminservice.ArticleQuery(adminArticleQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, err)
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	util.ToJSON(res, rw)
}

// ArticleUpdate 文章更新
func (AdminHanlder) ArticleUpdate(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin article update")
	path := mux.Vars(req)
	articleIDStr, _ := path["article_id"]
	articleID, err := strconv.Atoi(articleIDStr)
	if err != nil {
		RespondBadRequest(rw, "article_id format err")
		return
	}
	status, _ := path["status"]
	if err = adminservice.ArticleUpdate(uint64(articleID), data.CommonType(status)); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// PhotoQuery 照片查询
func (AdminHanlder) PhotoQuery(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin photo query")
	adminPhotoQuery := req.Context().Value(AdminPhotoQueryKey{}).(*data.AdminPhotoQuery)
	res, err := adminservice.PhotoQuery(adminPhotoQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, err)
			return
		}
		if err == service.ErrPageIndexOutOfRange {
			RespondBadRequest(rw, err.Error())
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	util.ToJSON(res, rw)
}

// PhotoUpdate 照片更新
func (AdminHanlder) PhotoUpdate(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin article update")
	path := mux.Vars(req)
	fileName, _ := path["file_name"]
	status, _ := path["status"]
	if err := adminservice.PhotoUpdate(fileName, data.CommonType(status)); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// CommentQuery 评论查询
func (AdminHanlder) CommentQuery(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin comment query")
	adminCommentQuery := req.Context().Value(AdminCommentQueryKey{}).(*data.AdminCommentQuery)
	res, err := adminservice.CommentQuery(adminCommentQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, err)
			return
		}
		if err == service.ErrPageIndexOutOfRange {
			RespondBadRequest(rw, err.Error())
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	util.ToJSON(res, rw)
}

// CommentUpdate 评论更新
func (AdminHanlder) CommentUpdate(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin comment update")
	path := mux.Vars(req)
	commentIDStr, _ := path["comment_id"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		RespondBadRequest(rw, "comment_id format err")
		return
	}
	status, _ := path["status"]
	if err = adminservice.CommentUpdate(uint64(commentID), data.CommonType(status)); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// UserQuery 用户查询
func (AdminHanlder) UserQuery(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin user query")
	adminUserQuery := req.Context().Value(AdminUserQueryKey{}).(*data.AdminUserQuery)
	res, err := adminservice.UserQuery(adminUserQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondNotFound(rw, err)
			return
		}
		RespondInternalServerError(rw, err)
		return
	}
	util.ToJSON(res, rw)
}

// UserUpdate 用户更新
func (AdminHanlder) UserUpdate(rw http.ResponseWriter, req *http.Request) {
	log.Info("admin user update")
	path := mux.Vars(req)
	userIDStr, _ := path["user_id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondBadRequest(rw, "user_id format err")
		return
	}
	status, _ := path["status"]
	if err = adminservice.UserUpdate(uint(userID), data.CommonType(status)); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetFeedback 获取用户反馈
func (AdminHanlder) GetFeedback(rw http.ResponseWriter, req *http.Request) {
	pageInfo := req.Context().Value(AdminFeedbackQueryKey{}).(*data.PageInfo)
	resp, err := adminservice.GetFeedback(pageInfo)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	util.ToJSON(resp, rw)
}

// GetRoutes 获取路由配置
func (admin *AdminHanlder) GetRoutes() []*config.Route {
	articleQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/article",
		Handler:         admin.ArticleQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission, MiddlewareAdminArticleMaintainQueryValidation},
	}
	articleUpdate := &config.Route{
		Method:          http.MethodPut,
		Path:            "/admin/article/{article_id:[0-9]+}/{status:[YBD]{1}}",
		Handler:         admin.ArticleUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission},
	}
	photoQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/photo",
		Handler:         admin.PhotoQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission, MiddlewareAdminPhotoQueryValidation},
	}
	photoUpdate := &config.Route{
		Method:          http.MethodPut,
		Path:            "/admin/photo/{file_name:[0-9a-zA-Z]+\\.[a-zA-Z]{3,4}}/{status:[YBD]{1}}",
		Handler:         admin.PhotoUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission},
	}
	commentQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/comment",
		Handler:         admin.CommentQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission, MiddlewareAdminCommentQueryValidation},
	}
	commentUpdate := &config.Route{
		Method:          http.MethodPut,
		Path:            "/admin/comment/{comment_id:[0-9]+}/{status:[YBD]{1}}",
		Handler:         admin.CommentUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission},
	}
	userQuery := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/user",
		Handler:         admin.UserQuery,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission, MiddlewareAdminUserQueryValidation},
	}
	userUpdate := &config.Route{
		Method:          http.MethodPut,
		Path:            "/admin/user/{user_id:[0-9]+}/{status:[YBD]{1}}",
		Handler:         admin.UserUpdate,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission},
	}
	feedback := &config.Route{
		Method:          http.MethodGet,
		Path:            "/admin/feedback",
		Handler:         admin.GetFeedback,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAdminPermission, MiddlewareAdminFeedbackQueryValidation},
	}
	return []*config.Route{articleQuery, articleUpdate, photoQuery, photoUpdate, commentQuery, commentUpdate, userQuery, userUpdate, feedback}
}
