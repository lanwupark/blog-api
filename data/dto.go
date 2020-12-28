package data

// LikeDTO 联表
type LikeDTO struct {
	Like
	UserLogin string `db:"user_login"`
}
