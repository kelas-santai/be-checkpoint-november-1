package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"meeting3/database"
	"meeting3/entity"
	"time"

	"github.com/gofiber/fiber/v2"
)

func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102")
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("ORD-%s-%s", timestamp, hex.EncodeToString(b))
}

func generateTrackingToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func CreateOrder(c *fiber.Ctx) error {

	type OrderRequest struct {
		TableNumber   uint               `json:"table_number"`
		CustomerName  string             `json:"customer_name"`
		CustomerEmail string             `json:"customer_email"`
		PaymentMethod string             `json:"payment_method"`
		Notes         string             `json:"notes"`
		OrderItems    []entity.OrderItem `json:"order_items"`
	}

	var request OrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Validasi table
	var table entity.Table
	database.DB.Where("number = ?", request.TableNumber).First(&table)
	if table.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Table tidak ditemukan",
		})
	}

	if !table.IsActive {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Table tidak aktif",
		})
	}

	//Hitung total amount
	var totalAmount float64
	for i := range request.OrderItems {
		//Validasi product
		var product entity.Product
		database.DB.Where("id = ?", request.OrderItems[i].ProductID).First(&product)
		if product.ID == 0 {
			return c.Status(400).JSON(fiber.Map{
				"pesan": fmt.Sprintf("Product dengan ID %d tidak ditemukan", request.OrderItems[i].ProductID),
			})
		}

		//Cek stock
		if product.Stock < request.OrderItems[i].Quantity {
			return c.Status(400).JSON(fiber.Map{
				"pesan": fmt.Sprintf("Stock product %s tidak mencukupi", product.Nama),
			})
		}

		//Hitung subtotal
		request.OrderItems[i].Subtotal = float64(request.OrderItems[i].Quantity) * request.OrderItems[i].Price
		totalAmount += request.OrderItems[i].Subtotal
	}

	//Buat order
	order := entity.Order{
		OrderNumber:   generateOrderNumber(),
		Table:         string(table.Number),
		CustomerName:  request.CustomerName,
		CustomerEmail: request.CustomerEmail,
		Status:        "pending",
		PaymentMethod: request.PaymentMethod,
		PaymentStatus: "unpaid",
		TotalAmount:   totalAmount,
		Notes:         request.Notes,
		TrackingToken: generateTrackingToken(),
	}

	//Simpan order
	database.DB.Create(&order)

	//Simpan order items
	for i := range request.OrderItems {
		request.OrderItems[i].IdOrder = order.IdOrder
		database.DB.Create(&request.OrderItems[i])

		//Update stock product
		var product entity.Product
		database.DB.Where("id = ?", request.OrderItems[i].ProductID).First(&product)
		product.Stock -= request.OrderItems[i].Quantity
		database.DB.Save(&product)
	}

	//Load order items untuk response
	database.DB.Preload("OrderItems.Product").Preload("Table").Where("id = ?", order.IdOrder).First(&order)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat order",
		"data":  order,
	})
}

func GetOrder(c *fiber.Ctx) error {

	id := c.Query("id")
	var order entity.Order
	database.DB.Preload("OrderItems").Where("id = ?", id).First(&order)

	if order.IdOrder == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Order tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data order",
		"data":  order,
	})
}

func GetOrderByTracking(c *fiber.Ctx) error {

	token := c.Params("token")
	var order entity.Order
	database.DB.Preload("OrderItems.Product").Preload("Table").Where("tracking_token = ?", token).First(&order)

	if order.IdOrder == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Order tIdOrderak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data order",
		"data":  order,
	})
}

func GetAllOrder(c *fiber.Ctx) error {

	var orders []entity.Order
	database.DB.Preload("OrderItems").Find(&orders)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua order",
		"data":  orders,
	})
}

func UpdateOrderStatus(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var order entity.Order
	database.DB.Where("id = ?", id).First(&order)

	if order.IdOrder == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Order di database untuk ID itu",
		})
	}

	type StatusRequest struct {
		Status string `json:"status"`
	}

	var request StatusRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Validasi status
	validStatuses := []string{"pending", "processing", "completed", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if request.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Status tidak valid. Harus pending, processing, completed, atau cancelled",
		})
	}

	//Update status
	order.Status = request.Status

	//Jika order dibatalkan, kembalikan stock
	if request.Status == "cancelled" && order.Status != "cancelled" {
		var orderItems []entity.OrderItem
		database.DB.Where("id_order = ?", order.IdOrder).Find(&orderItems)

		for _, item := range orderItems {
			var product entity.Product
			database.DB.Where("id = ?", item.ProductID).First(&product)
			product.Stock += item.Quantity
			database.DB.Save(&product)
		}
	}

	database.DB.Save(&order)

	//Load order items untuk response
	database.DB.Preload("OrderItems.Product").Preload("Table").Where("id = ?", order.IdOrder).First(&order)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update status order",
		"data":  order,
	})
}

func UpdatePaymentStatus(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var order entity.Order
	database.DB.Where("id = ?", id).First(&order)

	if order.IdOrder == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Order di database untuk ID itu",
		})
	}

	type PaymentRequest struct {
		PaymentStatus string `json:"payment_status"`
	}

	var request PaymentRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Validasi payment status
	if request.PaymentStatus != "paid" && request.PaymentStatus != "unpaid" {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Payment status tidak valid. Harus paid atau unpaid",
		})
	}

	order.PaymentStatus = request.PaymentStatus
	database.DB.Save(&order)

	//Load order items untuk response
	database.DB.Preload("OrderItems.Product").Preload("Table").Where("id = ?", order.IdOrder).First(&order)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update payment status order",
		"data":  order,
	})
}

func DeleteOrder(c *fiber.Ctx) error {

	id := c.Query("id")
	var order entity.Order
	database.DB.Where("id = ?", id).First(&order)

	if order.IdOrder == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Order di database",
		})
	}

	//Kembalikan stock jika order belum completed
	if order.Status != "completed" && order.Status != "cancelled" {
		var orderItems []entity.OrderItem
		database.DB.Where("id_order = ?", order.IdOrder).Find(&orderItems)

		for _, item := range orderItems {
			var product entity.Product
			database.DB.Where("id = ?", item.ProductID).First(&product)
			product.Stock += item.Quantity
			database.DB.Save(&product)
		}
	}

	//Hapus order items terlebih dahulu
	database.DB.Where("id_order = ?", order.IdOrder).Delete(&entity.OrderItem{})

	//Hapus order
	database.DB.Delete(&order)

	return c.JSON(fiber.Map{
		"pesan": "Order telah dihapus",
	})
}
