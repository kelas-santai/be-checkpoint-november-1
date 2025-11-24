package controllers

import (
	"meeting3/database"
	"meeting3/entity"

	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {

	var category entity.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Cek apakah nama category sudah ada
	var existing entity.Category
	database.DB.Where("nama = ?", category.Nama).First(&existing)
	if existing.IdCategory != 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Kategori dengan nama tersebut sudah ada",
		})
	}

	//Simpan ke database
	database.DB.Create(&category)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat category",
		"data":  category,
	})
}

func GetCategory(c *fiber.Ctx) error {

	id := c.Query("id")
	var category entity.Category
	database.DB.Preload("Product").Where("id = ?", id).First(&category)

	if category.IdCategory == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Category tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data category",
		"data":  category,
	})
}

func GetAllCategory(c *fiber.Ctx) error {

	var categories []entity.Category
	database.DB.Preload("Product").Find(&categories)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua category",
		"data":  categories,
	})
}

func UpdateCategory(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var category entity.Category
	database.DB.Where("id = ?", id).First(&category)

	if category.IdCategory == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Category di database untuk ID itu",
		})
	}

	var request entity.Category
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Update data category
	category.Nama = request.Nama
	category.IsActive = request.IsActive

	//Simpan ke database
	database.DB.Save(&category)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update category",
		"data":  category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {

	id := c.Query("id")
	var category entity.Category
	database.DB.Where("id = ?", id).First(&category)

	if category.IdCategory == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Category di database",
		})
	}

	//Cek apakah ada product yang menggunakan category ini
	var productCount int64
	database.DB.Model(&entity.Product{}).Where("id_category = ?", id).Count(&productCount)

	if productCount > 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak bisa menghapus category, masih ada product yang menggunakan category ini",
		})
	}

	//Hapus data dari database
	database.DB.Delete(&category)

	return c.JSON(fiber.Map{
		"pesan": "Category telah dihapus",
	})
}
