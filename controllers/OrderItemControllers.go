package controllers

import (
	"meeting3/database"
	"meeting3/entity"

	"github.com/gofiber/fiber/v2"
)

func GetOrderItem(c *fiber.Ctx) error {

	id := c.Query("id")
	var orderItem entity.OrderItem
	database.DB.Preload("Product").Preload("Order").Where("id = ?", id).First(&orderItem)

	if orderItem.IdOrderItem == 0 {
		return c.Status(404).JSON(fiber.Map{
			"pesan": "Order item tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil get data order item",
		"data":  orderItem,
	})
}

func GetOrderItemsByOrderID(c *fiber.Ctx) error {

	orderID := c.Query("order_id")
	var orderItems []entity.OrderItem
	database.DB.Preload("Product").Where("id_order = ?", orderID).Find(&orderItems)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua order item",
		"data":  orderItems,
	})
}

func UpdateOrderItem(c *fiber.Ctx) error {

	id := c.QueryInt("id")

	var orderItem entity.OrderItem
	database.DB.Where("id = ?", id).First(&orderItem)

	if orderItem.IdOrderItem == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Order Item di database untuk ID itu",
		})
	}

	//Cek status order
	var order entity.Order
	database.DB.Where("id = ?", orderItem.IdOrder).First(&order)

	if order.Status == "completed" || order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak bisa update order item karena order sudah completed atau cancelled",
		})
	}

	type UpdateRequest struct {
		Quantity int    `json:"quantity"`
		Notes    string `json:"notes"`
	}

	var request UpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	//Validasi product stock
	var product entity.Product
	database.DB.Where("id = ?", orderItem.ProductID).First(&product)

	//Hitung selisih quantity
	qtyDiff := request.Quantity - orderItem.Quantity

	if qtyDiff > 0 {
		//Jika quantity bertambah, cek stock
		if product.Stock < qtyDiff {
			return c.Status(400).JSON(fiber.Map{
				"pesan": "Stock product tidak mencukupi",
			})
		}
		product.Stock -= qtyDiff
	} else if qtyDiff < 0 {
		//Jika quantity berkurang, kembalikan stock
		product.Stock += -qtyDiff
	}

	database.DB.Save(&product)

	//Update order item
	orderItem.Quantity = request.Quantity
	orderItem.Subtotal = float64(request.Quantity) * orderItem.Price
	orderItem.Notes = request.Notes

	database.DB.Save(&orderItem)

	//Update total amount di order
	var allOrderItems []entity.OrderItem
	database.DB.Where("id_order = ?", orderItem.IdOrder).Find(&allOrderItems)

	var newTotal float64
	for _, item := range allOrderItems {
		newTotal += item.Subtotal
	}

	order.TotalAmount = newTotal
	database.DB.Save(&order)

	//Load relations untuk response
	database.DB.Preload("Product").Preload("Order").Where("id = ?", orderItem.IdOrderItem).First(&orderItem)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update order item",
		"data":  orderItem,
	})
}

func DeleteOrderItem(c *fiber.Ctx) error {

	id := c.Query("id")
	var orderItem entity.OrderItem
	database.DB.Where("id = ?", id).First(&orderItem)

	if orderItem.IdOrderItem == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Order Item di database",
		})
	}

	//Cek status order
	var order entity.Order
	database.DB.Where("id = ?", orderItem.IdOrder).First(&order)

	if order.Status == "completed" || order.Status == "cancelled" {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak bisa hapus order item karena order sudah completed atau cancelled",
		})
	}

	//Kembalikan stock
	var product entity.Product
	database.DB.Where("id = ?", orderItem.ProductID).First(&product)
	product.Stock += orderItem.Quantity
	database.DB.Save(&product)

	//Hapus order item
	database.DB.Delete(&orderItem)

	//Update total amount di order
	var allOrderItems []entity.OrderItem
	database.DB.Where("id_order = ?", orderItem.IdOrder).Find(&allOrderItems)

	var newTotal float64
	for _, item := range allOrderItems {
		newTotal += item.Subtotal
	}

	order.TotalAmount = newTotal
	database.DB.Save(&order)

	return c.JSON(fiber.Map{
		"pesan": "Order item telah dihapus",
	})
}
