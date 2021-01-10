package data

import (
	"net/http"
	"time"
)

// 一些回复的结构体

// TokenClaimsSubject token负荷 主体
type TokenClaimsSubject struct {
	UserID      uint
	UserLogin   string
	IsAdmin     bool
	GithubToken string
}

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

// NewResultResponse 结果返回
func NewResultResponse(data interface{}) *ResultResponse {
	return &ResultResponse{
		GenericResponse: GenericResponse{
			Successed: true,
			Code:      http.StatusOK,
		},
		Result: data,
	}
}

// GithubUserResponse github返回
type GithubUserResponse struct {
	Login     string //登录名
	ID        uint   `json:"id"`         //用户独有ID
	NodeID    string `json:"node_id"`    //
	AvatarURL string `json:"avatar_url"` //头像url
	URL       string `json:"url"`        //用户数据url
	Blog      string `json:",omitempty"` //博客
	Email     string `json:",omitempty"` //邮箱
	Localtion string `json:",omitempty"` //位置
	Name      string `json:",omitempty"` //名称
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

// ArticleResponse TreeView结构的Article
type ArticleResponse struct {
	ArticleID  uint64
	UserID     uint
	Title      string
	Content    string
	Comments   []*CommentResponse
	Hits       uint
	Stars      []*LikeResponse
	Favorities []*LikeResponse
	Status     CommonType `json:"-"`
	CreateAt   time.Time
}

// CommentResponse TreeView结构的Comment
type CommentResponse struct {
	CommentID uint64
	UserID    uint
	Content   string
	Status    CommonType
	Replies   []*CommentResponse
	CreateAt  time.Time
}

// GithubTokenResponse 请求github token的返回json数据
type GithubTokenResponse struct {
	Error            string
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"_,omitempty"`
}

// LikeResponse 喜欢 回复
type LikeResponse struct {
	UserID    uint
	UserLogin string
	CreateAt  time.Time
}

// ArticleMaintainResponse 文章大致查询
type ArticleMaintainResponse struct {
	ArticleID          uint64
	Title              string
	Categories         []string
	LastEditDate       time.Time
	LastEditDateString string
	LastEditUserID     uint
	LastEditUserLogin  string
	Hits               uint
	StarNumber         uint
	CommentNumber      uint
	FavoriteNumber     uint
	CreateAt           time.Time
}

// AddPhotoResponse 添加图片返回
type AddPhotoResponse struct {
	FileName         string
	OriginalFileName string
	FileSize         int
}

// AlbumMaintainResponse 相册简约信息
type AlbumMaintainResponse struct {
	AlbumID   uint64
	AlbumName string `bson:"title"`
	CoverName string
	Location  string
	CreateAt  time.Time
}

// UserInfo 用户页用户信息
type UserInfo struct {
	UserID           uint
	UserLogin        string
	Email            string
	DaysJoined       uint
	FavoritedNumber  uint
	StaredNumber     uint
	CreateAt         time.Time
	ArticleMaintains []*ArticleMaintainResponse
	AlbumMaintains   []*AlbumMaintainResponse
}
