package controllers

import (
	"meeting3/database"
	"meeting3/entity"
	"meeting3/tools"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {

	var Product entity.Product
	if err := c.BodyParser(&Product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//buat untuk penyimpanan gambar
	gambar, _ := c.FormFile("gambar")

	//untuk menghapus spaci dari nama file gambar
	fileName := tools.RemoveSpaci(gambar.Filename)

	//Simpan Data ke Folder

	if err := c.SaveFile(gambar, "./public/product/"+fileName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Gagal Save File Gambar",
			"err":   err,
		})
	}

	Product.Gambar = fileName
	//Menyimpan Ke database
	database.DB.Create(&Product)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat product",
		"data":  Product,
	})
}

func GetProduct(c *fiber.Ctx) error {

	const baseUrl = `http://localhost:3000`

	id := c.Query("id")
	var product entity.Product
	database.DB.Where("id = ?", id).First(&product)

	if product.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Product tidak ditemukan",
		})
	}

	product.Gambar = baseUrl + "/public/product/" + product.Gambar
	return c.JSON(fiber.Map{
		"pesan": "berhasil get data product",
		"data":  product,
	})
}

func GetAllProduct(c *fiber.Ctx) error {

	const baseUrl = `http://localhost:3000`

	var products []entity.Product
	database.DB.Find(&products)

	//Update gambar URL untuk setiap product
	for i := range products {
		products[i].Gambar = baseUrl + "/public/product/" + products[i].Gambar
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua product",
		"data":  products,
	})
}

func UpdateProduct(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var product entity.Product
	database.DB.Where("id = ?", id).First(&product)

	if product.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Product di database untuk ID itu",
		})
	}

	var request entity.Product
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Update data product
	product.Nama = request.Nama
	product.Harga = request.Harga
	product.IdCategory = request.IdCategory
	product.Description = request.Description
	product.Stock = request.Stock
	product.IsAvailable = request.IsAvailable

	//Cek apakah ada gambar baru
	gambar, err := c.FormFile("gambar")
	if err == nil {
		//Hapus gambar lama
		oldImagePath := "./public/product/" + product.Gambar
		os.Remove(oldImagePath)

		//Simpan gambar baru
		fileName := tools.RemoveSpaci(gambar.Filename)
		if err := c.SaveFile(gambar, "./public/product/"+fileName); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"pesan": "Gagal Save File Gambar",
				"err":   err,
			})
		}
		product.Gambar = fileName
	}

	//Simpan ke database
	database.DB.Save(&product)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update product",
		"data":  product,
	})
}

func DeleteProduct(c *fiber.Ctx) error {

	id := c.Query("id")
	var product entity.Product
	database.DB.Where("id = ?", id).First(&product)

	if product.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Product di database",
		})
	}

	//Hapus file gambar
	imagePath := "./public/product/" + product.Gambar
	os.Remove(imagePath)

	//Hapus data dari database
	database.DB.Delete(&product)

	return c.JSON(fiber.Map{
		"pesan": "Product telah dihapus",
	})
}
