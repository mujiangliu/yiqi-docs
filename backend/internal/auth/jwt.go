// backend/internal/auth/jwt.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷
type Claims struct {
	UID  uint   `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// JWTIssuer 签发与解析 JWT
type JWTIssuer struct {
	secret []byte
}

// NewJWTIssuer 创建一个签发器
func NewJWTIssuer(secret string) *JWTIssuer {
	return &JWTIssuer{secret: []byte(secret)}
}

// Issue 签发一个带过期时间的 token
func (j *JWTIssuer) Issue(uid uint, role string, ttl time.Duration) (string, error) {
	claims := Claims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secret)
}

// Parse 解析并校验 token
func (j *JWTIssuer) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
