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

// SelectMostNames 查询用户最常选的文章分类
func (CategoryDao) SelectMostNames(size int) ([]string, error) {
	db := conn.DB
	var res []string
	if err := db.Select(&res, "SELECT name FROM categories GROUP BY `name` ORDER BY COUNT(0) DESC LIMIT ?", size); err != nil {
		return nil, err
	}
	return res, nil
}

// SelectArticleIDsByCategoryName 根据分类名查询articleID
func (CategoryDao) SelectArticleIDsByCategoryName(name string) (res []uint64, err error) {
	db := conn.DB
	if err = db.Select(&res, "SELECT article_id FROM categories WHERE `name` = ?", name); err != nil {
		return nil, err
	}
	return
}
