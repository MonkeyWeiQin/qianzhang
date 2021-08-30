package middleware

import (
	"battle_rabbit/define"
	"battle_rabbit/server/login/proto"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var (
	JwtObj           *JWT
	JwtSecret        = ""
	TokenExpired     = errors.New("授权已过期")
	TokenNotValidYet = errors.New("授权未激活")
	TokenMalformed   = errors.New("认证授权错误")
	TokenInvalid     = errors.New("无效的认证")
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(JwtSecret),
	}
}
func InitJwtAuth(secret string) {
	JwtSecret = secret
	JwtObj = NewJWT()
}

/*
	json web token 认证中间件
*/
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*		c.Next()
				return*/
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			proto.ResponseFailWithCode(c, define.MsgCode401)
			c.Abort()
			return
		}
		//if service.IsBlacklist(token) {
		//	response.Send(response.ERROR, gin.H{
		//		"reload": true,
		//	}, "您的帐户异地登陆或令牌失效", c)
		//	c.Abort()
		//	return
		//}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			c.Abort()
			proto.ResponseFailWithCode(c, define.MsgCode401)
			return
		}
		// 判断是否要还有半个小时就过期了,是的话就刷新token
		if claims.ExpiresAt-time.Now().Unix() < 1800 {
			claims.NotBefore = time.Now().Unix() - 10
			claims.ExpiresAt = time.Now().Add(time.Hour * 3).Unix()
			newToken, err := j.CreateToken(*claims)
			if err != nil {
				c.Abort()
				proto.ResponseFailWithCode(c, define.MsgCode401)
				return
			}
			c.Header("x-token", newToken)
			c.Header("expires-at", strconv.FormatInt(claims.ExpiresAt, 10))
		}
		c.Set("claims", claims)
		userId, err := strconv.Atoi(claims.Id)
		if err != nil || userId == 0 {
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims jwt.StandardClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
