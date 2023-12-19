package utils

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Error().Msg("error hashing password")
		return ""
	}

	return string(bytes)
}

func PasswordMatchesHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
