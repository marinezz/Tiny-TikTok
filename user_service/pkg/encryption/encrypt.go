// 密码加密

package encryption

import "golang.org/x/crypto/bcrypt"

const PassWordCost = 12 // 加密难度，相当与hash被迭代2^12次。默认是10 ，推荐使用12

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

// VerifyPasswordWithHash 方法验证传入的密码与数据库中的密码哈希是否匹配
func VerifyPasswordWithHash(inputPassword string, storedHash string) bool {
	// 将数据库中的密码哈希字符串解析为哈希字节
	storedHashBytes := []byte(storedHash)

	// 使用 bcrypt 提供的 CompareHashAndPassword 方法比较密码
	err := bcrypt.CompareHashAndPassword(storedHashBytes, []byte(inputPassword))
	return err == nil
}
