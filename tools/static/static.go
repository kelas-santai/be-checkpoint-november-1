package static

import "github.com/gofiber/fiber/v2"

func StaticProduct(c *fiber.Ctx) error {

	path := c.Params("name")
	//2.png
	//send berdasarkan nama file yang masuk
	return c.SendFile("./public/product/" + path)

}
