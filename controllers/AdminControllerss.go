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

	//Cek apakah email atau no_telpon sudah ada
	var existing entity.Admin
	database.DB.Where("email = ? OR no_telpon = ?", data["email"], data["no_telpon"]).First(&existing)
	if existing.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Email atau no telpon sudah terdaftar",
		})
	}

	newData := entity.Admin{
		Nama:     data["nama"],
		Email:    data["email"],
		Role:     data["role"],
		Password: tools.GeneratePassword(data["password"]),
		Alamat:   data["alamat"],
		NoTelpon: data["no_telpon"],
	}

	database.DB.Create(&newData)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat admin",
		"data":  newData,
	})
}

func GetAdmin(c *fiber.Ctx) error {

	id := c.Query("id")

	var admin entity.Admin
	database.DB.Where("id = ?", id).First(&admin)

	if admin.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Admin tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan data admin",
		"data":  admin,
	})
}

func GetAllAdmin(c *fiber.Ctx) error {

	var allAdmins []entity.Admin
	database.DB.Find(&allAdmins)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua admin",
		"data":  allAdmins,
	})
}

func UpdateAdmin(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var admin entity.Admin
	database.DB.Where("id = ?", id).First(&admin)

	if admin.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Admin di database untuk ID itu",
		})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Update data admin
	admin.Nama = data["nama"]
	admin.Alamat = data["alamat"]
	admin.Email = data["email"]
	admin.NoTelpon = data["no_telpon"]
	admin.Role = data["role"]

	//Update password jika ada
	if data["password"] != "" {
		admin.Password = tools.GeneratePassword(data["password"])
	}

	//Simpan ke database
	database.DB.Save(&admin)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update admin",
		"data":  admin,
	})
}

func DeleteAdmin(c *fiber.Ctx) error {

	id := c.Query("id")
	var admin entity.Admin
	database.DB.Where("id = ?", id).First(&admin)

	if admin.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Admin di database",
		})
	}

	//Hapus data dari database
	database.DB.Delete(&admin)

	return c.JSON(fiber.Map{
		"pesan": "Admin telah dihapus",
	})
}
