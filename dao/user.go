package dao

import (
	"sync"

	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

var (
	conn    = config.GetConnection()
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
	db := conn.DB
	var users []*data.User
	// select * from users
	db.Find(&users)
	return users
}
