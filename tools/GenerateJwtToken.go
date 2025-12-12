package tools

import (
	"meeting3/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(entitys interface{}) string {
	var name, role, types string
	var status bool
	var idAdmin float64 // Mengubah tipe data menjadi string

	switch e := entitys.(type) {
	case entity.Admin:
		name = e.Nama
		idAdmin = float64((int(e.ID))) // Contoh penggunaan ID, sesuaikan dengan kebutuhan Anda
		role = "1"
		types = "admin"
		status = e.Status
		name = e.Nama
	default:
		return ""
	}

	claims := jwt.MapClaims{
		"name":     name,
		"id_admin": idAdmin,
		"status":   status,
		"role":     role,
		"type":     types,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte("Elaut@3123!"))
	return t
}
