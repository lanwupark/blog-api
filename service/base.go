package service

import (
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
)

var (
	articledao  *dao.ArticleDao
	likedao     *dao.LikeDao
	categorydao *dao.CategoryDao
	userdao     *dao.UserDao
	albumdao    *dao.AlbumDao
	friendao    *dao.FriendDao
	conn        = config.GetConnection()
)

func init() {
	articledao = dao.NewArticleDao()
	likedao = dao.NewLikeDao()
	categorydao = dao.NewCategoryDao()
	userdao = dao.NewUserDao()
	albumdao = dao.NewAlbumDao()
	friendao = dao.NewFriendDao()
}
