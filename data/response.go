package data

// 一些回复的结构体

// GithubUserResponse github返回
type GithubUserResponse struct {
	Login     string //登录名
	ID        string `json:"id"`         //用户独有ID
	NodeID    string `json:"node_id"`    //
	AvatarURL string `json:"avatar_url"` //头像url
	URL       string `json:"url"`        //用户数据url
	Blog      string //博客
	Email     string //邮箱
	Localtion string //位置
	Name      string //名称
}

// UserResponse 返回给前端 不一样的json序列化
type UserResponse struct {
	Login     string
	NodeID    string
	AvatarURL string
	URL       string
	Blog      string
	Email     string
	Localtion string
	Name      string
}
