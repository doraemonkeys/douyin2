package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 对传入字符串进行加密
func BcryptHash(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// BcryptMatch 对传入的加密字符串进行比对,str为明文
func BcryptMatch(hash string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
