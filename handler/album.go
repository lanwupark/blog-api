package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/lanwupark/blog-api/util"
)

var (
	albumservice = service.NewAlbumService()
)

// AlbumHander handler
type AlbumHander struct{}

// NewAlbumHander 相册处理器
func NewAlbumHander() *AlbumHander {
	return &AlbumHander{}
}

// AlbumIDContextKey 上下文
type AlbumIDContextKey struct{}

// AddAlbumContextKey 上下文
type AddAlbumContextKey struct{}

// AddPhoto 添加照片
func (album *AlbumHander) AddPhoto(rw http.ResponseWriter, req *http.Request) {
	albumID := req.Context().Value(AlbumIDContextKey{}).(uint64)
	// 在中间件里获取
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	fileName, ok := mux.Vars(req)["file_name"]
	if !ok || fileName == "" {
		RespondBadRequest(rw, "file name format err")
		return
	}
	addPhotoResp, err := albumservice.AddPhoto(user.UserID, albumID, fileName, req.Body)
	if err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	resp := data.NewResultResponse(addPhotoResp)
	util.ToJSON(resp, rw)
}

// NewAlbum 添加相册
func (album *AlbumHander) NewAlbum(rw http.ResponseWriter, req *http.Request) {
	// 在中间件里获取
	user := req.Context().Value(UserHandler{}).(*data.TokenClaimsSubject)
	albumReq := req.Context().Value(AddAlbumContextKey{}).(*data.AddAlbumRequest)
	if err := albumservice.NewAlbum(user.UserID, albumReq); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// CancelNewAlbum 取消新建相册
func (album *AlbumHander) CancelNewAlbum(rw http.ResponseWriter, req *http.Request) {
	albumID := req.Context().Value(AlbumIDContextKey{}).(uint64)
	if err := albumservice.CancelNewAlbum(albumID); err != nil {
		RespondInternalServerError(rw, err)
		return
	}
	RespondStatusOk(rw)
}

// GetRoutes 获取路由
func (album *AlbumHander) GetRoutes() []*config.Route {
	addAlbum := &config.Route{
		Method:          http.MethodPost,
		Path:            "/album",
		Handler:         album.NewAlbum,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareAddAlbumValidation},
	}
	addPhoto := &config.Route{
		Method:          http.MethodPost,
		Path:            "/album/photo/{album_id:[0-9]+}/{file_name:[a-zA-Z]+\\.[a-zA-Z]{3,4}}", //正则判断
		Handler:         album.AddPhoto,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckAlbumIDValidation},
	}
	cancelNewAlbum := &config.Route{
		Method:          http.MethodDelete,
		Path:            "/album/cancel/{album_id:[0-9]+}",
		Handler:         album.CancelNewAlbum,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareRequireAuthorization, MiddlewareCheckAlbumIDValidation},
	}
	return []*config.Route{addAlbum, addPhoto, cancelNewAlbum}
}
