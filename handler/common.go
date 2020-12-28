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
	commonservice = service.NewCommonService()
)

// CommonHandler ...
type CommonHandler struct{}

type FeedbackContextKey struct{}

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

// Feedback 反馈
func (CommonHandler) Feedback(rw http.ResponseWriter, req *http.Request) {
	feedback := req.Context().Value(FeedbackContextKey{}).(*data.FeedbackRequest)
	user := req.Context().Value(UserHandler{})
	var userToken data.TokenClaimsSubject
	if user != nil {
		userToken = *user.(*data.TokenClaimsSubject)
	}
	commonservice.AddFeedback(&userToken, feedback)
	RespondStatusOk(rw)
}

// GetRoutes 实现接口
func (ch *CommonHandler) GetRoutes() []*config.Route {
	generateID := &config.Route{
		Method:          http.MethodGet,
		Path:            "/common/generate_id",
		Handler:         ch.GenerateID,
		MiddlewareFuncs: []mux.MiddlewareFunc{},
	}
	feedback := &config.Route{
		Method:          http.MethodPost,
		Path:            "/feedback",
		Handler:         ch.Feedback,
		MiddlewareFuncs: []mux.MiddlewareFunc{MiddlewareOptionalAuthorization, MiddlewareFeedbackValidation},
	}
	return []*config.Route{generateID, feedback}
}
