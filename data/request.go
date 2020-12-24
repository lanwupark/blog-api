package data

// 一些请求的结构体

// AddArticleRequest 添加文章请求
type AddArticleRequest struct {
	Title      string   `validate:"required,min=10"`
	Categories []string `validate:"gt=0"` //切片长度大于0
	Content    string   `validate:"required,min=50"`
}
