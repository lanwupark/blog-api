package data

import (
	"database/sql"
	"time"
)

// User 用户
type User struct {
	UserID    uint   `db:"user_id" validate:"required"`
	UserLogin string `db:"user_login" validate:"required"`
	IsAdmin   bool   `db:"is_admin"`
	Status    string
	Email     sql.NullString
	Location  sql.NullString
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

// Article 文章 存在mongo里
type Article struct {
	ArticleID uint64
	UserID    uint
	Title     string
	Content   string
	Comments  []*Comment
	Status    string
	CreateAt  time.Time
}

// TreeView 将Article存在mongo里的结果转为树形结构 以便更直观的显示  实现:2个map+2个for循环 小小算法 可笑可笑
func (article *Article) TreeView() *ArticleResponse {
	res := &ArticleResponse{
		ArticleID: article.ArticleID,
		UserID:    article.UserID,
		Title:     article.Title,
		Content:   article.Content,
		Status:    article.Status,
		CreateAt:  article.CreateAt,
	}
	// 定义两个map 评论ID map(唯一)和 reply to map (同一条评论可能有多个回复，所以是切片)
	commentIDMap, replyToMap := map[uint64]*CommentResponse{}, map[uint64][]*CommentResponse{}
	for _, comment := range article.Comments {
		commentResponse := &CommentResponse{
			CommentID: comment.CommentID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			Status:    comment.Status,
			Replies:   []*CommentResponse{},
			CreateAt:  comment.CreateAt,
		}
		// put 元素到map里
		commentIDMap[comment.CommentID] = commentResponse
		// 恢复某一评论或文章的集合
		replyToMap[comment.ReplyTo] = append(replyToMap[comment.ReplyTo], commentResponse)
	}
	// 把ReplyTo加到与CommentResponse CommentID相等的Comments字段下
	for _, comment := range article.Comments {
		// 获取该评论
		commentResponse := commentIDMap[comment.CommentID]
		// 获取回复该评论的slice集合 并将其加入到 该评论的回复集合中
		commentResponse.Replies = replyToMap[comment.CommentID]
	}
	// 回复map中 键与文章的id相等的切片就是最大的子评论
	res.Comments = replyToMap[res.ArticleID]
	return res
}

// Comment 嵌套存在mongo里
type Comment struct {
	CommentID uint64 // 评论ID
	ReplyTo   uint64 // 回复评论或者文章的ID
	UserID    uint   //用户ID
	Content   string //内容
	Status    string //状态
	CreateAt  time.Time
}
