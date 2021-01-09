package dao

import (
	"context"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
	"go.mongodb.org/mongo-driver/bson"
)

// ArticleDao ...
type ArticleDao struct{}

// NewArticleDao ...
func NewArticleDao() *ArticleDao {
	return &ArticleDao{}
}

// SelectOne 查询一个 正常状态的
func (ArticleDao) SelectOne(articleID uint64) (*data.Article, error) {
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	var article data.Article
	res := coll.FindOne(context.TODO(), bson.M{"articleid": articleID, "status": data.Normal})
	if res.Err() != nil {
		return nil, res.Err()
	}
	if err := res.Decode(&article); err != nil {
		return nil, err
	}
	return &article, nil
}

// HitAddition 点击数增加
func (ArticleDao) HitAddition(articleID uint64) {
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"articleid": articleID}, bson.D{
		{"$inc", bson.D{
			{"hits", 1},
		}},
	})
	if err != nil {
		log.WithError(err).Errorf("add hit error for: %d", articleID)
	}
}

// Select 批量查询文档
func (ArticleDao) Select(articleIDs []uint64) ([]*data.Article, error) {
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	articles := []*data.Article{}
	filter := bson.D{
		{"articleid", bson.D{
			{"$in", articleIDs},
		}},
	}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &articles); err != nil {
		return nil, err
	}
	return articles, nil
}
