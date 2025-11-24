package controllers

import (
	"meeting3/database"
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

	//data bisa langsung di simpan ke databases
	database.DB.Create(&Users)

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

	//Cek dulu data email ga boleh sama, dan juga data no telpon itu ga boleh sama
	//get  data ke database

	var exixting entity.Users
	//dia akan mencari semua data kalau kita pake fungsi .find()
	//database.DB.Find()

	//kalo mencari 1 data saja kita pake .First atau  akan mencari data yang paling atas
	//database.DB.First()

	//kalo mencari 1 data saja kita pake .Last data yang paling bawah
	//database.DB.Last()

	//mengambil dat
	database.DB.Where("email = ? OR no_telpon = ? ", user.Email, user.NoTelpon).First(&exixting)
	if exixting.ID != 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "gagal buat akun, karena ada data user dengan nama" + exixting.Nama,
		})
	}

	database.DB.Create(&user)

	return c.JSON(fiber.Map{
		"pesan": "berhasil membuat akun",
		"data":  user,
	})

}

func GetUserByParameter(c *fiber.Ctx) error {

	//id
	id := c.Query("id")
	var user entity.Users
	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(fiber.Map{
		"Pesan": "Berhasil Get Data User By Id",
		"data":  user,
	})
}

func GetAllUser(c *fiber.Ctx) error {

	var allUsers []entity.Users
	database.DB.Find(&allUsers)

	return c.JSON(fiber.Map{
		"pesan": "berhasil mendapatkan semua user",
		"data":  allUsers,
	})
}

func UpdateUsers(c *fiber.Ctx) error {

	//ambil data dari database berdasarkan id
	id := c.QueryInt("id")

	//Cari data dengan Parameter Id

	var User entity.Users
	database.DB.Where("id = ?", id).First(&User)
	if User.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "Tidak ada Data di database untuk ID itu",
		})
	}

	var request entity.Users
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"pesan": "gagal untuk body parser",
			"err":   err,
		})
	}
	User.Nama = request.Nama
	User.Alamat = request.Alamat
	User.Email = request.Email
	User.Password = request.Password
	User.NoTelpon = request.NoTelpon

	//Simpan database
	database.DB.Save(&User)

	return c.JSON(fiber.Map{
		"pesan": "berhasil update data",
		"data":  User,
	})
}

func DeleteUsers(c *fiber.Ctx) error {
	id := c.Query("id")
	var user entity.Users
	database.DB.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"pesan": "Tidak ada Data di database",
		})
	}
	//Tinggal Hapus datanya
	database.DB.Delete(&user)
	return c.JSON(fiber.Map{
		"pesan": "Data telah di delete",
	})
}

func GetRawQuery(c *fiber.Ctx) error {

	type Result struct {
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}

	//tempat nampng data
	var result []Result

	//Query mentah
	raw := `SELECT nama, email FROM users`
	database.DB.Raw(raw).Scan(&result)

	return c.JSON(fiber.Map{
		"data": result,
	})
}
