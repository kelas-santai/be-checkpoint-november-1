package entity

type Users struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Alamat    string `json:"alamat"`
	NoTelpon  string `json:"no_telpon"`
	Status    bool   `json:"status"` //aktif atau tidak
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//seperti kasir atau super admin
type Admin struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Role      string `json:"role"` //kasir atau super admin
	Password  string `json:"-"`
	Alamat    string `json:"alamat"`
	NoTelpon  string `json:"no_telpon"`
	Status    bool   `json:"status"` //aktif atau tidak
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
