package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("生成密码哈希失败:%s", err)
		return "", err
	}
	return string(passwordHash), nil
}

func CompareHashAndPassword(passwordHash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
