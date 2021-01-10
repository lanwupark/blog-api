package data

// LikeDTO 联表
type LikeDTO struct {
	Like
	UserLogin string `db:"user_login"`
}

type ArticleCalculateDTO struct {
	ArticleID uint64
	UserID    uint
	Hits      uint
	Comments  []*Comment
}

// ArticleCalculate 计算文章数据用
type ArticleCalculate struct {
	ArticleID      uint64
	UserID         uint
	Hits           uint
	StarNumber     int
	CommentNumber  int
	FavoriteNumber int
}

var (
	// FavoriteScore 收藏得分
	FavoriteScore = 1000
	// StarScore 喜欢得分
	StarScore = 500
	// CommentScore 评论得分
	CommentScore = 400
	// HitScore 点击得分
	HitScore = 200
)

// ByRule 排序用 排序规则: 收藏数*1000+点赞数*500+评论数*400+点击量*200
type ByRule []*ArticleCalculate

func (a ByRule) Len() int {
	return len(a)
}

func (a ByRule) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRule) Less(i, j int) bool {
	sumI := a[i].FavoriteNumber*FavoriteScore + a[i].StarNumber*StarScore + a[i].CommentNumber*CommentScore + int(a[i].Hits)*HitScore
	sumJ := a[j].FavoriteNumber*FavoriteScore + a[j].StarNumber*StarScore + a[j].CommentNumber*CommentScore + int(a[j].Hits)*HitScore
	return sumI > sumJ
}
