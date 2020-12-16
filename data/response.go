package data

import "net/http"

// 一些回复的结构体

// GenericResponse 通用返回题
type GenericResponse struct {
	Successed bool   // 是否成功
	Code      int    `json:",omitempty"` //状态码
	Message   string `json:",omitempty"` //消息
}

// ResultListResponse 多结果集返回
type ResultListResponse struct {
	GenericResponse
	PageInfo
	ResultList interface{} `json:",omitempty"` //结果集
}

// ResultResponse 单结果返回
type ResultResponse struct {
	GenericResponse
	Result interface{} `json:",omitempty"` //结果集
}

// PageInfo 分页信息
type PageInfo struct {
	PageSize  uint `json:",omitempty"`
	PageIndex uint `json:",omitempty"`
}

// NewFailedResponse 新的错误回应(带状态码)
func NewFailedResponse(msg string, code int) *GenericResponse {
	return &GenericResponse{
		Successed: false,
		Code:      code,
		Message:   msg,
	}
}

// NewSuccessResponse 成功返回
func NewSuccessResponse() *GenericResponse {
	return &GenericResponse{
		Successed: true,
		Code:      http.StatusOK,
	}
}

// NewResultListResponse 结果集返回
func NewResultListResponse(data interface{}) *ResultListResponse {
	return &ResultListResponse{
		GenericResponse: GenericResponse{
			Successed: true,
			Code:      http.StatusOK,
		},
		ResultList: data,
	}
}

// GithubUserResponse github返回
type GithubUserResponse struct {
	Login     string //登录名
	ID        string `json:"id"`         //用户独有ID
	NodeID    string `json:"node_id"`    //
	AvatarURL string `json:"avatar_url"` //头像url
	URL       string `json:"url"`        //用户数据url
	Blog      string //博客
	Email     string //邮箱
	Localtion string //位置
	Name      string //名称
}

// UserResponse 返回给前端 不一样的json序列化
type UserResponse struct {
	Login     string
	NodeID    string
	AvatarURL string
	URL       string
	Blog      string
	Email     string
	Localtion string
	Name      string
}
