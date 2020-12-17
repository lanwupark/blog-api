package dao

import (
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
	var users []*data.User
	db.Select(&users, "SELECT * FROM users")
	return users
}
