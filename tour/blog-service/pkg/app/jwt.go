package app

import (
	"blog-service/global"
	"blog-service/pkg/util"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims JWT 载荷结构体
type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

// GetJWTSecret 获取 JWT 密钥
func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

// GenerateToken 生成 JWT Token
func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	exprieTime := nowTime.Add(global.JWTSetting.Expire)
	// 创建 JWT 载荷
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exprieTime.Unix(),
			IssuedAt:  global.DBEngine.RowsAffected,
		},
	}

	// 根据 Claims 结构体创建 Token 实例
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// SignedString 方法使用指定的密钥签名并获得完整的编码后的字符串
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

// 解析和校验Token
func ParseToken(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回 *Token。
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	// Vaild 验证基于时间的声明
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
