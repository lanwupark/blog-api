package data

import (
	"time"
)

// User 用户
type User struct {
	UserID    uint   `db:"user_id" validate:"required"`
	UserLogin string `db:"user_login" validate:"required"`
	IsAdmin   bool   `json:"IsAdmin" db:"is_admin"`
	Status    string
	CreateAt  time.Time `db:"create_at"`
	UpdateAt  time.Time `db:"update_at"`
}

// Category 文章分类
type Category struct {
	CategoryID int       `json:",omitempty"`
	UserID     uint      `json:"UserID" validate:"required"`
	Name       string    `json:"Name" validate:"required"`
	InDate     string    `json:"-"`
	EditDate   time.Time `json:"-"`
}
