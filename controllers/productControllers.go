package controllers

import (
	"meeting3/entity"
	"meeting3/tools"

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

	return c.JSON(fiber.Map{
		"pesan": Product,
	})
}

func GetProduct(c *fiber.Ctx) error {

	const baseUrl = `http://localhost:3000`

	produc := entity.Product{
		Nama:   "Susi Tuna",
		Harga:  "30000",
		Gambar: "2.png",
	}

	produc.Gambar = baseUrl + "/public/product/" + produc.Gambar
	return c.JSON(fiber.Map{
		"pesan": "berhasil get data product",
		"data":  produc,
	})
}
