package handler

import (
	"fmt"
	"net/http"

	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

// RespondInternalServerError 500
func RespondInternalServerError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	resp := data.NewFailedResponse(fmt.Sprintf("internal server error: %s", err.Error()), http.StatusInternalServerError)
	util.ToJSON(resp, rw)
}

// RespondNotFound 404
func RespondNotFound(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusNotFound)
	resp := data.NewFailedResponse(err.Error(), http.StatusNotFound)
	util.ToJSON(resp, rw)
}

// RespondStatusOk 正常返回
func RespondStatusOk(rw http.ResponseWriter) {
	resp := data.NewSuccessResponse()
	util.ToJSON(resp, rw)
}

// RespondBadRequest 400
func RespondBadRequest(rw http.ResponseWriter, msg string) {
	rw.WriteHeader(http.StatusBadRequest)
	resp := data.NewFailedResponse(msg, http.StatusBadRequest)
	util.ToJSON(resp, rw)
}
