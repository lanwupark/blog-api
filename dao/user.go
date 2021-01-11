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

// Upsert 更新或者插入
func (UserDao) Upsert(gur *data.GithubUserResponse) (*data.User, error) {
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

// SelectUserLoginByUserID 根据用户id搜用户名
func (UserDao) SelectUserLoginByUserID(userID uint) (string, error) {
	db := conn.DB
	var res string
	err := db.Get(&res, "SELECT user_login FROM users WHERE user_id=?", userID)
	return res, err
}

// SelectUserIDByUserLogin 根据用户名搜用户id
func (UserDao) SelectUserIDByUserLogin(userLogin string) (uint, error) {
	db := conn.DB
	var res uint
	err := db.Get(&res, "SELECT user_id FROM users WHERE user_login=?", userLogin)
	return res, err
}

// SelectUserByUserIDAndType 搜索用户
func (UserDao) SelectUserByUserIDAndType(userID uint, status data.CommonType) (*data.User, error) {
	db := conn.DB
	var user data.User
	err := db.Get(&user, "SELECT * FROM users WHERE user_id=? AND status=?", userID, status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
