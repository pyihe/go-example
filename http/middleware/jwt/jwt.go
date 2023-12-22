package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	httpnet "github.com/pyihe/go-example/http"
)

const (
	keyJwtJTI = "jti"
)

// WithJWT JWT中间件，访问API时做Token校验
func WithJWT(method jwt.SigningMethod, publicKey interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var token *jwt.Token
		var header = c.Request.Header
		var tokenStr = header.Get(httpnet.Authorization)

		if tokenArray := strings.Split(tokenStr, " "); len(tokenArray) == 2 {
			tokenStr = tokenArray[1]
		}
		token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		}, jwt.WithValidMethods([]string{method.Alg()}))

		if err != nil || token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			c.Set(httpnet.TokenKey, claims[keyJwtJTI])
		}
		c.Next()
	}
}

type Request struct {
	ClientID   string            // 客户端ID
	PrivateKey interface{}       // 对应签名类型的私钥
	Expire     time.Duration     // Token有效期
	Method     jwt.SigningMethod // WithJWT 签名类型
}

// Generate 生成JSON WEB TOKEN
// 生成普通的Token
func Generate(apply Request) (string, error) {
	var now = jwt.TimeFunc()
	var claims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(apply.Expire)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        apply.ClientID,
	}
	var token = jwt.NewWithClaims(apply.Method, claims)
	return token.SignedString(apply.PrivateKey)
}
