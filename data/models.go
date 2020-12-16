package data

import (
	"time"
)

// User 用户
type User struct {
	UserID    uint       `json:"UserId" db:"user_id" validate:"required"`
	UserLogin string     `json:"Login" db:"user_login" validate:"required"`
	IsAdmin   bool       `json:"IsAdmin" db:"is_admin"`
	CreateAt  *time.Time `json:"-" db:"create_at"`
	UpdateAt  *time.Time `json:"-" db:"update_at"`
}

// Category 文章分类
type Category struct {
	CategoryID int       `json:",omitempty"`
	UserID     uint      `json:"UserID" validate:"required"`
	Name       string    `json:"Name" validate:"required"`
	InDate     string    `json:"-"`
	EditDate   time.Time `json:"-"`
}
