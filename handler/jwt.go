package handler

// jwt 包 当用户通过github授权后 生成token用
import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lanwupark/blog-api/data"
)

var (
	secretKey = []byte("eansonG&%34R&R") //密钥 写死了的
)

// CreateToken 创建token
func CreateToken(user *data.User) (tokenString string, err error) {
	var subject string
	subject, err = data.ToJSONString(user)
	if err != nil {
		return "", err
	}
	claims := &jwt.StandardClaims{
		Issuer:    "eanson",                                  //签发人 我写我自己^_^
		NotBefore: time.Now().Unix(),                         //生效时间
		ExpiresAt: int64(time.Now().Add(time.Minute).Unix()), // 过期时间
		Subject:   subject,                                   //主体 json数据
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	return
}

// ParseToken 解析json数据
func ParseToken(tokenString string) (user *data.User, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// secretKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err = data.FromJSONString(claims["sub"].(string), &user)
		if err != nil {
			return nil, err
		}
		return
	}
	return nil, err
}
