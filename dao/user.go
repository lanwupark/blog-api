package dao

import (
	"sync"

	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

var (
	userdao *userDao
	once    sync.Once
)

// user数据访问层
type userDao struct{}

// GetUserDaoInstance 获取用户DAO实例
func GetUserDaoInstance() *userDao {
	once.Do(func() {
		userdao = &userDao{}
	})
	return userdao
}

// SelectAll 查询所有
func (userDao) SelectAll() []*data.User {
	db := config.GetDB()
	var users []*data.User
	db.Select(&users, "SELECT * FROM users")
	return users
}
