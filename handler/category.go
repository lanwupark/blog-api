package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
)

var (
	categorydao = dao.NewCategoryDao()
)

// CategoryHandler category hanlder
type CategoryHandler struct{}

// NewCategoryHandler 新建
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// AddCategory 添加分类
func (CategoryHandler) AddCategory(rw http.ResponseWriter, req *http.Request) {
	catogory := req.Context().Value(CategoryHandler{}).(*data.Category)
	id, err := categorydao.InsertOneToMongo(catogory)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(data.NewFailedResponse(err.Error(), http.StatusInternalServerError), rw)
	}
	data.ToJSON(data.NewResultResponse(id), rw)
}

// GetRoutes 实现接口
func (c *CategoryHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:          http.MethodPost,
		Path:            "/category",
		Handler:         c.AddCategory,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCategoryValidation},
	}
	return []*config.Route{route}
}
