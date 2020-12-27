package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

// CommonHandler ...
type CommonHandler struct{}

// NewCommonHandler new
func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

// GenerateID 生成id
func (CommonHandler) GenerateID(rw http.ResponseWriter, req *http.Request) {
	id := util.MustGetNextID()
	resp := data.NewResultResponse(id)
	util.ToJSON(resp, rw)
}

// GetRoutes 实现接口
func (ch *CommonHandler) GetRoutes() []*config.Route {
	generateID := &config.Route{
		Method:          http.MethodGet,
		Path:            "/common/generate_id",
		Handler:         ch.GenerateID,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	return []*config.Route{generateID}
}
