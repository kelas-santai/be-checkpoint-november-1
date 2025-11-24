package entity

type Category struct {
	IdCategory uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama       string    `json:"nama"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	Product    []Product `gorm:"foreignKey:IdCategory" json:"products"`
	CreateAt   string    `json:"create_at"`
	UpdateAt   string    `json:"update_at"`
}

type Product struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	IdCategory  uint   `json:"id_category"`
	Nama        string `json:"nama"`
	Description string `gorm:"type:text" json:"description"`
	Harga       string `json:"harga"`
	Gambar      string `json:"gambar"`
	Stock       int    `gorm:"default:100" json:"stock"`
	IsAvailable bool   `gorm:"default:true" json:"is_available"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

//mennyimpan informasi qr code dan nomor meja
type Table struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Number    int    `json:"number"`   //nomor meha
	QRCode    string `json:"qr_code"`  // ini adalah gambar atau path dari qr code
	QRToken   string `json:"qr_token"` //token acak untuk generate qr code
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
