package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/go-redis/redis/v8"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// LikeArticle 喜欢文章 两种类型 如果没有该文档 会返回mongo.ErrNoDocuments错误
func (ArticleService) LikeArticle(articleID uint64, userID uint, likeType data.LikeType) error {
	if _, err := articledao.SelectOne(articleID); err != nil {
		// 没找到
		if err == mongo.ErrNoDocuments {
			log.WithError(err).Warn("no document")
			return err
		}
		log.WithError(err).Error("find article error")
		return err
	}
	like := &data.Like{
		ArticleID: articleID,
		UserID:    userID,
		Type:      likeType,
	}
	db := conn.DB
	sqlTx := db.MustBegin()
	defer sqlTx.Commit()
	// 先删除之前的
	_, err := sqlTx.Exec("DELETE FROM `like` WHERE article_id=? AND user_id=? AND type=? ", like.ArticleID, like.UserID, like.Type)
	if err != nil {
		sqlTx.Rollback()
		return err
	}
	_, err = sqlTx.Exec("INSERT INTO `like` (article_id,user_id,type) VALUES (?,?,?)", like.ArticleID, like.UserID, like.Type)
	if err != nil {
		sqlTx.Rollback()
	}
	return err
}

// CancelLikeArticle 取消喜欢
func (ArticleService) CancelLikeArticle(articleID uint64, userID uint, likeType data.LikeType) error {
	like := &data.Like{
		ArticleID: articleID,
		UserID:    userID,
		Type:      likeType,
	}
	db := conn.DB
	// 先删除之前的
	_, err := db.Exec("DELETE FROM `like` WHERE article_id=? AND user_id=? AND type=? ", like.ArticleID, like.UserID, like.Type)
	return err
}

// GetArticleDetail 获取文章详情
func (ArticleService) GetArticleDetail(articleID uint64) (*data.ArticleResponse, error) {
	article, err := articledao.SelectOne(articleID)
	if err != nil {
		return nil, err
	}
	articleResp := article.TreeView()
	likes, err := likedao.SelectByArticleID(articleID)
	if err != nil {
		return nil, err
	}
	stars, favoraties := []*data.LikeResponse{}, []*data.LikeResponse{}
	for _, like := range likes {
		likeResp := &data.LikeResponse{
			UserID:    like.UserID,
			UserLogin: like.UserLogin,
			CreateAt:  like.CreateAt,
		}
		if like.Type == data.Favorite {
			favoraties = append(favoraties, likeResp)
		}
		if like.Type == data.Star {
			stars = append(stars, likeResp)
		}
	}
	articleResp.Stars = stars
	articleResp.Favorities = favoraties
	// 增加点击数
	go articledao.HitAddition(articleID)
	return articleResp, nil
}

