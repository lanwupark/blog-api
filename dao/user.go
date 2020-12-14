package dao

import (
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

var (
	conn = config.GetConnection()
)

// UserDao user数据访问层
type UserDao struct {
}

// SelectAll 查询所有
func (UserDao) SelectAll() []*data.User {
	db := conn.DB
	var users []*data.User
	// select * from users
	db.Find(&users)
	return users
}
