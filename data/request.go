package data

import "time"

// 一些请求的结构体

// MongoCondition 过滤条件
type MongoCondition string

var (
	// MongoTag tag name
	MongoTag = "mongo"

	// MongoEqual 相等
	MongoEqual MongoCondition = "equal"

	// MongoLessThan 小于
	MongoLessThan MongoCondition = "lt"

	// MongoGreatThan 大于
	MongoGreatThan MongoCondition = "gt"

	// MongoLike like
	MongoLike MongoCondition = "like"
)

// AddArticleRequest 添加文章请求
type AddArticleRequest struct {
	Title      string   `validate:"required,min=10"`
	Categories []string `validate:"gt=0"` //切片长度大于0
	Content    string   `validate:"required,min=50"`
}

// AddCommentRequest 添加评论或者恢复
type AddCommentRequest struct {
	ReplyTo uint64 `validate:"required"`
	Content string `validate:"required,min=1"`
	UserID  uint
}

// LikeArticleRequest 喜欢类型
type LikeArticleRequest struct {
	LikeType LikeType `validate:"required,oneof=S F"` //S 和 F
}

// FeedbackRequest 请求
type FeedbackRequest struct {
	Description string `validate:"required"`
	Contact     string
}

// AddAlbumRequest 添加相册请求
type AddAlbumRequest struct {
	AlbumID     uint64 `validate:"required"`
	Title       string `validate:"required,min=10"`
	CoverName   string
	Description string
	Location    string
	PhotoList   []string //要添加的uuid文件名集合
}

// EditAlbumRequest 编辑相册请求
type EditAlbumRequest struct {
	Title           string
	Description     string
	Location        string
	CoverName       string
	DeletePhotoList []string
}

// UpdateFriendType 更新好友状态请求类型
type UpdateFriendType string

const (
	// Send 请求发送方
	Send UpdateFriendType = "S"
	// Receive 请求接收方
	Receive UpdateFriendType = "R"
)

// UpdateFriendStatus 更新好友请求状态
type UpdateFriendStatus string

const (
	// Yes 添加
	Yes UpdateFriendStatus = "Y"
	// Decline 拒绝
	Decline UpdateFriendStatus = "D"
)

// UpdateFriendStatusRequest 更新好友状态请求
type UpdateFriendStatusRequest struct {
	FriendUserLogin string             `validate:"required"`
	Status          UpdateFriendStatus `validate:"required,oneof=Y D"`
	Type            UpdateFriendType   `validate:"required,oneof=S R"`
}

// ArticleMaintainQuery 文章大概查询
type ArticleMaintainQuery struct {
	Content      string `schema:"content"`
	CategoryName string `schema:"category_name"`
	PageInfo
}

// DateInterval 时间区间
type DateInterval struct {
	DateFrom *time.Time `schema:"date_from" mongo:"gt"`
	DateTo   *time.Time `schema:"date_to" mongo:"lt"`
}

// AdminArticleQuery 管理文章查询
type AdminArticleQuery struct {
	ArticleID uint64     `schema:"article_id" mongo:"equal"`
	UserLogin string     `schema:"user_login"`
	UserID    uint       `schema:"-" mongo:"equal"`
	Title     string     `schema:"title" mongo:"like"`
	Content   string     `schema:"content" mongo:"like"`
	Status    CommonType `schema:"status" mongo:"equal"`
	PageInfo
	DateInterval
}

// AdminPhotoQuery ...
type AdminPhotoQuery struct {
	AlbumID   uint64     `schema:"album_id" mongo:"equal"`
	UserLogin string     `schema:"user_login"`
	UserID    uint       `schema:"-" mongo:"equal"`
	Status    CommonType `schama:"status" mongo:"equal"`
	PageInfo
	DateInterval
}

// AdminCommentQuery ...
type AdminCommentQuery struct {
	ArticleID uint64     `schema:"article_id" mongo:"equal"`
	CommentID uint64     `schema:"comment_id"`
	UserLogin string     `schema:"user_login"`
	UserID    uint       `schema:"-" mongo:"equal"`
	Status    CommonType `schema:"status" mongo:"equal"`
	PageInfo
	DateInterval
}

// AdminUserQuery ...
type AdminUserQuery struct {
	UserID    uint       `schema:"user_id"`
	UserLogin string     `schema:"user_login"`
	Status    CommonType `schema:"status"`
	PageInfo
	DateInterval
}
