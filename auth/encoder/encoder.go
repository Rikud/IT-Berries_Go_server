package encoder

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), err
}

func ComparePasswords(dbHash string, password string) bool {
	plainHash := []byte(password)
	hashedPassword := []byte(dbHash)
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainHash)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
