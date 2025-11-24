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

	//--------------------product--------------
	product.Post("/create-product", controllers.CreateProduct)
	product.Get("/get-product", controllers.GetProduct)

	//--------------------static file------------------

	app.Get("/public/product/:name", static.StaticProduct)

	log.Fatal(app.Listen(":3000"))
}
