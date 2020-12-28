package dao

import "github.com/lanwupark/blog-api/data"

// LikeDao like表 data access
type LikeDao struct{}

func NewLikeDao() *LikeDao {
	return &LikeDao{}
}

// SelectByArticleID 按文章主键搜索
func (LikeDao) SelectByArticleID(articleID uint64) ([]*data.LikeDTO, error) {
	db := conn.DB
	likes := []*data.LikeDTO{}
	err := db.Select(&likes, "SELECT l.*,u.user_login FROM `like` l JOIN users u ON  l.user_id=u.user_id AND article_id=?", articleID)
	if err != nil {
		return nil, err
	}
	return likes, nil
}

// SelectArticleIDs 查询属于该id的集合
func (LikeDao) SelectArticleIDs(userID uint, likeType data.LikeType) ([]uint64, error) {
	db := conn.DB
	res := []uint64{}
	err := db.Select(&res, "SELECT article_id FROM `like` WHERE user_id=? AND type=?", userID, likeType)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SelectCountByArticleIDAndType 按照某一类型 查找该文章的 被点赞或被收藏的数量
func (LikeDao) SelectCountByArticleIDAndType(articleID uint64, likeType data.LikeType) (uint, error) {
	db := conn.DB
	var res uint
	err := db.Get(&res, "SELECT count(0) FROM `like` WHERE article_id=? AND type=?", articleID, likeType)
	return res, err
}
