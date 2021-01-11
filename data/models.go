package data

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

const (
	// MongoCollectionArticle 集合
	MongoCollectionArticle = "article"
	// MongoCollectionAlbum 集合
	MongoCollectionAlbum = "album"
	// MongoCollectionFeedback 集合
	MongoCollectionFeedback = "feedback"
)

// LikeType 喜欢的类型
type LikeType string

// CommonType 通用的类型
type CommonType string

// FriendType 好友类型
type FriendType string

const (
	// Star 点赞
	Star LikeType = "S"
	// Favorite 收藏
	Favorite LikeType = "F"
)

const (
	// Normal 正常
	Normal CommonType = "Y"
	// Blocked 拉黑
	Blocked CommonType = "B"
	// Deleted 被删除
	Deleted CommonType = "D"
)

const (
	// Added 已添加
	Added FriendType = "Y"
	// Declined 已拒绝
	Declined FriendType = "D"
	// Waitting 等待中
	Waitting FriendType = "W"
)

// User 用户
type User struct {
	UserID    uint   `db:"user_id"`
	UserLogin string `db:"user_login"`
	IsAdmin   bool   `db:"is_admin"`
	Status    string
	Email     sql.NullString
	Location  sql.NullString
	Blog      sql.NullString
	CreateAt  time.Time `db:"create_at"`
	UpdateAt  time.Time `db:"update_at"`
}

// Category 文章分类
type Category struct {
	CategoryID int
	ArticleID  uint64 `db:"article_id"`
	Name       string
	CreateAt   time.Time `db:"create_at"`
	UpdateAt   time.Time `db:"update_at"`
}

// Like 点赞或收藏
type Like struct {
	LikeID    uint   `db:"like_id"`
	ArticleID uint64 `db:"article_id"`
	UserID    uint   `db:"user_id"`
	Type      LikeType
	CreateAt  time.Time `db:"create_at"`
	UpdateAt  time.Time `db:"update_at"`
}

// Friend 好友列表
type Friend struct {
	FriendsID  uint `db:"friends_id"`
	FromUserID uint `db:"from_user_id"`
	ToUserID   uint `db:"to_user_id"`
	Status     FriendType
	CreateAt   time.Time `db:"create_at"`
	UpdateAt   time.Time `db:"update_at"`
}

// Article 文章 存在mongo里
type Article struct {
	ArticleID uint64
	UserID    uint
	Title     string
	Content   string
	Comments  []*Comment
	Hits      uint
	Status    CommonType
	CreateAt  time.Time
	UpdateAt  time.Time
}

// Comment 嵌套存在mongo里
type Comment struct {
	CommentID uint64     // 评论ID
	ReplyTo   uint64     // 回复评论或者文章的ID
	UserID    uint       //用户ID
	Content   string     //内容
	Status    CommonType //状态
	CreateAt  time.Time
}

// Album 相册
type Album struct {
	AlbumID     uint64
	UserID      uint
	Title       string
	Description string
	CoverName   string
	Location    string
	Hits        uint
	Status      CommonType
	Photos      []*Photo
	CreateAt    time.Time
	UpdateAt    time.Time
}

// Photo 相片
type Photo struct {
	Name         string
	OriginalName string
	FileSize     uint64
	Status       CommonType
	CreateAt     time.Time
	UpdateAt     time.Time
}

// TreeView 将Article存在mongo里的结果转为树形结构 以便更直观的显示  实现:2个map+2个for循环 小小算法 可笑可笑
// 同时  **这里还会过滤一些被删除和被拉黑的评论**
func (article *Article) TreeView() *ArticleResponse {
	res := &ArticleResponse{
		ArticleID: article.ArticleID,
		UserID:    article.UserID,
		Title:     article.Title,
		Content:   article.Content,
		Hits:      article.Hits,
		Status:    article.Status,
		CreateAt:  article.CreateAt,
	}
	// 定义两个map 评论ID map(唯一)和 reply to map (同一条评论可能有多个回复，所以是切片)
	commentIDMap, replyToMap := map[uint64]*CommentResponse{}, map[uint64][]*CommentResponse{}
	for _, comment := range article.Comments {
		// 过滤被拉黑的评论
		if comment.Status != Normal {
			continue
		}
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
		// 过滤被拉黑的评论
		if comment.Status != Normal {
			continue
		}
		// 获取该评论
		commentResponse := commentIDMap[comment.CommentID]
		// 获取回复该评论的slice集合 并将其加入到 该评论的回复集合中
		commentResponse.Replies = replyToMap[comment.CommentID]
	}
	// 回复map中 键与文章的id相等的切片就是最大的子评论
	res.Comments = replyToMap[res.ArticleID]
	return res
}

// GetLastEditDateAndUserID 获取最后次修改的时间 和该用户
// 逻辑:优先级校验 文章更新时间UpdateAt (创建时createAt和updateAt相同)  (因为在mongo里是push的所以评论最新评论始终在最后面) > 文章创建时间
func (article *Article) GetLastEditDateAndUserID() (ltime time.Time, userID uint) {
	ltime = article.UpdateAt
	userID = article.UserID
	comments := article.Comments
	for i := len(comments) - 1; i >= 0; i-- {
		// 正常状态并且大于文章更新时间
		if comments[i].Status == Normal && comments[i].CreateAt.Unix() > ltime.Unix() {
			ltime = comments[i].CreateAt
			userID = comments[i].UserID
			return
		}
	}
	return
}

// DuplicateStructField 复制结构体字段:复制约束 两个名称相同,类型相同
func DuplicateStructField(src interface{}, desc interface{}) error {
	// 判断是否是指针类型
	if reflect.TypeOf(src).Kind() != reflect.Ptr || reflect.TypeOf(desc).Kind() != reflect.Ptr {
		return errors.New("the param shoud be a pointer to the struct type")
	}
	// Elem()会去指针指向的值
	if reflect.TypeOf(src).Elem().Kind() != reflect.Struct || reflect.TypeOf(desc).Elem().Kind() != reflect.Struct {
		return errors.New("the param shoud be a pointer to the struct type")
	}
	// 解指针 获取 value
	srcValue := reflect.ValueOf(src).Elem()
	// 解指针 获取 type
	srcType := reflect.TypeOf(src).Elem()
	//
	descValue := reflect.ValueOf(desc)
	var (
		srcField, descField reflect.StructField
		ok                  bool
	)
	for i := 0; i < srcValue.NumField(); i++ {
		srcField = srcType.Field(i)
		// 判断类型里面有没有 名称相同的 field
		if descField, ok = (descValue).Elem().Type().FieldByName(srcField.Name); !ok {
			continue
		}
		// 判断类型是否相同
		if srcField.Type == descField.Type {
			// 找到该字段value
			srcFielValue := srcValue.Field(i)
			// 解指针 设置指针
			descValue.Elem().FieldByName(srcField.Name).Set(srcFielValue)
		}
	}
	return nil
}

// Feedback 反馈
type Feedback struct {
	UserID      uint
	UserLogin   string //为了后面好看 就不满足范式了
	Description string
	Contact     string
}
