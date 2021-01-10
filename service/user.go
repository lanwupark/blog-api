package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/lanwupark/blog-api/data"
)

// UserService ...
type UserService struct{}

// NewUserService new
func NewUserService() *UserService {
	return &UserService{}
}

// GetUserInfo 获取用户信息
func (UserService) GetUserInfo(userID uint) (*data.UserInfo, error) {
	user, err := userdao.SelectUserByUserIDAndType(userID, data.Normal)
	if err != nil {
		return nil, err
	}
	// star 总数
	num1, err := likedao.SelectCountByUserIDAndType(userID, data.Star)
	if err != nil {
		return nil, err
	}
	// favorite 总数
	num2, err := likedao.SelectCountByUserIDAndType(userID, data.Favorite)
	if err != nil {
		return nil, err
	}
	// 加入天数
	daysJoined := uint(time.Now().Sub(user.CreateAt).Hours() / 24)
	userinfo := &data.UserInfo{
		UserID:          user.UserID,
		UserLogin:       user.UserLogin,
		Email:           user.Email.String,
		DaysJoined:      daysJoined,
		StaredNumber:    num1,
		FavoritedNumber: num2,
		CreateAt:        user.CreateAt,
	}
	// 获取文章简要信息 从redis找出
	rds := conn.Redis
	key := strings.Replace(RedisRankKeyUserArticleKey, "${user_id}", strconv.Itoa(int(user.UserID)), 1)
	sliceCmd := rds.LRange(context.TODO(), key, 0, -1)
	if sliceCmd.Err() != nil {
		return nil, sliceCmd.Err()
	}
	articles := []uint64{}
	for _, val := range sliceCmd.Val() {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		articles = append(articles, uint64(intVal))
	}
	articleMaintains := make([]*data.ArticleMaintainResponse, len(articles))
	articleservice := NewArticleSrrvice()
	// 查询文章简要信息
	for index, val := range articles {
		resp, err := articleservice.GetArticleMaintain(val)
		if err != nil {
			return nil, err
		}
		articleMaintains[index] = resp
	}
	albumresp, err := albumdao.FindMaintainByUserID(user.UserID)
	if err != nil {
		return nil, err
	}

	userinfo.ArticleMaintains = articleMaintains
	userinfo.AlbumMaintains = albumresp
	return userinfo, nil
}
