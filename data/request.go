package data

// 一些请求的结构体

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
