package entity

type Users struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Alamat   string `json:"alamat"`
	NoTelpon string `json:"no_telpon"`
}

type Admin struct {
	ID       uint   `json:"id"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Alamat   string `json:"alamat"`
	NoTelpon string `json:"no_telpon"`
}
