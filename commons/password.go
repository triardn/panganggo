package commons

import (
	"log"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password []byte) string {
	// Use GenerateFromPassword to hash & salt password
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPassword string, plainPassword []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func ValidatePassword(s string) bool {
	var (
		hasMinMaxLen = false
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
	)

	if len(s) >= 6 && len(s) <= 64 {
		hasMinMaxLen = true
	}

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinMaxLen && hasUpper && hasLower && hasNumber && hasSpecial
}
