package repository

type CartDB struct {
	CartID   string
	UserID   string
	Products []ProductCartDB
}

type OrderDetailDB struct {
	Products  int32
	Discounts int32
	Shipping  int32
	Order     int32
}

type OrderDB struct {
	CartID string
	Totals OrderDetailDB
}

type ProductDB struct {
	ProductID string
	Name      string
	Category  string
	Price     float32
}

type ProductCartDB struct {
	ProductDB
	Quantity int32
}
