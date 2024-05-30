package lib

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(i []byte) string {
	hash, err := bcrypt.GenerateFromPassword(i, -1)

	if err != nil {
		log.Fatal("SECURITY ERROR / hashing not working.")
	}

	return string(hash)
}

func ValidateHash(h []byte, i []byte) bool {
	return bcrypt.CompareHashAndPassword(h, i) == nil
}
