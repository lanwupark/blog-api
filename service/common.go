package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
)

const (
	// ArticlesLeaderBoardSize 排行榜大小
	ArticlesLeaderBoardSize = 100

	// CategoriesLeaderBoardSize 分类排行
	CategoriesLeaderBoardSize = 20

	// RedisRankKeyArticle 文章排行名
	RedisRankKeyArticle = "article_rank"
	// RedisRankKeyUserArticleKey 每个用户的文章排行键
	RedisRankKeyUserArticleKey = "user_article_rank_${user_id}"
	// RedisRankKeyCategoryKey redis分类排行
	RedisRankKeyCategoryKey = "categories_rank"
	// RedisRankKeyCategoryArticleKey redis每个分类的文章排行
	RedisRankKeyCategoryArticleKey = "category_article_rank_${category}"
)

// CommonService 小需求写这里
type CommonService struct{}

// NewCommonService ...
func NewCommonService() *CommonService {
	return &CommonService{}
}

// AddFeedback 添加反馈
func (CommonService) AddFeedback(user *data.TokenClaimsSubject, request *data.FeedbackRequest) {
	feedback := &data.Feedback{
		UserID:      user.UserID,
		Description: request.Description,
		Contact:     request.Contact,
		UserLogin:   user.UserLogin,
		UpdateAt:    time.Now(),
	}
	coll := conn.MongoDB.Collection(data.MongoCollectionFeedback)
	coll.InsertOne(context.TODO(), &feedback)
}

// CalculateSort 计算排行:包括每个用户最火的文章 最火的分类 全局最火的文章 全部存入Redis中
func CalculateSort() {
	log.Info("set categories rank begin")
	articleData, err := articledao.FindAllCalculateData()
	if err != nil {
		panic(err)
	}
	// 查询所有文章的数据
	for _, val := range articleData {
		num1, err := likedao.SelectCountByArticleIDAndType(val.ArticleID, data.Favorite)
		if err != nil {
			panic(err)
		}
		num2, err := likedao.SelectCountByArticleIDAndType(val.ArticleID, data.Star)
		if err != nil {
			panic(err)
		}
		val.FavoriteNumber = int(num1)
		val.StarNumber = int(num2)
	}
	// 按排序规则排序
	sort.Sort(data.ByRule(articleData))
	// 存储排行
	rds := conn.Redis
	userMap := map[uint][]string{}
	categoryMap := map[string][]string{}
	ctx := context.TODO()
	rds.Del(ctx, RedisRankKeyArticle)
	articleIDs := []string{}
	for index, val := range articleData {
		// 只需要ArticlesLeaderBoardSize个
		if index < ArticlesLeaderBoardSize {
			articleIDs = append(articleIDs, strconv.Itoa(int(val.ArticleID)))
		}
		// 存储用户的文章
		userMap[val.UserID] = append(userMap[val.UserID], strconv.Itoa(int(val.ArticleID)))
		// 存储分类的文章
		categories, err := categorydao.SelectNamesByArticleID(val.ArticleID)
		if err != nil {
			panic(err)
		}
		for _, category := range categories {
			categoryMap[category] = append(categoryMap[category], strconv.Itoa(int(val.ArticleID)))
		}
	}
	// 存储文章排行
	if len(articleIDs) > 0 {
		res := rds.RPush(ctx, RedisRankKeyArticle, articleIDs)
		log.Info("set article rank")
		if res.Err() != nil {
			panic(res.Err())
		}
	}

	// 存储每个用户的文章排行
	for userid, val := range userMap {
		key := strings.Replace(RedisRankKeyUserArticleKey, "${user_id}", strconv.Itoa(int(userid)), 1)
		res := rds.Del(ctx, key)
		if res.Err() != nil {
			panic(res.Err())
		}
		res = rds.RPush(ctx, key, val)
		if res.Err() != nil {
			panic(res.Err())
		}
	}
	// 存储每个分类的文章排行
	for category, articles := range categoryMap {
		key := strings.Replace(RedisRankKeyCategoryArticleKey, "${category}", category, 1)
		res := rds.Del(ctx, key)
		if res.Err() != nil {
			panic(res.Err())
		}
		res = rds.RPush(ctx, key, articles)
		if res.Err() != nil {
			panic(res.Err())
		}
	}
	// 存储分类排行
	categories, err := categorydao.SelectMostNames(CategoriesLeaderBoardSize)
	if err != nil {
		panic(err)
	}
	// 存储常用分类
	rds.Del(ctx, RedisRankKeyCategoryKey)
	if len(categories) > 0 {
		res := rds.RPush(ctx, RedisRankKeyCategoryKey, categories)
		if res.Err() != nil {
			panic(res.Err())
		}
	}
	log.Info("set categories rank over")
}
