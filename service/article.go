package service

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ()

// ArticleService 文章服务
type ArticleService struct{}

// NewArticleSrrvice new
func NewArticleSrrvice() *ArticleService {
	return &ArticleService{}
}

// AddArticle 添加文章
func (ArticleService) AddArticle(article *data.Article, categories []string) (uint64, error) {
	articleID := util.MustGetNextID()
	article.ArticleID = articleID
	article.Status = data.Normal         // Y正常
	article.Hits = 0                     // 点击数
	article.Comments = []*data.Comment{} //空
	article.CreateAt = time.Now()
	article.UpdateAt = time.Now()
	categs := make([]*data.Category, len(categories))
	for index, name := range categories {
		category := &data.Category{
			ArticleID: articleID,
			Name:      name,
		}
		categs[index] = category
	}
	ctx := context.Background()
	var err error
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	// 事务操作mongo
	_, err = coll.InsertOne(ctx, article)
	if err != nil {
		log.WithError(err).Error("insert article failed")
		return 0, err
	}
	sqlTx := conn.DB.MustBegin()
	defer sqlTx.Commit()
	// sql
	for _, v := range categs {
		_, err = sqlTx.Exec("INSERT INTO categories (article_id,name) VALUES (?,?)", v.ArticleID, v.Name)
		if err != nil {
			log.WithError(err).Error("insert category failed")
			break
		}
	}
	if err != nil {
		log.Warn("rollback transaction")
		sqlTx.Rollback()
	}
	return articleID, err
}

// EditArticle 编辑文章 更新mongo 删除分类
func (ArticleService) EditArticle(article *data.Article, categories []string) error {
	articleID := article.ArticleID
	log.Infof("update article   id:%d", articleID)
	// mongo
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	res, err := coll.UpdateOne(context.TODO(),
		bson.D{
			{"articleid", articleID},
		},
		bson.D{
			{"$set", bson.D{
				{"title", article.Title},
				{"content", article.Content},
				{"updateat", time.Now()},
			}},
		})
	if err != nil {
		return err
	}
	// 没得
	if res.ModifiedCount == 0 {
		return nil
	}
	// sql
	categs := make([]*data.Category, len(categories))
	for index, name := range categories {
		category := &data.Category{
			ArticleID: articleID,
			Name:      name,
		}
		categs[index] = category
	}
	sqlTx := conn.DB.MustBegin()
	defer sqlTx.Commit()
	if _, err = sqlTx.Exec("DELETE FROM categories WHERE article_id = ?", articleID); err != nil {
		// 回滚
		log.WithError(err).Error("delete category error")
		sqlTx.Rollback()
		return err
	}
	for _, v := range categs {
		_, err = sqlTx.Exec("INSERT INTO categories (article_id,name) VALUES (?,?)", v.ArticleID, v.Name)
		if err != nil {
			log.WithError(err).Error("insert category failed")
			break
		}
	}
	if err != nil {
		log.Warn("rollback transaction")
		sqlTx.Rollback()
	}
	return err
}

// AddComment 添加评论
func (ArticleService) AddComment(articleID uint64, commentReq *data.AddCommentRequest) (uint64, error) {
	commentID := util.MustGetNextID()
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	comment := &data.Comment{
		CommentID: commentID,
		ReplyTo:   commentReq.ReplyTo,
		UserID:    commentReq.UserID,
		Content:   commentReq.Content,
		Status:    data.Normal,
		CreateAt:  time.Now(),
	}
	res, err := coll.UpdateOne(context.TODO(),
		bson.D{
			{"articleid", articleID},
		}, bson.M{
			"$push": bson.M{
				"comments": comment,
			},
		})
	if err != nil {
		log.WithError(err).Error("add comment error")
	} else {
		log.Infof("add success modified count:%d", res.ModifiedCount)
	}
	return commentID, err
}

// LikeArticle 喜欢文章 两种类型
func (ArticleService) LikeArticle(articleID uint64, userID uint, likeType data.LikeType) error {
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	res := coll.FindOne(context.TODO(), bson.D{
		{"articleid", articleID},
	})
	// 没找到
	if err := res.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			log.WithError(err).Error("find article error")
			return err
		}
		// 忽略
		return nil
	}
	like := &data.Like{
		ArticleID: articleID,
		UserID:    userID,
		Type:      likeType,
	}
	db := conn.DB
	_, err := db.Exec("INSERT INTO `like` (article_id,user_id,type) VALUES (?,?,?)", like.ArticleID, like.UserID, like.Type)
	return err
}
