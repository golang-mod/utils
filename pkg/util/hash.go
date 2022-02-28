package util

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword 密码加密
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword 密码比对
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}
