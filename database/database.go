package database

import (
	"fmt"
	"meeting3/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	user := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	databaseName := "meeting4-november"
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	con := user + ":" + password + "@tcp" + "(" + host + ":" + port + ")/" + databaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(con), &gorm.Config{})
	if err != nil {
		fmt.Println("eror")
		fmt.Println(err.Error())
	}
	DB = db

	//SAya akan bikin migrassi untuk table yang akan di buat

	DB.AutoMigrate(
		&entity.Users{},
		&entity.Admin{},
		&entity.Category{},
		&entity.Product{},

		&entity.Order{},
		&entity.OrderItem{},
		&entity.Table{},
	)

}
