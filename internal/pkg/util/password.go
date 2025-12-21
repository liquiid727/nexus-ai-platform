package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash 检查密码是否与哈希值匹配
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GeneratePasswordResetToken 生成密码重置令牌
func GeneratePasswordResetToken() string {
	return GenerateRandomString(32)
}

// GenerateRandomString 生成X位随机字符串
func GenerateRandomString(x int) string {
	panic("unimplemented")
}
