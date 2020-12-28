package dao

import (
	"database/sql"

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

// UpSert 更新或者插入
func (UserDao) UpSert(gur *data.GithubUserResponse) (*data.User, error) {
	db := conn.DB
	var user data.User
	stmtx, err := db.Preparex("SELECT * FROM users WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	if err = stmtx.Get(&user, gur.ID); err != nil {
		if err == sql.ErrNoRows {
			//数据库里没有该用户 insert
			db.MustExec("INSERT INTO users (user_id,user_login,email,location,blog) VALUES (?,?,?,?,?)", gur.ID, gur.Login, gur.Email, gur.Localtion, gur.Blog)
		} else {
			return nil, err
		}
	} else {
		// 更新 有时候用户会更改信息
		// 判断是否有改变 有点sb 但是不用动脑O(∩_∩)O
		if user.UserLogin != gur.Login || user.Email.String != gur.Email || user.Location.String != gur.Localtion || user.Blog.String != gur.Blog {
			db.MustExec("UPDATE users SET user_login=?, email=?,location=?, blog=? WHERE user_id=?", gur.Login, gur.Email, gur.Localtion, gur.Blog, user.UserID)
			stmtx.Get(&user, gur.ID)
		}
	}
	return &user, nil
}

// SelectUserLoginByUserId 根据用户id搜用户名
func (UserDao) SelectUserLoginByUserId(userID uint) (string, error) {
	db := conn.DB
	var res string
	err := db.Get(&res, "SELECT user_login FROM users WHERE user_id=?", userID)
	return res, err
}
