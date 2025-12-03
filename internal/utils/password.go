package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return ""
	}

	return string(hashedPass)
}

func CheckPass(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))

	return err == nil
}
