package tools

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(pass string) string {

	hasilGenerate, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	hasilString := string(hasilGenerate)

	return hasilString
}

func ComparePassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
