package services

import (
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateId()(string){

	guid := xid.New()
	id := guid.String()

	return id
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

