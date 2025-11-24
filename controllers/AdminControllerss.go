package controllers

import (
	"meeting3/database"
	"meeting3/entity"
	"meeting3/tools"

	"github.com/gofiber/fiber/v2"
)

func CreateAdmin(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "gagal parser body",
			"err":   err.Error(),
		})
	}

	newData := entity.Admin{
		Nama:     data["nama"],
		Email:    data["email"],
		Password: tools.GeneratePassword(data["password"]),
		Alamat:   data["alamat"],
		NoTelpon: data["no_telpon"],
	}

	database.DB.Create(&newData)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat admin",
	})
}

func GetAdmin(c *fiber.Ctx) error {

	id := c.Query("id")

	var admin entity.Admin
	database.DB.Where("id = ?", id).First(&admin)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan data admin",
		"data":  admin,
	})
}
