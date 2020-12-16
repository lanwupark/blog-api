package handler

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

var (
	categoryOnce    sync.Once
	categoryhandler *categoryHandler
)

type categoryHandler struct{}

// GetCategoryHandlerInstance 新建
func GetCategoryHandlerInstance() *categoryHandler {
	categoryOnce.Do(func() {
		categoryhandler = &categoryHandler{}
	})
	return categoryhandler
}

func (categoryHandler) AddCategory(rw http.ResponseWriter, req *http.Request) {
	catogory := req.Context().Value(categoryHandler{}).(*data.Category)
	data.ToJSON(catogory, rw)
}

func (c *categoryHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:          http.MethodPost,
		Path:            "/category",
		Handler:         c.AddCategory,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareCategoryValidation},
	}
	return []*config.Route{route}
}
