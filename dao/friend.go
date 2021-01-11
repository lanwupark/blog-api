package dao

import "github.com/lanwupark/blog-api/data"

// FriendDao friend 表
type FriendDao struct{}

func NewFriendDao() *FriendDao {
	return &FriendDao{}
}

// Upsert 更新或插入
func (FriendDao) Upsert(fromUserID, toUserID uint, status data.FriendType) error {
	db := conn.DB
	var count int
	err := db.Get(&count, "SELECT COUNT(0) FROM friends WHERE from_user_id=? AND to_user_id=?", fromUserID, toUserID)
	if err != nil {
		return err
	}
	if count != 0 {
		_, err = db.Exec("UPDATE friends SET status=? WHERE from_user_id=? AND to_user_id=?", status, fromUserID, toUserID)
	} else {
		_, err = db.Exec("INSERT INTO friends (from_user_id,to_user_id,status) VALUES (?,?,?)", fromUserID, toUserID, status)
	}
	return err
}

// Update 更新数据
func (FriendDao) Update(fromUserID, toUserID uint, status data.FriendType) error {
	_, err := conn.DB.Exec("UPDATE friends SET status=? WHERE from_user_id=? AND to_user_id=?", status, fromUserID, toUserID)
	return err
}

// SelectByFromUserID 通过from_user_id 查询数据
func (FriendDao) SelectByFromUserID(userID uint) ([]*data.Friend, error) {
	db := conn.DB
	var friends []*data.Friend
	if err := db.Select(&friends, "SELECT * FROM friends WHERE from_user_id=?", userID); err != nil {
		return nil, err
	}
	return friends, nil
}

// SelectByToUserID 通过to_user_id 查询数据
func (FriendDao) SelectByToUserID(userID uint) ([]*data.Friend, error) {
	db := conn.DB
	var friends []*data.Friend
	if err := db.Select(&friends, "SELECT * FROM friends WHERE to_user_id=?", userID); err != nil {
		return nil, err
	}
	return friends, nil
}
