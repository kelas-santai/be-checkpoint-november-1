package entity

type Category struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type Product struct {
	ID         uint   `json:"id"`
	IdCategory uint   `json:"id_category"`
	Nama       string `json:"nama"`
	Harga      string `json:"harga"`
	Gambar     string `json:"gambar"`
	CreateAt   string `json:"create_at"`
	UpdateAt   string `json:"update_at"`
}