// DeleteArticleOrComment 删除文章或评论
func (ArticleService) DeleteArticleOrComment(id uint64, userID uint) error {
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	// 先尝试删除文章
	res, err := coll.UpdateOne(context.TODO(), bson.M{
		"articleid": id,
		"userid":    userID,
	}, bson.D{
		{"$set",
			bson.D{
				{"status", data.Deleted},
			},
		},
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.WithError(err).Error("update error")
		return err
	}
	if res.ModifiedCount > 0 {
		log.Info("delete article")
		return nil
	}
	// 再尝试删除评论
	res, err = coll.UpdateOne(context.TODO(), bson.D{
		{"comments.commentid", id},
		{"comments.userid", userID},
	}, bson.D{
		{"$set", bson.D{
			{"comments.$.status", data.Deleted},
		},
		},
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.WithError(err).Error("update error")
		return err
	}
	if res.ModifiedCount > 0 {
		log.Info("delete comment")
	}
	return nil
}

// GetFavoriteList 获取收藏夹
func (articleservice ArticleService) GetFavoriteList(userID uint) ([]*data.ArticleMaintainResponse, error) {
	log.Infof("get user:%d favorite list", userID)
	articleIDs, err := likedao.SelectArticleIDs(userID, data.Favorite)
	if err != nil {
		return nil, err
	}
	articles, err := articledao.Select(articleIDs)
	if err != nil {
		log.WithError(err).Info("get articles error")
		return nil, err
	}
	articilMaintainsResponse := make([]*data.ArticleMaintainResponse, len(articles))
	for index, article := range articles {
		resp, err := articleservice.setArticleMaintainResponse(article)
		if err != nil {
			return nil, err
		}
		// 赋值
		articilMaintainsResponse[index] = resp
	}
	return articilMaintainsResponse, nil
}

// GetArticleMaintain 获取文章大概
func (articleservice ArticleService) GetArticleMaintain(articleID uint64) (*data.ArticleMaintainResponse, error) {
	article, err := articledao.SelectOne(articleID)
	if err != nil {
		return nil, err
	}
	resp, err := articleservice.setArticleMaintainResponse(article)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ArticleService) setArticleMaintainResponse(article *data.Article) (*data.ArticleMaintainResponse, error) {
	resp := &data.ArticleMaintainResponse{
		ArticleID:     article.ArticleID,
		Title:         article.Title,
		Hits:          article.Hits,
		CommentNumber: uint(len(article.Comments)),
		CreateAt:      article.CreateAt,
	}
	// 点赞总数
	starNumber, err := likedao.SelectCountByArticleIDAndType(article.ArticleID, data.Star)
	if err != nil {
		log.WithError(err).Info("get star count error")
		return nil, err
	}
	// 收藏总数
	favoriteNumber, err := likedao.SelectCountByArticleIDAndType(article.ArticleID, data.Favorite)
	if err != nil {
		log.WithError(err).Info("get favorite count error")
		return nil, err
	}
	// 获取最后修改时间
	lastEditDate, lastEditUserID := article.GetLastEditDateAndUserID()
	lastEditUserLogin, err := userdao.SelectUserLoginByUserID(lastEditUserID)
	if err != nil {
		return nil, err
	}
	// 查询该文章的分类
	categories, err := categorydao.SelectNamesByArticleID(article.ArticleID)
	if err != nil {
		return nil, err
	}
	resp.StarNumber = starNumber
	resp.FavoriteNumber = favoriteNumber
	resp.LastEditUserID = lastEditUserID
	resp.LastEditDate = lastEditDate
	resp.LastEditDateString = util.GetIntervalString(lastEditDate, time.Now().UTC()) //获取时间间隔
	resp.LastEditUserLogin = lastEditUserLogin
	resp.Categories = categories
	return resp, nil
}

// GetUsualCategories 获取常用分类
func (ArticleService) GetUsualCategories() ([]string, error) {
	rds := conn.Redis
	res := rds.LRange(context.TODO(), RedisRankKeyCategoryKey, 0, -1)
	return res.Val(), res.Err()
}

// ArticleMaintainQuery 文章大概查询
// 逻辑:
// 如果 category name 和 content都为空 直接从redis里查热门数据article id再查mongo
// 如果只有category name的话 从redis里查出该分类的集合
// 剩下的就直接查mongo
func (articleservice ArticleService) ArticleMaintainQuery(query *data.ArticleMaintainQuery) ([]*data.ArticleMaintainResponse, *data.PageInfo, error) {
	resp := []*data.ArticleMaintainResponse{}
	rds := conn.Redis
	// 设置默认分页数据
	if query.PageIndex <= data.DefaultPageIndex {
		query.PageIndex = data.DefaultPageIndex
	}
	if query.PageSize <= 0 {
		query.PageSize = data.DefaultPageSize
	}
	// 计算skip limit
	skip := (query.PageIndex - 1) * query.PageSize
	end := query.PageIndex*query.PageSize - 1
	limit := query.PageSize
	// result
	var articles []*data.Article
	var err error
	if strings.TrimSpace(query.Content) == "" {
		var res *redis.StringSliceCmd
		// 查热门数据
		if strings.TrimSpace(query.CategoryName) == "" {
			res = rds.LRange(context.TODO(), RedisRankKeyArticle, skip, end)
		} else { // 查某个分类的数据
			key := strings.Replace(RedisRankKeyCategoryArticleKey, "${category}", query.CategoryName, 1)
			res = rds.LRange(context.TODO(), key, skip, end)
		}
		if res.Err() != nil {
			return nil, nil, res.Err()
		}
		strRes := res.Val()
		articleIDs := []uint64{}
		for _, val := range strRes {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, nil, err
			}
			articleIDs = append(articleIDs, uint64(intVal))
		}
		// 查article mongo里 的不是有顺序的
		articles, err = articledao.Select(articleIDs)
		if err != nil {
			return nil, nil, err
		}
	} else {
		filter := bson.D{
			{"status", data.Normal},
			{"$or", bson.A{
				bson.D{{"title", primitive.Regex{Pattern: query.Content, Options: "i"}}},
				bson.D{{"content", primitive.Regex{Pattern: query.Content, Options: "i"}}},
			}}}
		// 如果有category 那么需要找出具体的articles
		if strings.TrimSpace(query.CategoryName) != "" {
			// 直接查mongo
			categoryArticleIDs, err := categorydao.SelectArticleIDsByCategoryName(query.CategoryName)
			if err != nil {
				return nil, nil, err
			}
			filter = append(filter, bson.E{
				Key: "articleid",
				Value: bson.D{
					{"$in", categoryArticleIDs},
				}},
			)
		}
		coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
		findoptions := &options.FindOptions{}
		// 设置分页
		findoptions.SetSkip(skip)
		findoptions.SetLimit(limit)
		// 查article
		cursor, err := coll.Find(context.TODO(), filter, findoptions)
		if err != nil {
			return nil, nil, err
		}
		err = cursor.All(context.TODO(), &articles)
		if err != nil {
			return nil, nil, err
		}
	}
	// 设置 article maintain response
	for _, val := range articles {
		maintain, err := articleservice.setArticleMaintainResponse(val)
		if err != nil {
			return nil, nil, err
		}
		resp = append(resp, maintain)
	}
	// 排序
	sort.Sort(data.ArticleMaintainSortRule(resp))
	return resp, &data.PageInfo{query.PageIndex, query.PageSize}, nil
}
