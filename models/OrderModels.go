package models

import (
	"time"
)

type AddOrderModel struct {
	UserID         string `db:"user_id" json:"userID"`
	AddressID      string `db:"address_id" json:"addressID"`
	CartID         string `db:"cart_id" json:"cartID"`
	CartDetails    CartModel
	PaymentByCard  string    `db:"payment_by_card" json:"paymentByCard"`
	PaymentByCod   bool      `db:"payment_by_cod" json:"paymentByCod"`
	DateOfDelivery time.Time `db:"date_of_delivery" json:"dateOfDelivery"`
}
