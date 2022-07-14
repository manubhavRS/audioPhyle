package helper

import (
	"audioPhile/models"
	"audioPhile/utilities"
	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	"log"
)

func AddOrderHelper(order models.AddOrderModel, tx *sqlx.Tx) error {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertBuilder := psql.Insert("orders").Columns("product_id", "user_id", "address_id", "cost", "payment_by_card", "payment_by_cod", "date_of_delivery")
	for i, cart := range order.CartDetails.Products {
		insertBuilder.Values(cart.ProductID, order.UserID, order.AddressID, (1.2*order.CartDetails.Products[i].Price)+utilities.ShippingCharges, order.PaymentByCard, order.PaymentByCod, order.DateOfDelivery)
	}
	sql, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("AddCartProductsHelper Error: %v", err)
		return err
	}
	return nil
}
