package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func Password(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func EqualPassword(hashedPass, newPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(newPass))
	return err == nil
}
