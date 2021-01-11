package service

import (
	"context"
	"errors"
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

// ErrNotMyself 好友不能是自己
var ErrNotMyself = errors.New("friend must not be yourself")

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

// UpdateFriendStatus 更新好友状态
func (UserService) UpdateFriendStatus(meID uint, updateReq *data.UpdateFriendStatusRequest) (err error) {
	friendID, err := userdao.SelectUserIDByUserLogin(updateReq.FriendUserLogin)
	if err != nil {
		return err
	}
	// 好友不能是自己
	if friendID == meID {
		return ErrNotMyself
	}
	// 发送添加好友请求
	if updateReq.Status == data.Yes && updateReq.Type == data.Send {
		err = friendao.Upsert(meID, friendID, data.Waitting)
	}
	// 回应同意好友请求
	if updateReq.Status == data.Yes && updateReq.Type == data.Receive {
		err = friendao.Update(friendID, meID, data.Added)
	}
	// 回应拒绝好友请求
	if updateReq.Status == data.Decline && updateReq.Type == data.Receive {
		err = friendao.Update(friendID, meID, data.Declined)
	}
	return
}

// GetFriendList 获取好友集合
// 逻辑
//  1. 假设user是发送方(from_user_id) 如果status等于add:已添加 , status==wait 等待对方同意 status==decline 被对方拒绝
//  2. 假设user是发送方(to_user_id) 如果status等于add:已添加 , status==wait 等待同意或拒绝对方请求 status==decline 已拒绝对方
func (UserService) GetFriendList(userID uint) ([]*data.FriendListResponse, error) {
	friendsResp := []*data.FriendListResponse{}
	friends1, err := friendao.SelectByFromUserID(userID)
	if err != nil {
		return nil, err
	}
	// 这里的toUser就是该用户好友 Type是Send
	for _, val := range friends1 {
		userLogin, err := userdao.SelectUserLoginByUserID(val.ToUserID)
		if err != nil {
			return nil, err
		}
		friendsResp = append(friendsResp, &data.FriendListResponse{FriendUserID: val.ToUserID, FriendUserLogin: userLogin, Type: data.Send})
	}
	friends2, err := friendao.SelectByToUserID(userID)
	if err != nil {
		return nil, err
	}
	// 这里的from user就是该用户的好友 Type是Receive
	for _, val := range friends2 {
		userLogin, err := userdao.SelectUserLoginByUserID(val.FromUserID)
		if err != nil {
			return nil, err
		}
		friendsResp = append(friendsResp, &data.FriendListResponse{FriendUserID: val.FromUserID, FriendUserLogin: userLogin, Type: data.Receive})
	}
	// 合成一个
	friends1 = append(friends1, friends2...)
	for index, val := range friends1 {
		friendsResp[index].Status = val.Status
		friendsResp[index].CreateAt = val.CreateAt
		friendsResp[index].UpdateAt = val.UpdateAt
	}
	return friendsResp, nil
}
