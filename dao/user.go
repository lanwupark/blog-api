package dao

import (
	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
)

// UserDao 数据访问层
type UserDao struct{}

// NewUserDao 获取用户DAO实例
func NewUserDao() *UserDao {
	return &UserDao{}
}

// SelectAll 查询所有
func (UserDao) SelectAll() []*data.User {
	db := conn.DB
	// 别使用纯声明的方式 这样切片为nil
	users := []*data.User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	log.Infof("%v", users)
	return users
}
