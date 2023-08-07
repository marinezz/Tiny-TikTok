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
