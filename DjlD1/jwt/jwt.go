package jwt

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义声明结构体
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成 JWT
func GenerateJWT(username string) (string, error) {
	jwtKey, err := os.ReadFile("key/Private.key")
	if err != nil {
		log.Fatalf("读取私钥失败: %v", err)
	}
	// 设置声明内容
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 2小时后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                    //立即生效  time.Now().Add(5 * time.Minute) 表示五分钟后才有效
			Issuer:    "djluser",                                         //签发者
			Subject:   "登录验证",                                            //主题
		},
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(jwtKey) //创建私钥对象
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}
	// 创建 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 使用密钥签名并生成字符串
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func Verifytoken(token string) (string, error) {
	jwtKey, err := os.ReadFile("key/Public.key")
	if err != nil {
		log.Fatalf("读取公钥失败: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(jwtKey) //创建公钥对象
	if err != nil {
		log.Fatalf("解析公钥钥失败: %v", err)
	}
	token1, err := jwt.ParseWithClaims(token, &Claims{}, func(token1 *jwt.Token) (interface{}, error) {

		if _, ok := token1.Method.(*jwt.SigningMethodRSA); !ok { // 确保算法一致
			return nil, jwt.ErrTokenUnverifiable
		}
		return publicKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token1.Claims.(*Claims); ok && token1.Valid {
		return claims.Username, nil
	}
	return "", jwt.ErrTokenInvalidClaims
}
