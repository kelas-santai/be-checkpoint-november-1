package main

import (
	"log"
	"meeting3/controllers"
	"meeting3/database"
	"meeting3/tools/static"

	"github.com/gofiber/fiber/v2"
)

func Hello(c *fiber.Ctx) error {

	type User struct {
		Id     uint
		Nama   string
		Alamat string
	}

	users := []User{
		{
			Id:     1,
			Nama:   "Bagja",
			Alamat: "Tangerang",
		},
		{
			Id:     2,
			Nama:   "Lazwardi",
			Alamat: "Tangerang",
		},
	}

	return c.JSON(fiber.Map{
		"data":  users,
		"pesan": "berhasil mendapatkan data dari database",
	})
}

func main() {
	app := fiber.New()

	database.Connect()

	product := app.Group("product")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/data", Hello)
	//----------------users --------------
	app.Post("/create-user-cara-1", controllers.CreateUserCara1)
	app.Post("/create-user-cara-2", controllers.CreateUserCara2)
	app.Get("/get-data-user-by-id", controllers.GetUserByParameter)
	app.Get("/all-users", controllers.GetAllUser)
	app.Put("/update-user", controllers.UpdateUsers)
	app.Delete("/delete-user", controllers.DeleteUsers)
	app.Get("/get-raw-query", controllers.GetRawQuery)
	//app.Delete("/delete-user")

	//--------------------admin --------------------
	admin := app.Group("/admin")
	admin.Post("/create-admin", controllers.CreateAdmin)
	admin.Get("/get-admin", controllers.GetAdmin)
	admin.Get("/all-admin", controllers.GetAllAdmin)
	admin.Put("/update-admin", controllers.UpdateAdmin)
	admin.Delete("/delete-admin", controllers.DeleteAdmin)

	//auth Admin
	//app.Post("/admin/login", controllers.AdminLogin)

	//--------------------product--------------
	product.Post("/create-product", controllers.CreateProduct)
	product.Get("/get-product", controllers.GetProduct)
	product.Get("/all-product", controllers.GetAllProduct)
	product.Put("/update-product", controllers.UpdateProduct)
	product.Delete("/delete-product", controllers.DeleteProduct)

	//--------------------category--------------
	category := app.Group("/category")
	category.Post("/create-category", controllers.CreateCategory)
	category.Get("/get-category", controllers.GetCategory)
	category.Get("/all-category", controllers.GetAllCategory)
	category.Put("/update-category", controllers.UpdateCategory)
	category.Delete("/delete-category", controllers.DeleteCategory)

	//--------------------table--------------
	table := app.Group("/table")
	table.Post("/create-table", controllers.CreateTable)
	table.Get("/get-table", controllers.GetTable)
	table.Get("/token/:token", controllers.GetTableByToken)
	table.Get("/all-table", controllers.GetAllTable)
	table.Put("/update-table", controllers.UpdateTable)
	table.Put("/regenerate-qr", controllers.RegenerateQRToken)
	table.Delete("/delete-table", controllers.DeleteTable)

	//--------------------order--------------
	order := app.Group("/order")
	order.Post("/create-order", controllers.CreateOrder)
	order.Get("/get-order", controllers.GetOrder)
	order.Get("/tracking/:token", controllers.GetOrderByTracking)
	order.Get("/all-order", controllers.GetAllOrder)
	order.Put("/update-status", controllers.UpdateOrderStatus)
	order.Put("/update-payment", controllers.UpdatePaymentStatus)
	order.Delete("/delete-order", controllers.DeleteOrder)

	//--------------------order-item--------------
	orderItem := app.Group("/order-item")
	orderItem.Get("/get-order-item", controllers.GetOrderItem)
	orderItem.Get("/by-order", controllers.GetOrderItemsByOrderID)
	orderItem.Put("/update-order-item", controllers.UpdateOrderItem)
	orderItem.Delete("/delete-order-item", controllers.DeleteOrderItem)

	//--------------------static file------------------

	app.Get("/public/product/:name", static.StaticProduct)

	log.Fatal(app.Listen(":3000"))
}
