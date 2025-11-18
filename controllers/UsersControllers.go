package controllers

import (
	"fmt"
	"meeting3/entity"

	"github.com/gofiber/fiber/v2"
)

func CreateUserCara1(c *fiber.Ctx) error {

	// 	//{
	//     "nama":"Bagja Lazwardi",
	//     "email":"rendi@gmail.com",
	//     "password":"1231431",
	//     "alamat":"Tangerang"
	// }
	var Users entity.Users
	if err := c.BodyParser(&Users); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat akun",
		"data":  Users,
	})

}

func CreateUserCara2(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}
	user := entity.Users{
		Nama:     data["nama"],
		Alamat:   data["alamat"],
		Email:    data["email"],
		Password: data["password"],
		NoTelpon: data["no_telpon"],
	}
	fmt.Println(data)
	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat akun",
		"data":  user,
	})

}

func GetUserByParameter(c *fiber.Ctx) error {

	//params
	nama := c.Params("nama")

	//id
	id := c.Query("id")

	category := c.Query("category")

	return c.JSON(fiber.Map{
		"Pesan":      "Berhasil Get Data Parameter",
		"data":       nama,
		"data_query": id,
		"category":   category,
	})
}

func UpdateUsers(c *fiber.Ctx) error {

	//ambil data dari database berdasarkan id
	id := c.QueryInt("id")

	//Cari data dengan Parameter Id
	user := entity.Users{
		ID:       uint(id),
		Nama:     "Bagja",
		Alamat:   "Tangerang",
		Email:    "testing@gmail.com",
		Password: "12312312",
		NoTelpon: "31231231",
	}

	dataLama := user

	var request entity.Users
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}
	user.Nama = request.Nama
	user.Alamat = request.Alamat
	user.Email = request.Email
	user.Password = request.Password

	//Simpan database

	return c.JSON(fiber.Map{
		"pesan":     "berhasil update data",
		"data_lama": dataLama,
		"data":      user,
	})
}

func DeleteUsers(c *fiber.Ctx) error {

	return nil
}
