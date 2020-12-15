package data

import (
	"time"
)

// User 用户
type User struct {
	UserID    uint      `json:"UserId" gorm:"primaryKey"`
	UserLogin string    `json:"Login"`
	IsAdmin   bool      `json:"IsAdmin"`
	CreateAt  time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

// Category 文章分类
type Category struct {
	CategoryID int       `json:"CategoryID"`
	UserID     string    `json:"UserID"`
	Name       string    `json:"Name"`
	InDate     string    `json:"InDate"`
	EditDate   time.Time `json:"EditDate"`
}
