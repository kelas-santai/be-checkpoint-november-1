package entity

type Order struct {
	IdOrder       uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber   string      `json:"order_number"`
	Table         string      `json:"table"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	Status        string      `json:"status"` // pending, processing, completed, cancelled
	PaymentMethod string      `json:"payment_method"`
	PaymentStatus string      `json:"payment_status"` // unpaid, paid
	TotalAmount   float64     `json:"total_amount"`
	Notes         string      `json:"notes"`
	TrackingToken string      `json:"tracking_token"`
	OrderItems    []OrderItem `gorm:"foreignKey:IdOrder" json:"order_items"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

// OrderItem untuk detail item dalam pesanan
type OrderItem struct {
	IdOrderItem uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	IdOrder     uint    `json:"id_order"`
	ProductID   uint    `json:"product_id"`
	Product     string  `json:"product"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
	Notes       string  `json:"notes"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
