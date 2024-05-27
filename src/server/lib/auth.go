package lib

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
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

func ReadPassword() []byte {
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Password can't read.")
	}

	return password
}
