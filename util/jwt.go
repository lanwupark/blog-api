package util

// jwt 包 当用户通过github授权后 生成token用
import (
	"fmt"
	"time"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lanwupark/blog-api/data"
)

var (
	secretKey   = []byte("eansonG&%34R&R") //密钥 写死了的
	expiredTime = time.Minute * 30         //过期时间 30分钟
)

// CreateToken 创建token
func CreateToken(subject *data.TokenClaimsSubject) (tokenString string, err error) {
	var sub string
	sub, err = ToJSONString(subject)
	if err != nil {
		return "", err
	}
	claims := &jwt.StandardClaims{
		Issuer:    "eanson",                                  //签发人 我写我自己^_^
		NotBefore: time.Now().Unix(),                         //生效时间
		ExpiresAt: int64(time.Now().Add(expiredTime).Unix()), // 过期时间
		Subject:   sub,                                       //主体 json数据
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	return
}

// ParseToken 解析json数据
func ParseToken(tokenString string) (subject *data.TokenClaimsSubject, err error) {
	mapClaims, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	err = FromJSONString(mapClaims["sub"].(string), &subject)
	if err != nil {
		return nil, err
	}
	return
}

// parseToken 获取负荷
func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// secretKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// RefreshToken 刷新token 如果过期时间小于expiredTime*0.2 那么就刷新Token
func RefreshToken(tokenString string) (newToken string, success bool) {
	success = false
	claims, err := parseToken(tokenString)
	if err != nil {
		return
	}
	// 当前时间戳
	now := time.Now().Unix()
	created, expired := int64(claims["nbf"].(float64)), int64(claims["exp"].(float64))
	// 小小计算
	lastSeconds := expired - now
	// 剩余过期时间 in [0,20%] 重新刷新Token
	if lastSeconds > 0 && lastSeconds <= (expired-created)*2/10 {
		var subject data.TokenClaimsSubject
		err = FromJSONString(claims["sub"].(string), &subject)
		if err != nil {
			return
		}
		newToken, err = CreateToken(&subject)
		if err != nil {
			return
		}
		log.Infof("refresh token %+v\t token:\n", subject.UserID, newToken)
		success = true
	}
	return
}
