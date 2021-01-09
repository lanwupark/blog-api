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
