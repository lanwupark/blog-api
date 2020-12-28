package dao

type CategoryDao struct{}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{}
}

// SelectNamesByArticleID 搜索某文章的分类
func (CategoryDao) SelectNamesByArticleID(articleID uint64) ([]string, error) {
	db := conn.DB
	var res []string
	if err := db.Select(&res, "SELECT name FROM categories WHERE article_id=?", articleID); err != nil {
		return nil, err
	}
	return res, nil
}
