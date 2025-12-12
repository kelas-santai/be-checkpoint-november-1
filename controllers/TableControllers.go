package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"meeting3/database"
	"meeting3/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Tambahab
// Ok
func CreateTable(c *fiber.Ctx) error {

	var table map[string]string
	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	numberMejaInt, _ := strconv.Atoi(table["number_table"])

	//Cek apakah nomor meja sudah ada
	var existing entity.Table
	database.DB.Where("number = ?", numberMejaInt).First(&existing)
	if existing.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Nomor meja sudah terdaftar",
		})
	}

	//Generate QR Token
	QRToken := generateToken()

	//Generate QR Code URL

	//link aplikasi frontend yang akan diakses ketika QR code di scan
	baseUrl := "http://localhost:3000"

	urlCode := baseUrl + "/order/table/" + QRToken

	fileNameQrCode := fmt.Sprintf("table_%d_qr.png", numberMejaInt)

	filePath := "./public/qrcode/" + fileNameQrCode
	// Buat QR code
	err := qrcode.WriteFile(urlCode, qrcode.Medium, 256, filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"pesan": "gagal membuat QR code",
			"err":   err.Error(),
		})
	}
	QRCode := fmt.Sprintf("%s/order/table/%s", baseUrl, QRToken)
	//Simpan ke database

	newData := entity.Table{
		Number:   numberMejaInt,
		QRCode:   QRCode,
		QRToken:  QRToken,
		IsActive: true,
	}
	database.DB.Create(&newData)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat table",
		"data":  table,
	})
}

func GetTable(c *fiber.Ctx) error {

	id := c.Query("id")
	var table entity.Table
	database.DB.Where("id = ?", id).First(&table)

	if table.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Table tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data table",
		"data":  table,
	})
}

func GetTableByToken(c *fiber.Ctx) error {

	token := c.Params("token")
	var table entity.Table
	database.DB.Where("qr_token = ?", token).First(&table)

	if table.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Table tidak ditemukan",
		})
	}

	if !table.IsActive {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Table tidak aktif",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data table",
		"data":  table,
	})
}

func GetAllTable(c *fiber.Ctx) error {

	var tables []entity.Table
	database.DB.Find(&tables)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua table",
		"data":  tables,
	})
}

func UpdateTable(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var table entity.Table
	database.DB.Where("id = ?", id).First(&table)

	if table.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Table di database untuk ID itu",
		})
	}

	var request entity.Table
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Update data table
	table.Number = request.Number
	table.IsActive = request.IsActive

	//Simpan ke database
	database.DB.Save(&table)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update table",
		"data":  table,
	})
}

func DeleteTable(c *fiber.Ctx) error {

	id := c.Query("id")
	var table entity.Table
	database.DB.Where("id = ?", id).First(&table)

	if table.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Table di database",
		})
	}

	//Hapus data dari database
	database.DB.Delete(&table)

	return c.JSON(fiber.Map{
		"pesan": "Table telah dihapus",
	})
}

func RegenerateQRToken(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var table entity.Table
	database.DB.Where("id = ?", id).First(&table)

	if table.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Table di database untuk ID itu",
		})
	}

	//Generate QR Token baru
	table.QRToken = generateToken()

	//Generate QR Code URL baru
	baseUrl := "http://localhost:3000"
	table.QRCode = fmt.Sprintf("%s/order/table/%s", baseUrl, table.QRToken)

	//Simpan ke database
	database.DB.Save(&table)

	return c.JSON(fiber.Map{
		"pesan": "berhasil regenerate QR token",
		"data":  table,
	})
}
