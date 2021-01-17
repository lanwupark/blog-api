package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/dao"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
)

var (
	oauthHandler         *OAuthHandler
	oauthOnce            sync.Once
	tokenRequestTemplate *template.Template //请求模板
	httpclient           *http.Client       //http client
	userdao              = dao.NewUserDao()
)

// 应用配置
var properties *config.ApplicationProperties = &config.GetConfigs().ApplicationProperties

const tokenRequestURLTemplate = `https://github.com/login/oauth/access_token?client_id={{.ClientID}}&client_secret={{.ClientSecret}}&code={{.AuthCode}}`

// github user api
const userAPI = `https://api.github.com/user`

func init() {
	// 初始化请求模板
	tokenRequestTemplate = template.Must(template.New("tokenRequest").Parse(tokenRequestURLTemplate))
	// 初始化http client
	httpclient = new(http.Client)
}

type tokenRequestParam struct {
	ClientID     string
	ClientSecret string
	AuthCode     string
}

func newTokenRequestParam(authCode string) *tokenRequestParam {
	return &tokenRequestParam{properties.OAuthClientID, properties.OAuthClientSecret, authCode}
}

// OAuthHandler 获取github第三方授权
type OAuthHandler struct{}

// NewOAuthHandler 获取OAuthHandler单例对象
func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

// LoginOAuth 获取登录令牌
// 1. 获取查询参数code
// 2. 向gitub请求获取token
// 3. 向github api查询用户数据
// 4. 查询数据库 添加数据或更新
// 5. 重定向回写JWT  url: {{referer}}/oauth/token?set_token={{blog-token}}&github_token={{github-token}}
func (OAuthHandler) LoginOAuth(rw http.ResponseWriter, req *http.Request) {
	if !parametersCheck(rw, req) { //1
		return
	}
	githubTokenResponse, err := tokenRequestGet(rw, req) //2
	if err != nil {
		panic(err)
	}
	githubUserResponse, err := userInfoRequestGet(rw, req, githubTokenResponse) //3
	if err != nil {
		panic(err)
	}
	// ---------------------4---------------------
	user, err := userdao.Upsert(githubUserResponse)
	if err != nil {
		panic(err)
	}
	// ---------------------5---------------------
	tokenSubject := &data.TokenClaimsSubject{
		UserID:      user.UserID,
		UserLogin:   user.UserLogin,
		IsAdmin:     user.IsAdmin,
		GithubToken: githubTokenResponse.AccessToken,
	}
	token, err := util.CreateToken(tokenSubject)
	if err != nil {
		panic(err)
	}
	referer := req.Referer()
	referer = strings.TrimSuffix(referer, "/")
	redirectURL := fmt.Sprintf("%s/oauth/token?set_token=%s&github_token=%s", referer, token, githubTokenResponse.AccessToken)
	log.Infof("url:%s", redirectURL)
	// 删除 content type 底层代码会加上text/html
	rw.Header().Del("Content-Type")
	// 重定向
	http.Redirect(rw, req, redirectURL, http.StatusFound)
}

// GetRoutes 实现接口
func (o *OAuthHandler) GetRoutes() []*config.Route {
	route := &config.Route{
		Method:  http.MethodGet,
		Path:    "/oauth/redirect",
		Handler: o.LoginOAuth,
	}
	return []*config.Route{route}
}

func acceptJSONHeader(req *http.Request) {
	req.Header.Add("Accept", "application/json")
}

func mustNewGet(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func clientDo(req *http.Request) *http.Response {
	resp, err := httpclient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

// 参数校验 必须要有code和referer
func parametersCheck(rw http.ResponseWriter, req *http.Request) bool {
	// referer url  必须要有referer
	if req.Referer() == "" {
		rw.WriteHeader(http.StatusForbidden)
		util.ToJSON(data.NewFailedResponse("are you serious?", http.StatusBadRequest), rw)
		return false
	}
	code := req.URL.Query().Get("code")
	// 没有或者长度为0 或者为""
	if code == "" {
		rw.WriteHeader(http.StatusBadRequest)
		util.ToJSON(data.NewFailedResponse("param a does not exist", http.StatusBadRequest), rw)
		return false
	}
	return true
}

// 获取 github token 请求
func tokenRequestGet(rw http.ResponseWriter, req *http.Request) (*data.GithubTokenResponse, error) {
	code := req.URL.Query().Get("code")
	tokenRequestParam := newTokenRequestParam(code)
	urlBuf := new(bytes.Buffer)
	err := tokenRequestTemplate.Execute(urlBuf, tokenRequestParam) //解析模板
	if err != nil {
		return nil, err
	}
	tokenRequestURL := urlBuf.String()
	log.Infof("token request url:%s \t referer:%s", tokenRequestURL, req.Referer())
	// get request to github
	tokenReq := mustNewGet(tokenRequestURL)
	acceptJSONHeader(tokenReq)
	tokenResp := clientDo(tokenReq)
	bytes, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, err
	}
	respJSON := string(bytes)
	log.Infof("github response status:%d  body:%s", tokenResp.StatusCode, respJSON)
	var githubTokenResponse data.GithubTokenResponse
	if err = util.FromJSONString(respJSON, &githubTokenResponse); err != nil {
		return nil, err
	}
	// 获取token失败
	if githubTokenResponse.AccessToken == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte(respJSON))
		return nil, errors.New("authorized error")
	}
	return &githubTokenResponse, nil
}

// userInfoRequestGet 获取用户资料请求
func userInfoRequestGet(rw http.ResponseWriter, req *http.Request, tokenResponse *data.GithubTokenResponse) (*data.GithubUserResponse, error) {
	getUserRequest := mustNewGet(userAPI)
	acceptJSONHeader(getUserRequest)
	getUserRequest.Header.Add("Authorization", fmt.Sprintf("token %s", tokenResponse.AccessToken))
	userResp := clientDo(getUserRequest)
	bytes, err := ioutil.ReadAll(userResp.Body)
	if err != nil {
		return nil, err
	}
	respJSON := string(bytes)
	var githubUserResponse data.GithubUserResponse
	if err = util.FromJSONString(respJSON, &githubUserResponse); err != nil {
		return nil, err
	}
	return &githubUserResponse, nil
}
