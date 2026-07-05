// backend/internal/auth/password.go
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword 用 bcrypt 哈希密码
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword 校验密码与哈希是否匹配
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
