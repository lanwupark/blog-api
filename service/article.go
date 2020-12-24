package service

import "github.com/lanwupark/blog-api/data"

// ArticleSrrvice 文章服务
type ArticleSrrvice struct{}

// NewArticleSrrvice new
func NewArticleSrrvice() *ArticleSrrvice {
	return &ArticleSrrvice{}
}

// AddArticle 添加文章
func (as *ArticleSrrvice) AddArticle(article *data.Article) {

}
