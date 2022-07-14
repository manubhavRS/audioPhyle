package models

type AddCartModel struct {
	UserID    string                `db:"user_id" json:"userID"`
	Products  []AddCartProductModel `db:"products" json:"products"`
	TotalCost float64               `db:"total_cost" json:"totalCost"`
}
type UpdateCartModel struct {
	CartID    string                `db:"id" json:"cartID"`
	UserID    string                `db:"user_id" json:"userID"`
	Products  []AddCartProductModel `db:"products" json:"products"`
	TotalCost float64               `db:"total_cost" json:"totalCost"`
}
type AddCartProductModel struct {
	ProductID string  `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
}
type ProductPriceModel struct {
	ProductID string  `db:"id" json:"id"`
	Price     float64 `db:"price" json:"price"`
}
type CartModel struct {
	UserID                string                `db:"user_id" json:"userID"`
	Products              []AddCartProductModel `db:"products" json:"products"`
	Cost                  []float64             `db:"price" json:"price"`
	TotalCostWithShipping float64
}
