package data

import (
	"time"
)

// User 用户
type User struct {
	UserID    uint       `json:"UserId" db:"user_id"`
	UserLogin string     `json:"Login" db:"user_login"`
	IsAdmin   bool       `json:"IsAdmin" db:"is_admin"`
	CreateAt  *time.Time `json:"-" db:"create_at"`
	UpdateAt  *time.Time `json:"-" db:"update_at"`
}

// Category 文章分类
type Category struct {
	CategoryID int       `json:"CategoryID"`
	UserID     string    `json:"UserID"`
	Name       string    `json:"Name"`
	InDate     string    `json:"InDate"`
	EditDate   time.Time `json:"EditDate"`
}
